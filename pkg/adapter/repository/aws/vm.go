package aws

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
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
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
	"strings"
)

const (
	tagKeyValueSeparator = "="
)

func (aws *awsRepository) FindVM(ctx context.Context, opt option.Option) (*instance.VM, errors.Error) {
	if !opt.SetType(reflect.TypeOf(resourceRepo.FindResourceOption{}).String()).Validate() {
		return nil, errors.InvalidOption.WithMessage(fmt.Sprintf("invalid option : want %s, got %+v", reflect.TypeOf(resourceRepo.FindResourceOption{}).String(), opt.Get()))
	}
	req := opt.Get().(resourceRepo.FindResourceOption)
	awsVM := &v1beta1.Instance{}
	if err := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: req.Name}, awsVM); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, errors.NotFound.WithMessage(fmt.Sprintf("vm %s not found", req.Name))
		}
		return nil, errors.KubernetesError.WithMessage(fmt.Sprintf("unable to get vm %s", req.Name))
	}

	awsKeyPair, err := aws._getKeyPair(ctx, awsVM)
	if !err.IsOk() {
		return nil, err
	}

	mod, err := aws.toModelVM(awsVM, awsKeyPair)
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
	awsVMList := &v1beta1.InstanceList{}
	listOpt := &client.ListOptions{
		LabelSelector: client.MatchingLabelsSelector{Selector: labels.SelectorFromSet(req.Labels)},
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
	//TODO implement me
	//panic("implement me")
	return &instance.VMCollection{}, errors.OK
}

func (aws *awsRepository) CreateVM(ctx context.Context, vm *instance.VM) errors.Error {
	if exists, err := aws.VMExists(ctx, vm); !err.IsOk() {
		return err
	} else if exists {
		return errors.Conflict.WithMessage(fmt.Sprintf("vm %s already exists in subnet %s", vm.IdentifierName.VM, vm.IdentifierID.Subnetwork))
	}
	awsVM, err := aws.toAWSVM(ctx, vm)
	if !err.IsOk() {
		return err
	}
	if err := kubernetes.Client().Client.Create(ctx, awsVM); err != nil {
		if k8serrors.IsAlreadyExists(err) {
			return errors.Conflict.WithMessage(fmt.Sprintf("vm %s already exists", vm.IdentifierName.VM))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to create vm %s", vm.IdentifierName.VM))
	}

	if err := aws._createKeyPair(ctx, vm); !err.IsOk() {
		return err
	}
	return errors.Created
}

func (aws *awsRepository) UpdateVM(ctx context.Context, vm *instance.VM) errors.Error {
	existingVM := &v1beta1.Instance{}
	if err := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: vm.IdentifierID.VM}, existingVM); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("vm %s not found", vm.IdentifierID.VM))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to get vm %s", vm.IdentifierID.Subnetwork))
	}
	awsVM, err := aws.toAWSVM(ctx, vm)
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

	if err := aws._updateKeyPair(ctx, vm); !err.IsOk() {
		return err
	}
	return errors.NoContent
}

func (aws *awsRepository) DeleteVM(ctx context.Context, vm *instance.VM) errors.Error {
	awsVM, err := aws.toAWSVM(ctx, vm)
	if !err.IsOk() {
		return err
	}
	if err := kubernetes.Client().Client.Delete(ctx, awsVM); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("vm %s not found", vm.IdentifierName.VM))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to delete vm %s", vm.IdentifierName.VM))
	}

	if err := aws._deleteKeyPair(ctx, vm); !err.IsOk() {
		return err
	}
	return errors.NoContent
}

func (aws *awsRepository) VMExists(ctx context.Context, vm *instance.VM) (bool, errors.Error) {
	awsVMs := &v1beta1.InstanceList{}
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
	if err := kubernetes.Client().Client.List(ctx, awsVMs, listOpt); err != nil {
		return false, errors.KubernetesError.WithMessage("unable to list vm")
	}
	return len(awsVMs.Items) > 0, errors.OK
}

func (aws *awsRepository) toModelVM(vm *v1beta1.Instance, keyPair *v1beta1.KeyPair) (*instance.VM, errors.Error) {
	id := identifier.VM{}
	name := identifier.VM{}
	if err := id.IDFromLabels(vm.Labels); !err.IsOk() {
		return nil, err
	}
	if err := name.NameFromLabels(vm.Labels); !err.IsOk() {
		return nil, err
	}
	tags := map[string]string{}
	for _, tag := range vm.Spec.ForProvider.Tags {
		split := strings.Split(*tag, tagKeyValueSeparator)
		if len(split) != 2 {
			continue
		}
		tags[split[0]] = split[1]
	}

	sshKeys := crossplane.FromSSHKeySecretLabels(keyPair.Labels)
	publicIP := ""
	if vm.Status.AtProvider.PublicIP != nil {
		publicIP = *vm.Status.AtProvider.PublicIP
	}
	hasPublicIp, err := strconv.ParseBool(vm.ObjectMeta.Labels[crossplane.VMPublicIPLabel])
	if err != nil {
		return nil, errors.InternalError.WithMessage("unable to parse public ip label")
	}
	vmOS := instance.VMOS{}
	if vm.Spec.ForProvider.AMI != nil {
		vmOS = toVmOS(vm.Spec.ForProvider.AMI)
	}
	return &instance.VM{
		Metadata: metadata.Metadata{
			Status:  metadata.StatusFromKubernetesStatus(vm.Status.Conditions),
			Managed: vm.Spec.ResourceSpec.DeletionPolicy == v1.DeletionDelete,
			Tags:    tags,
		},
		IdentifierID:   id,
		IdentifierName: name,
		AssignPublicIP: hasPublicIp,
		PublicIP:       publicIP,
		Zone:           *vm.Spec.ForProvider.Region,
		MachineType:    *vm.Spec.ForProvider.InstanceType,
		Auths:          sshKeys,
		Disks:          aws.toVMDiskCollection(vm.Spec.ForProvider.RootBlockDevice),
		OS:             vmOS,
	}, errors.OK
}

func (aws *awsRepository) toAWSVM(ctx context.Context, vm *instance.VM) (*v1beta1.Instance, errors.Error) {
	sshKeysLabels := crossplane.ToSSHKeySecretLabels(vm.Auths)
	asPublicIPLabel := map[string]string{crossplane.VMPublicIPLabel: strconv.FormatBool(vm.AssignPublicIP)}
	instanceLabels := lo.Assign(crossplane.GetBaseLabels(ctx.Value(context.ProjectIDKey).(string)), vm.IdentifierID.ToIDLabels(), vm.IdentifierName.ToNameLabels(), sshKeysLabels, asPublicIPLabel)
	subnetID := identifier.Subnetwork{Provider: vm.IdentifierID.Provider, VPC: vm.IdentifierID.VPC, Network: vm.IdentifierID.Network, Subnetwork: vm.IdentifierID.Subnetwork}
	keyName := fmt.Sprintf("%s-keypair", vm.IdentifierID.VM)

	awsDiskList, err := aws.toAWSVMDiskList(vm.Disks, vm.OS)
	if !err.IsOk() {
		return nil, err
	}

	return &v1beta1.Instance{
		ObjectMeta: metav1.ObjectMeta{
			Name:        vm.IdentifierID.VM,
			Labels:      instanceLabels,
			Annotations: crossplane.GetAnnotations(vm.Metadata.Managed, vm.IdentifierName.Network),
		},
		Spec: v1beta1.InstanceSpec{
			ResourceSpec: v1.ResourceSpec{
				DeletionPolicy: crossplane.GetDeletionPolicy(vm.Metadata.Managed),
				ProviderConfigReference: &v1.Reference{
					Name: vm.IdentifierID.VM,
				},
			},
			ForProvider: v1beta1.InstanceParameters{
				AMI:                      &vm.OS.Name,
				AssociatePublicIPAddress: &vm.AssignPublicIP,
				InstanceType:             &vm.MachineType,
				KeyName:                  &keyName,
				Region:                   &vm.Zone,
				RootBlockDevice:          *awsDiskList,
				SubnetIDSelector: &v1.Selector{
					MatchLabels: subnetID.ToIDLabels(),
				},
			},
		},
	}, errors.OK
}

func (aws *awsRepository) toModelVMCollection(ctx context.Context, instanceList *v1beta1.InstanceList) (*instance.VMCollection, errors.Error) {
	items := instance.VMCollection{}
	for _, item := range instanceList.Items {
		keyPair, err := aws._getKeyPair(ctx, &item)
		if !err.IsOk() {
			return nil, err
		}

		vm, err := aws.toModelVM(&item, keyPair)
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

func (aws *awsRepository) toAWSVMDiskList(disks []instance.VMDisk, os instance.VMOS) (*[]v1beta1.RootBlockDeviceParameters, errors.Error) {
	var bootDisks []v1beta1.RootBlockDeviceParameters
	for _, disk := range disks {
		awsDisk, err := aws.toAWSVMDisk(disk)
		if !err.IsOk() {
			return nil, err
		}
		bootDisks = append(bootDisks, *awsDisk)
	}
	return &bootDisks, errors.OK
}

func (aws *awsRepository) toAWSVMDisk(disk instance.VMDisk) (*v1beta1.RootBlockDeviceParameters, errors.Error) {
	if disk.Type != instance.DiskTypeSSD {
		return nil, errors.BadRequest.WithMessage("Root block requires an SSD disk in aws.")
	}

	diskSize := float64(disk.SizeGib)
	diskType := "gp2"
	//diskMode := toAWSDiskMode(disk.Mode)
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
