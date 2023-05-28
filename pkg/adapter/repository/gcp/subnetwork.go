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
	"github.com/samber/lo"
	"github.com/upbound/provider-gcp/apis/compute/v1beta1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sync"
)

func (gcp *gcpRepository) FindSubnetwork(ctx context.Context, opt option.Option) (*resource.Subnetwork, errors.Error) {
	if !opt.SetType(reflect.TypeOf(resourceRepo.FindResourceOption{}).String()).Validate() {
		return nil, errors.InvalidOption.WithMessage(fmt.Sprintf("invalid option : want %s, got %+v", reflect.TypeOf(resourceRepo.FindResourceOption{}).String(), opt.Get()))
	}
	req := opt.Get().(resourceRepo.FindResourceOption)
	gcpSubnetwork := &v1beta1.Subnetwork{}
	if err := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, gcpSubnetwork); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, errors.NotFound.WithMessage(fmt.Sprintf("subnetwork %s not found in namespace %s", req.Name, req.Namespace))
		}
		return nil, errors.KubernetesError.WithMessage(fmt.Sprintf("unable to get subnetwork %s in namespace %s", req.Name, req.Namespace))
	}
	mod, err := gcp.toModelSubnetwork(gcpSubnetwork)
	if !err.IsOk() {
		return nil, err
	}
	return mod, errors.OK

}

func (gcp *gcpRepository) FindAllSubnetworks(ctx context.Context, opt option.Option) (*resource.SubnetworkCollection, errors.Error) {
	panic("implement me")
}

func (gcp *gcpRepository) FindAllRecursiveSubnetworks(ctx context.Context, opt option.Option, ch *resourceRepo.SubnetworkChannel) {
	if !opt.SetType(reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String()).Validate() {
		ch.ErrorChannel <- errors.InvalidOption.WithMessage(fmt.Sprintf("invalid option : want %s, got %+v", reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String(), opt.Get()))
		return
	}
	req := opt.Get().(resourceRepo.FindAllResourceOption)
	gcpSubnetworkList := &v1beta1.SubnetworkList{}
	listOpt := &client.ListOptions{
		Namespace:     req.Namespace,
		LabelSelector: client.MatchingLabelsSelector{Selector: labels.SelectorFromSet(req.Labels)},
	}
	if err := kubernetes.Client().Client.List(ctx, gcpSubnetworkList, listOpt); err != nil {
		ch.ErrorChannel <- errors.KubernetesError.WithMessage(fmt.Sprintf("unable to get subnetworks in namespace %s", req.Namespace))
		return
	}
	subnetworkCollection, err := gcp.toModelSubnetworkCollection(gcpSubnetworkList)
	if !err.IsOk() {
		ch.ErrorChannel <- err
		return
	}
	wg := &sync.WaitGroup{}
	vmChannels := make([]*resourceRepo.VMChannel, 0)
	for _, subnet := range *subnetworkCollection {
		wg.Add(1)
		optVM := option.Option{}
		vmCh := &resourceRepo.VMChannel{
			WaitGroup:    ch.WaitGroup,
			Channel:      make(chan *resource.VMCollection),
			ErrorChannel: make(chan errors.Error),
		}
		vmChannels = append(vmChannels, vmCh)
		go gcp.FindAllRecursiveVMs(ctx, optVM, vmCh)
		select {
		case errCh := <-vmCh.ErrorChannel:
			ch.ErrorChannel <- errCh
		case vms := <-vmCh.Channel:
			subnet.VMs = *vms
		}
	}

	subnetworks, err := gcp.toModelSubnetworkCollection(gcpSubnetworkList)
	go func() {
		wg.Wait()
		for _, c := range vmChannels {
			close(c.Channel)
			close(c.ErrorChannel)
		}
	}()
	if !err.IsOk() {
		ch.ErrorChannel <- err
	} else {
		ch.Channel <- subnetworks
	}
}

func (gcp *gcpRepository) CreateSubnetwork(ctx context.Context, subnetwork *resource.Subnetwork) errors.Error {
	gcpSubnetwork := gcp.toGCPSubnetwork(ctx, subnetwork)
	if err := kubernetes.Client().Client.Create(ctx, gcpSubnetwork); err != nil {
		if k8serrors.IsAlreadyExists(err) {
			return errors.Conflict.WithMessage(fmt.Sprintf("subnetwork %s already exists in namespace %s", subnetwork.IdentifierName.Subnetwork, subnetwork.Metadata.Namespace))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to create subnetwork %s in namespace %s", subnetwork.IdentifierName.Subnetwork, subnetwork.Metadata.Namespace))
	}
	return errors.Created
}

func (gcp *gcpRepository) UpdateSubnetwork(ctx context.Context, subnetwork *resource.Subnetwork) errors.Error {
	gcpSubnetwork := gcp.toGCPSubnetwork(ctx, subnetwork)
	if err := kubernetes.Client().Client.Update(ctx, gcpSubnetwork); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("subnetwork %s not found in namespace %s", subnetwork.IdentifierName.Subnetwork, subnetwork.Metadata.Namespace))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to update subnetwork %s in namespace %s", subnetwork.IdentifierName.Subnetwork, subnetwork.Metadata.Namespace))
	}
	return errors.NoContent
}

func (gcp *gcpRepository) DeleteSubnetwork(ctx context.Context, subnetwork *resource.Subnetwork) errors.Error {
	gcpSubnetwork := gcp.toGCPSubnetwork(ctx, subnetwork)
	if err := kubernetes.Client().Client.Delete(ctx, gcpSubnetwork); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("subnetwork %s not found in namespace %s", subnetwork.IdentifierName.Subnetwork, subnetwork.Metadata.Namespace))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to delete subnetwork %s in namespace %s", subnetwork.IdentifierName.Subnetwork, subnetwork.Metadata.Namespace))
	}
	return errors.NoContent
}

func (gcp *gcpRepository) DeleteSubnetworkCascade(ctx context.Context, subnetwork *resource.Subnetwork) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (gcp *gcpRepository) toModelSubnetwork(subnet *v1beta1.Subnetwork) (*resource.Subnetwork, errors.Error) {
	id := identifier.Subnetwork{}
	name := identifier.Subnetwork{}
	if err := id.IDFromLabels(subnet.Labels); !err.IsOk() {
		return nil, err
	}
	if err := name.NameFromLabels(subnet.Labels); !err.IsOk() {
		return nil, err
	}
	return &resource.Subnetwork{
		Metadata: metadata.Metadata{
			Managed:   subnet.Spec.ResourceSpec.DeletionPolicy == v1.DeletionDelete,
			Namespace: subnet.ObjectMeta.Namespace,
		},
		IdentifierID:   id,
		IdentifierName: name,
		Region:         *subnet.Spec.ForProvider.Region,
		IPCIDRRange:    *subnet.Spec.ForProvider.IPCidrRange,
	}, errors.OK
}

func (gcp *gcpRepository) toGCPSubnetwork(ctx context.Context, subnet *resource.Subnetwork) *v1beta1.Subnetwork {
	resLabels := lo.Assign(crossplane.GetBaseLabels(ctx.Value(context.ProjectIDKey).(string)), subnet.IdentifierID.ToIDLabels(), subnet.IdentifierName.ToNameLabels())
	return &v1beta1.Subnetwork{
		ObjectMeta: metav1.ObjectMeta{
			Name:        subnet.IdentifierID.Network,
			Namespace:   subnet.Metadata.Namespace,
			Labels:      resLabels,
			Annotations: crossplane.GetAnnotations(subnet.Metadata.Managed, subnet.IdentifierName.Subnetwork),
		},
		Spec: v1beta1.SubnetworkSpec{
			ResourceSpec: v1.ResourceSpec{
				DeletionPolicy: crossplane.GetDeletionPolicy(subnet.Metadata.Managed),
				ProviderConfigReference: &v1.Reference{
					Name: subnet.IdentifierID.Provider,
				},
			},
			ForProvider: v1beta1.SubnetworkParameters_2{
				Network:     &subnet.IdentifierName.Network,
				Project:     &subnet.IdentifierName.VPC,
				Region:      &subnet.Region,
				IPCidrRange: &subnet.IPCIDRRange,
			},
		},
	}
}

func (gcp *gcpRepository) toModelSubnetworkCollection(subnetworkList *v1beta1.SubnetworkList) (*resource.SubnetworkCollection, errors.Error) {
	var items resource.SubnetworkCollection
	for _, item := range subnetworkList.Items {
		subnet, err := gcp.toModelSubnetwork(&item)
		if !err.IsOk() {
			return nil, err
		}
		items[item.ObjectMeta.Annotations[crossplane.ExternalNameAnnotationKey]] = *subnet
	}
	return &items, errors.OK
}
