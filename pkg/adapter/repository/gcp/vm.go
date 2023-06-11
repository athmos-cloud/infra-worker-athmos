package gcp

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
	"github.com/upbound/provider-gcp/apis/compute/v1beta1"
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

func (gcp *gcpRepository) FindVM(ctx context.Context, opt option.Option) (*instance.VM, errors.Error) {
	if !opt.SetType(reflect.TypeOf(resourceRepo.FindResourceOption{}).String()).Validate() {
		return nil, errors.InvalidOption.WithMessage(fmt.Sprintf("invalid option : want %s, got %+v", reflect.TypeOf(resourceRepo.FindResourceOption{}).String(), opt.Get()))
	}
	req := opt.Get().(resourceRepo.FindResourceOption)
	gcpVM := &v1beta1.Instance{}
	if err := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: req.Name}, gcpVM); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, errors.NotFound.WithMessage(fmt.Sprintf("vm %s not found", req.Name))
		}
		return nil, errors.KubernetesError.WithMessage(fmt.Sprintf("unable to get vm %s", req.Name))
	}
	mod, err := gcp.toModelVM(gcpVM)
	if !err.IsOk() {
		return nil, err
	}
	return mod, errors.OK
}

func (gcp *gcpRepository) FindAllRecursiveVMs(ctx context.Context, opt option.Option, ch *resourceRepo.VMChannel) {
	if !opt.SetType(reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String()).Validate() {
		ch.ErrorChannel <- errors.InvalidOption.WithMessage(fmt.Sprintf("invalid option : want %s, got %+v", reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String(), opt.Get()))
		return
	}
	req := opt.Get().(resourceRepo.FindAllResourceOption)
	gcpVMList := &v1beta1.InstanceList{}
	listOpt := &client.ListOptions{
		LabelSelector: client.MatchingLabelsSelector{Selector: labels.SelectorFromSet(req.Labels)},
	}
	if err := kubernetes.Client().Client.List(ctx, gcpVMList, listOpt); err != nil {
		ch.ErrorChannel <- errors.KubernetesError.WithMessage("unable to list vm")
		return
	}
	if firewalls, err := gcp.toModelVMCollection(gcpVMList); !err.IsOk() {
		ch.ErrorChannel <- err
	} else {
		ch.Channel <- firewalls
	}
}

func (gcp *gcpRepository) FindAllVMs(ctx context.Context, opt option.Option) (*instance.VMCollection, errors.Error) {
	if !opt.SetType(reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String()).Validate() {
		return nil, errors.InvalidOption.WithMessage(fmt.Sprintf("invalid option : want %s, got %+v", reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String(), opt.Get()))
	}
	req := opt.Get().(resourceRepo.FindAllResourceOption)
	gcpVMList := &v1beta1.InstanceList{}
	kubeOptions := &client.ListOptions{
		LabelSelector: client.MatchingLabelsSelector{Selector: labels.SelectorFromSet(req.Labels)},
	}
	if err := kubernetes.Client().Client.List(ctx, gcpVMList, kubeOptions); err != nil {
		return nil, errors.KubernetesError.WithMessage(fmt.Sprintf("unable to list vm"))
	}
	modVms, err := gcp.toModelVMCollection(gcpVMList)
	if !err.IsOk() {
		return nil, err
	}

	return modVms, errors.OK
}

func (gcp *gcpRepository) CreateVM(ctx context.Context, vm *instance.VM) errors.Error {
	if exists, err := gcp.VMExists(ctx, vm); !err.IsOk() {
		return err
	} else if exists {
		return errors.Conflict.WithMessage(fmt.Sprintf("vm %s already exists in subnet %s", vm.IdentifierName.VM, vm.IdentifierID.Subnetwork))
	}
	gcpVM := gcp.toGCPVM(ctx, vm)
	if err := kubernetes.Client().Client.Create(ctx, gcpVM); err != nil {
		if k8serrors.IsAlreadyExists(err) {
			return errors.Conflict.WithMessage(fmt.Sprintf("vm %s already exists", vm.IdentifierName.VM))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to create vm %s", vm.IdentifierName.VM))
	}
	return errors.Created
}

func (gcp *gcpRepository) UpdateVM(ctx context.Context, vm *instance.VM) errors.Error {
	existingVM := &v1beta1.Instance{}
	if err := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: vm.IdentifierID.VM}, existingVM); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("vm %s not found", vm.IdentifierID.VM))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to get vm %s", vm.IdentifierID.VM))
	}
	gcpVM := gcp.toGCPVM(ctx, vm)
	existingVM.Spec = gcpVM.Spec
	existingVM.Labels = gcpVM.Labels
	if err := kubernetes.Client().Client.Update(ctx, existingVM); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("vm %s not found", vm.IdentifierName.VM))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to update vm %s", vm.IdentifierName.VM))
	}
	return errors.NoContent
}

func (gcp *gcpRepository) DeleteVM(ctx context.Context, vm *instance.VM) errors.Error {
	gcpVM := &v1beta1.Instance{}
	if err := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: vm.IdentifierID.VM}, gcpVM); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("vm %s not found", vm.IdentifierID.VM))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to get vm %s", vm.IdentifierID.VM))
	}
	if err := kubernetes.Client().Client.Delete(ctx, gcpVM); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("vm %s not found", vm.IdentifierName.VM))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to delete vm %s", vm.IdentifierName.VM))
	}

	return errors.NoContent
}

func (gcp *gcpRepository) VMExists(ctx context.Context, vm *instance.VM) (bool, errors.Error) {
	gcpVMs := &v1beta1.InstanceList{}
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
	if err := kubernetes.Client().Client.List(ctx, gcpVMs, listOpt); err != nil {
		return false, errors.KubernetesError.WithMessage("unable to list vm")
	}
	return len(gcpVMs.Items) > 0, errors.OK
}

func (gcp *gcpRepository) toModelVM(vm *v1beta1.Instance) (*instance.VM, errors.Error) {
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
	sshKeys := crossplane.FromSSHKeySecretLabels(vm.Labels)
	publicIP := ""
	if vm.Status.AtProvider.NetworkInterface != nil {
		publicIP = *vm.Status.AtProvider.NetworkInterface[0].AccessConfig[0].NATIP
	}
	hasPublicIp, err := strconv.ParseBool(vm.ObjectMeta.Labels[crossplane.VMPublicIPLabel])
	if err != nil {
		return nil, errors.InternalError.WithMessage("unable to parse public ip label")
	}
	vmOS := instance.VMOS{}
	if vm.Spec.ForProvider.BootDisk != nil {
		vmOS = toVmOS(&vm.Spec.ForProvider.BootDisk[0])
	}
	return &instance.VM{
		Metadata: metadata.Metadata{
			Managed: vm.Spec.ResourceSpec.DeletionPolicy == v1.DeletionDelete,
			Tags:    tags,
		},
		IdentifierID:   id,
		IdentifierName: name,
		AssignPublicIP: hasPublicIp,
		PublicIP:       publicIP,
		Zone:           *vm.Spec.ForProvider.Zone,
		MachineType:    *vm.Spec.ForProvider.MachineType,
		Auths:          sshKeys,
		Disks:          gcp.toVMDiskCollection(vm.Spec.ForProvider.BootDisk),
		OS:             vmOS,
	}, errors.OK
}

func (gcp *gcpRepository) toGCPVM(ctx context.Context, vm *instance.VM) *v1beta1.Instance {
	sshKeysLabels := crossplane.ToSSHKeySecretLabels(vm.Auths)
	asPublicIPLabel := map[string]string{crossplane.VMPublicIPLabel: strconv.FormatBool(vm.AssignPublicIP)}
	instanceLabels := lo.Assign(crossplane.GetBaseLabels(ctx.Value(context.ProjectIDKey).(string)), vm.IdentifierID.ToIDLabels(), vm.IdentifierName.ToNameLabels(), sshKeysLabels, asPublicIPLabel)
	networkID := identifier.Network{Provider: vm.IdentifierID.Provider, VPC: vm.IdentifierID.VPC, Network: vm.IdentifierID.Network}
	subnetID := identifier.Subnetwork{Provider: vm.IdentifierID.Provider, VPC: vm.IdentifierID.VPC, Network: vm.IdentifierID.Network, Subnetwork: vm.IdentifierID.Subnetwork}

	netInterface := []v1beta1.NetworkInterfaceParameters{
		{
			NetworkSelector: &v1.Selector{
				MatchLabels: networkID.ToIDLabels(),
			},
			SubnetworkSelector: &v1.Selector{
				MatchLabels: subnetID.ToIDLabels(),
			},
		},
	}
	if vm.AssignPublicIP {
		netInterface[0].AccessConfig = []v1beta1.AccessConfigParameters{{}}
	}

	return &v1beta1.Instance{
		ObjectMeta: metav1.ObjectMeta{
			Name:        vm.IdentifierID.VM,
			Labels:      instanceLabels,
			Annotations: crossplane.GetAnnotations(vm.Metadata.Managed, vm.IdentifierName.VM),
		},
		Spec: v1beta1.InstanceSpec{
			ResourceSpec: v1.ResourceSpec{
				DeletionPolicy: crossplane.GetDeletionPolicy(vm.Metadata.Managed),
				ProviderConfigReference: &v1.Reference{
					Name: vm.IdentifierID.VM,
				},
			},
			ForProvider: v1beta1.InstanceParameters{
				Project:     &vm.IdentifierID.VPC,
				MachineType: &vm.MachineType,
				Zone:        &vm.Zone,
				BootDisk:    gcp.toGCPVMDiskList(vm.Disks, vm.OS),
				Metadata: map[string]*string{
					"ssh-keys": sshKeysToString(vm.Auths),
				},
				NetworkInterface: netInterface,
			},
		},
	}
}

func (gcp *gcpRepository) toModelVMCollection(instanceList *v1beta1.InstanceList) (*instance.VMCollection, errors.Error) {
	items := instance.VMCollection{}
	for _, item := range instanceList.Items {
		vm, err := gcp.toModelVM(&item)
		if !err.IsOk() {
			return nil, err
		}
		items[vm.IdentifierName.VM] = *vm
	}
	return &items, errors.OK
}

func toVmOS(disk *v1beta1.BootDiskParameters) instance.VMOS {
	if disk == nil {
		return instance.VMOS{}
	}
	return instance.VMOS{
		ID:   *disk.InitializeParams[0].Image,
		Name: *disk.InitializeParams[0].Image,
	}
}

func (gcp *gcpRepository) toVMDiskCollection(disks []v1beta1.BootDiskParameters) []instance.VMDisk {
	var ret []instance.VMDisk
	for _, disk := range disks {
		ret = append(ret, gcp.toVMDisk(&disk))
	}
	return ret
}

func (gcp *gcpRepository) toVMDisk(disk *v1beta1.BootDiskParameters) instance.VMDisk {
	return instance.VMDisk{
		AutoDelete: *disk.AutoDelete,
		SizeGib:    int(*disk.InitializeParams[0].Size),
		Type:       fromGCPDiskType(*disk.InitializeParams[0].Type),
		Mode:       fromGCPDiskMode(*disk.Mode),
	}
}

func (gcp *gcpRepository) toGCPVMDiskList(disks []instance.VMDisk, os instance.VMOS) []v1beta1.BootDiskParameters {
	var bootDisks []v1beta1.BootDiskParameters
	for _, disk := range disks {
		bootDisks = append(bootDisks, gcp.toGCPVMDisk(disk, os))
	}
	return bootDisks
}

func (gcp *gcpRepository) toGCPVMDisk(disk instance.VMDisk, os instance.VMOS) v1beta1.BootDiskParameters {
	diskSize := float64(disk.SizeGib)
	diskType := toGCPDiskType(disk.Type)
	diskMode := toGCPDiskMode(disk.Mode)
	return v1beta1.BootDiskParameters{
		AutoDelete: &disk.AutoDelete,
		Mode:       &diskMode,
		InitializeParams: []v1beta1.InitializeParamsParameters{
			{
				Size:  &diskSize,
				Type:  &diskType,
				Image: &os.ID,
			},
		},
	}
}

func sshKeysToString(sshKeys model.SSHKeyList) *string {
	var ret string
	for _, key := range sshKeys {
		ret += fmt.Sprintf("%s:%s\n", key.Username, key.PublicKey)
	}
	ret = strings.TrimSuffix(ret, "\n")
	return &ret
}
