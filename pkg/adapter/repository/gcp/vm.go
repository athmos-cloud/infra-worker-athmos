package gcp

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/crossplane"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	resourceRepo "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository/resource"
	v1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/upbound/provider-gcp/apis/compute/v1beta1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

const (
	tagKeyValueSeparator = "="
)

func (gcp *gcpRepository) FindVM(ctx context.Context, opt option.Option) (*resource.VM, errors.Error) {
	if !opt.SetType(reflect.TypeOf(resourceRepo.FindResourceOption{}).String()).Validate() {
		return nil, errors.InvalidOption.WithMessage(fmt.Sprintf("invalid option : want %s, got %+v", reflect.TypeOf(resourceRepo.FindResourceOption{}).String(), opt.Get()))
	}
	req := opt.Get().(resourceRepo.FindResourceOption)
	gcpVM := &v1beta1.Instance{}
	if err := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, gcpVM); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, errors.NotFound.WithMessage(fmt.Sprintf("vm %s not found in namespace %s", req.Name, req.Namespace))
		}
		return nil, errors.KubernetesError.WithMessage(fmt.Sprintf("unable to get vm %s in namespace %s", req.Name, req.Namespace))
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
		Namespace:     req.Namespace,
		LabelSelector: client.MatchingLabelsSelector{Selector: labels.SelectorFromSet(req.Labels)},
	}
	if err := kubernetes.Client().Client.List(ctx, gcpVMList, listOpt); err != nil {
		ch.ErrorChannel <- errors.KubernetesError.WithMessage(fmt.Sprintf("unable to list vm in namespace %s", req.Namespace))
		return
	}
	if firewalls, err := gcp.toModelVMCollection(gcpVMList); !err.IsOk() {
		ch.ErrorChannel <- err
	} else {
		ch.Channel <- firewalls
	}
}

func (gcp *gcpRepository) FindAllVMs(ctx context.Context, opt option.Option) (*resource.VMCollection, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (gcp *gcpRepository) CreateVM(ctx context.Context, vm *resource.VM) errors.Error {
	gcpVM := gcp.toGCPVM(ctx, vm)
	if err := kubernetes.Client().Client.Create(ctx, gcpVM); err != nil {
		if k8serrors.IsAlreadyExists(err) {
			return errors.Conflict.WithMessage(fmt.Sprintf("vm %s already exists in namespace %s", vm.IdentifierName.VM, vm.Metadata.Namespace))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to create vm %s in namespace %s", vm.IdentifierName.VM, vm.Metadata.Namespace))
	}
	return errors.Created
}

func (gcp *gcpRepository) UpdateVM(ctx context.Context, vm *resource.VM) errors.Error {
	gcpVM := gcp.toGCPVM(ctx, vm)
	if err := kubernetes.Client().Client.Update(ctx, gcpVM); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("vm %s not found in namespace %s", vm.IdentifierName.VM, vm.Metadata.Namespace))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to update vm %s in namespace %s", vm.IdentifierName.VM, vm.Metadata.Namespace))
	}
	return errors.NoContent
}

func (gcp *gcpRepository) DeleteVM(ctx context.Context, vm *resource.VM) errors.Error {
	gcpVM := gcp.toGCPVM(ctx, vm)
	if err := kubernetes.Client().Client.Delete(ctx, gcpVM); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("vm %s not found in namespace %s", vm.IdentifierName.VM, vm.Metadata.Namespace))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to delete vm %s in namespace %s", vm.IdentifierName.VM, vm.Metadata.Namespace))
	}
	return errors.NoContent
}

func (gcp *gcpRepository) toModelVM(vm *v1beta1.Instance) (*resource.VM, errors.Error) {
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
	return &resource.VM{
		Metadata: metadata.Metadata{
			Managed:   vm.Spec.ResourceSpec.DeletionPolicy == v1.DeletionDelete,
			Namespace: vm.ObjectMeta.Namespace,
			Tags:      tags,
		},
		IdentifierID:   id,
		IdentifierName: name,
		AssignPublicIP: false,
		PublicIP:       "",
		Zone:           *vm.Spec.ForProvider.Zone,
		MachineType:    *vm.Spec.ForProvider.MachineType,
		Auths:          resource.VMAuthList{},
		Disks:          resource.VMDiskList{},
		OS:             resource.VMOS{},
	}, errors.OK
}

func (gcp *gcpRepository) toGCPVM(ctx context.Context, vm *resource.VM) *v1beta1.Instance {
	//TODO implement me
	panic("implement me")
}

func (gcp *gcpRepository) toModelVMCollection(firewallList *v1beta1.InstanceList) (*resource.VMCollection, errors.Error) {
	var items resource.VMCollection
	for _, item := range firewallList.Items {
		vm, err := gcp.toModelVM(&item)
		if !err.IsOk() {
			return nil, err
		}
		items[item.ObjectMeta.Annotations[crossplane.ExternalNameAnnotationKey]] = *vm
	}
	return &items, errors.OK
}
