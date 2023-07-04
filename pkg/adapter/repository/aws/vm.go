package aws

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/aws/xrds"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/crossplane"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/instance"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	resourceRepo "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository/resource"
	v1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/samber/lo"
	"github.com/upbound/provider-aws/apis/ec2/v1beta1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	tagKeyValueSeparator = "="
)

func (aws *awsRepository) FindVM(ctx context.Context, opt option.Option) (*instance.VM, errors.Error) {
	if !opt.SetType(reflect.TypeOf(resourceRepo.FindResourceOption{}).String()).Validate() {
		return nil, errors.InvalidOption.WithMessage(fmt.Sprintf("invalid option : want %s, got %+v", reflect.TypeOf(resourceRepo.FindResourceOption{}).String(), opt.Get()))
	}
	req := opt.Get().(resourceRepo.FindResourceOption)
	awsVM := &xrds.VMInstance{}
	if err := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: ctx.Value(context.CurrentNamespace).(string)}, awsVM); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, errors.NotFound.WithMessage(fmt.Sprintf("vm %s not found", req.Name))
		}
		return nil, errors.KubernetesError.WithMessage(fmt.Sprintf("unable to get vm %s", req.Name))
	}

	mod, err := aws.toModelVM(awsVM)
	if !err.IsOk() {
		return nil, err
	}
	return mod, errors.OK
}

func (aws *awsRepository) FindAllRecursiveVMs(ctx context.Context, opt option.Option, ch *resourceRepo.VMChannel) {
	if !opt.SetType(reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String()).Validate() {
		ch.ErrorChannel <- errors.InvalidOption.WithMessage(fmt.Sprintf("invalid option : want %s, got %+v", reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String(), opt.Get()))
		return
	}
	req := opt.Get().(resourceRepo.FindAllResourceOption)
	awsVMList := &xrds.VMInstanceList{}
	listOpt := &client.ListOptions{
		LabelSelector: client.MatchingLabelsSelector{Selector: labels.SelectorFromSet(req.Labels)},
		Namespace:     ctx.Value(context.CurrentNamespace).(string),
	}
	if err := kubernetes.Client().Client.List(ctx, awsVMList, listOpt); err != nil {
		ch.ErrorChannel <- errors.KubernetesError.WithMessage(fmt.Sprintf("unable to list vm"))
		return
	}
	if firewalls, err := aws.toModelVMCollection(ctx, awsVMList); !err.IsOk() {
		ch.ErrorChannel <- err
	} else {
		ch.Channel <- firewalls
	}
}

func (aws *awsRepository) FindAllVMs(ctx context.Context, opt option.Option) (*instance.VMCollection, errors.Error) {
	if !opt.SetType(reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String()).Validate() {
		return nil, errors.BadRequest.WithMessage(fmt.Sprintf("invalid option : want %s, got %+v", reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String(), opt.Get()))
	}
	req := opt.Get().(resourceRepo.FindAllResourceOption)
	awsVMList := &xrds.VMInstanceList{}
	listOpt := &client.ListOptions{
		LabelSelector: client.MatchingLabelsSelector{Selector: labels.SelectorFromSet(req.Labels)},
		Namespace:     ctx.Value(context.CurrentNamespace).(string),
	}
	if err := kubernetes.Client().Client.List(ctx, awsVMList, listOpt); err != nil {
		fmt.Println(err.Error())
		return nil, errors.KubernetesError.WithMessage("unable to list vms")
	}

	vmList, err := aws.toModelVMCollection(ctx, awsVMList)
	if !err.IsOk() {
		return nil, err
	}
	return vmList, errors.OK
}

func (aws *awsRepository) CreateVM(ctx context.Context, vm *instance.VM) errors.Error {
	if exists, err := aws.VMExists(ctx, vm); !err.IsOk() {
		return err
	} else if exists {
		return errors.Conflict.WithMessage(fmt.Sprintf("vm %s already exists in subnet %s", vm.IdentifierName.VM, vm.IdentifierID.Subnetwork))
	}

	vmInstance, err := aws.toVMInstance(ctx, vm)
	if !err.IsOk() {
		return err
	}
	if err := kubernetes.Client().Client.Create(ctx, vmInstance); err != nil {
		fmt.Println(err.Error())
		if k8serrors.IsAlreadyExists(err) {
			return errors.Conflict.WithMessage(fmt.Sprintf("vm %s already exists", vm.IdentifierName.VM))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to create vm %s", vm.IdentifierName.VM))
	}

	return errors.Created
}

func (aws *awsRepository) UpdateVM(ctx context.Context, vm *instance.VM) errors.Error {
	existingVM := &xrds.VMInstance{}
	if err := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: vm.IdentifierID.VM, Namespace: ctx.Value(context.CurrentNamespace).(string)}, existingVM); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("vm %s not found", vm.IdentifierID.VM))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to get vm %s", vm.IdentifierID.Subnetwork))
	}

	awsVM, err := aws.toVMInstance(ctx, vm)
	if !err.IsOk() {
		return err
	}
	existingVM.Spec = awsVM.Spec
	existingVM.Labels = awsVM.Labels
	if err := kubernetes.Client().Client.Update(ctx, existingVM); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("vm %s not found", vm.IdentifierName.VM))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to update vm %s", vm.IdentifierName.VM))
	}
	return errors.NoContent
}

func (aws *awsRepository) DeleteVM(ctx context.Context, vm *instance.VM) errors.Error {
	awsVM, err := aws.toVMInstance(ctx, vm)
	if !err.IsOk() {
		return err
	}

	if err := kubernetes.Client().Client.Delete(ctx, awsVM); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("vm %s not found", vm.IdentifierName.VM))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to delete vm %s", vm.IdentifierName.VM))
	}
	return errors.NoContent
}

func (aws *awsRepository) VMExists(ctx context.Context, vm *instance.VM) (bool, errors.Error) {
	vmInstances := &xrds.VMInstanceList{}
	parentID := identifier.Subnetwork{
		Provider:   vm.IdentifierID.Provider,
		VPC:        vm.IdentifierID.VPC,
		Network:    vm.IdentifierID.Network,
		Subnetwork: vm.IdentifierID.Subnetwork,
	}
	searchLabels := lo.Assign(parentID.ToIDLabels(), map[string]string{identifier.VMNameKey: vm.IdentifierName.VM})
	listOpt := &client.ListOptions{
		LabelSelector: client.MatchingLabelsSelector{Selector: labels.SelectorFromSet(searchLabels)},
	}
	if err := kubernetes.Client().Client.List(ctx, vmInstances, listOpt); err != nil {
		return false, errors.KubernetesError.WithMessage("unable to list vm")
	}
	return len(vmInstances.Items) > 0, errors.OK
}

func (aws *awsRepository) toModelVM(vm *xrds.VMInstance) (*instance.VM, errors.Error) {
	id := identifier.VM{}
	name := identifier.VM{}
	if err := id.IDFromLabels(vm.Labels); !err.IsOk() {
		return nil, err
	}
	if err := name.NameFromLabels(vm.Labels); !err.IsOk() {
		return nil, err
	}

	sshKeys := crossplane.FromSSHKeySecretLabels(vm.Labels)
	publicIP := ""
	if vm.Status.PublicIp != nil {
		publicIP = *vm.Status.PublicIp
	}

	hasPublicIp, err := strconv.ParseBool(vm.ObjectMeta.Labels[crossplane.VMPublicIPLabel])
	if err != nil {
		return nil, errors.InternalError.WithMessage("unable to parse public ip label")
	}
	vmOS := instance.VMOS{}
	if vm.Spec.Parameters.Os != nil {
		vmOS = toVmOS(vm.Spec.Parameters.Os)
	}
	return &instance.VM{
		Metadata: metadata.Metadata{
			Status:  metadata.StatusFromKubernetesStatus(vm.Status.Conditions),
			Managed: vm.Spec.Parameters.DeletionPolicy == v1.DeletionDelete,
			Tags:    make(map[string]string),
		},
		IdentifierID:   id,
		IdentifierName: name,
		AssignPublicIP: hasPublicIp,
		PublicIP:       publicIP,
		Zone:           *vm.Spec.Parameters.Region,
		MachineType:    *vm.Spec.Parameters.MachineType,
		Auths:          sshKeys,
		Disks:          aws.toVMDiskCollection(vm.Spec.Parameters.Disks),
		OS:             vmOS,
	}, errors.OK
}

func (aws *awsRepository) toVMInstance(ctx context.Context, vm *instance.VM) (*xrds.VMInstance, errors.Error) {
	sshKeysLabels := crossplane.ToSSHKeySecretLabels(vm.Auths)
	asPublicIPLabel := map[string]string{crossplane.VMPublicIPLabel: strconv.FormatBool(vm.AssignPublicIP)}
	labels := lo.Assign(crossplane.GetBaseLabels(ctx.Value(context.ProjectIDKey).(string)), vm.IdentifierID.ToIDLabels(), vm.IdentifierName.ToNameLabels(), sshKeysLabels, asPublicIPLabel)

	keyPairRef := fmt.Sprintf("%s-keypair", vm.IdentifierID.VM)
	securityGroupRef := fmt.Sprintf("%s-security-group", vm.IdentifierID.VM)
	disks, err := aws.toAWSVMDiskList(vm.Disks, vm.OS)
	if !err.IsOk() {
		return nil, err
	}

	return &xrds.VMInstance{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: crossplane.GetAnnotations(vm.Metadata.Managed, vm.IdentifierName.VM),
			Labels:      labels,
			Name:        vm.IdentifierID.VM,
			Namespace:   ctx.Value(context.CurrentNamespace).(string),
		},
		Spec: xrds.VMInstanceSpec{
			Parameters: xrds.VMInstanceParameters{
				AssignPublicIp:   &vm.AssignPublicIP,
				DeletionPolicy:   crossplane.GetDeletionPolicy(vm.Metadata.Managed),
				Disks:            disks,
				KeyPairRef:       &keyPairRef,
				MachineType:      &vm.MachineType,
				NetworkRef:       &vm.IdentifierID.Network,
				Os:               &vm.OS.Name,
				ProviderRef:      &vm.IdentifierID.Provider,
				PublicKey:        &vm.Auths[0].PublicKey,
				Region:           &vm.Zone,
				SecurityGroupRef: &securityGroupRef,
				SubnetworkRef:    &vm.IdentifierID.Subnetwork,
				VmId:             &vm.IdentifierID.VM,
			},
		},
	}, errors.OK
}

func (aws *awsRepository) toModelVMCollection(ctx context.Context, instanceList *xrds.VMInstanceList) (*instance.VMCollection, errors.Error) {
	items := instance.VMCollection{}
	for _, item := range instanceList.Items {
		vm, err := aws.toModelVM(&item)
		if !err.IsOk() {
			return nil, err
		}
		items[vm.IdentifierName.VM] = *vm
	}
	return &items, errors.OK
}

func toVmOS(ami *string) instance.VMOS {
	if ami == nil {
		return instance.VMOS{}
	}
	return instance.VMOS{
		ID:   *ami,
		Name: *ami,
	}
}

func (aws *awsRepository) toVMDiskCollection(disks []v1beta1.RootBlockDeviceParameters) []instance.VMDisk {
	var ret []instance.VMDisk
	for _, disk := range disks {
		ret = append(ret, aws.toVMDisk(&disk))
	}
	return ret
}

func (aws *awsRepository) toVMDisk(disk *v1beta1.RootBlockDeviceParameters) instance.VMDisk {
	return instance.VMDisk{
		AutoDelete: *disk.DeleteOnTermination,
		SizeGib:    int(*disk.VolumeSize),
		Type:       instance.DiskTypeSSD,
		Mode:       instance.DiskModeReadWrite,
	}
}

func (aws *awsRepository) toAWSVMDiskList(disks []instance.VMDisk, os instance.VMOS) ([]v1beta1.RootBlockDeviceParameters, errors.Error) {
	var bootDisks []v1beta1.RootBlockDeviceParameters
	for _, disk := range disks {
		awsDisk, err := aws.toAWSVMDisk(disk)
		if !err.IsOk() {
			return nil, err
		}
		bootDisks = append(bootDisks, *awsDisk)
	}
	return bootDisks, errors.OK
}

func (aws *awsRepository) toAWSVMDisk(disk instance.VMDisk) (*v1beta1.RootBlockDeviceParameters, errors.Error) {
	if disk.Type != instance.DiskTypeSSD {
		return nil, errors.BadRequest.WithMessage("Root block requires an SSD disk in aws.")
	}

	diskSize := float64(disk.SizeGib)
	diskType := "gp2"
	return &v1beta1.RootBlockDeviceParameters{
		DeleteOnTermination: &disk.AutoDelete,
		VolumeSize:          &diskSize,
		VolumeType:          &diskType,
	}, errors.OK
}

func sshKeysToString(sshKeys model.SSHKeyList) *string {
	var ret string
	for _, key := range sshKeys {
		ret += fmt.Sprintf("%s:%s\n", key.Username, key.PublicKey)
	}
	ret = strings.TrimSuffix(ret, "\n")
	return &ret
}
