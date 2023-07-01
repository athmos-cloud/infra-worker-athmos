package gcp

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/crossplane"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/instance"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/network"
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
)

func (gcp *gcpRepository) FindSubnetwork(ctx context.Context, opt option.Option) (*network.Subnetwork, errors.Error) {
	if !opt.SetType(reflect.TypeOf(resourceRepo.FindResourceOption{}).String()).Validate() {
		return nil, errors.InvalidOption.WithMessage(fmt.Sprintf("invalid option : want %s, got %+v", reflect.TypeOf(resourceRepo.FindResourceOption{}).String(), opt.Get()))
	}
	req := opt.Get().(resourceRepo.FindResourceOption)
	gcpSubnetwork := &v1beta1.Subnetwork{}
	if err := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: req.Name}, gcpSubnetwork); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, errors.NotFound.WithMessage(fmt.Sprintf("subnetwork %s not found", req.Name))
		}
		return nil, errors.KubernetesError.WithMessage(fmt.Sprintf("unable to get subnetwork %s", req.Name))
	}
	mod, err := gcp.toModelSubnetwork(gcpSubnetwork)
	if !err.IsOk() {
		return nil, err
	}
	return mod, errors.OK

}

func (gcp *gcpRepository) FindAllSubnetworks(ctx context.Context, opt option.Option) (*network.SubnetworkCollection, errors.Error) {
	if !opt.SetType(reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String()).Validate() {
		return nil, errors.InvalidOption.WithMessage(fmt.Sprintf("invalid option : want %s, got %+v", reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String(), opt.Get()))
	}
	req := opt.Get().(resourceRepo.FindAllResourceOption)
	gcpSubnetworkList := &v1beta1.SubnetworkList{}
	listOpt := &client.ListOptions{
		LabelSelector: client.MatchingLabelsSelector{Selector: labels.SelectorFromSet(req.Labels)},
	}
	if err := kubernetes.Client().Client.List(ctx, gcpSubnetworkList, listOpt); err != nil {
		return nil, errors.KubernetesError.WithMessage(fmt.Sprintf("unable to get subnetworks"))
	}
	subnetworkCollection, err := gcp.toModelSubnetworkCollection(gcpSubnetworkList)
	if !err.IsOk() {
		return nil, err
	}
	return subnetworkCollection, errors.OK
}

func (gcp *gcpRepository) FindAllRecursiveSubnetworks(ctx context.Context, opt option.Option, ch *resourceRepo.SubnetworkChannel) {
	if !opt.SetType(reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String()).Validate() {
		ch.ErrorChannel <- errors.InvalidOption.WithMessage(fmt.Sprintf("invalid option : want %s, got %+v", reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String(), opt.Get()))
		return
	}
	req := opt.Get().(resourceRepo.FindAllResourceOption)
	gcpSubnetworkList := &v1beta1.SubnetworkList{}
	listOpt := &client.ListOptions{
		LabelSelector: client.MatchingLabelsSelector{Selector: labels.SelectorFromSet(req.Labels)},
	}
	if err := kubernetes.Client().Client.List(ctx, gcpSubnetworkList, listOpt); err != nil {
		ch.ErrorChannel <- errors.KubernetesError.WithMessage("unable to get subnetworks")
		return
	}
	subnetworkCollection, err := gcp.toModelSubnetworkCollection(gcpSubnetworkList)
	if !err.IsOk() {
		ch.ErrorChannel <- err
		return
	}
	vmChannels := make([]*resourceRepo.VMChannel, 0)
	subnetResult := &network.SubnetworkCollection{}
	for _, subnet := range *subnetworkCollection {
		optVM := option.Option{Value: resourceRepo.FindAllResourceOption{
			Labels: subnet.IdentifierID.ToIDLabels(),
		}}
		vmCh := &resourceRepo.VMChannel{
			Channel:      make(chan *instance.VMCollection),
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
		(*subnetResult)[subnet.IdentifierName.Subnetwork] = subnet
	}
	for _, c := range vmChannels {
		close(c.Channel)
		close(c.ErrorChannel)
	}

	if !err.IsOk() {
		ch.ErrorChannel <- err
	} else {
		ch.Channel <- subnetResult
	}
}

func (gcp *gcpRepository) CreateSubnetwork(ctx context.Context, subnetwork *network.Subnetwork) errors.Error {
	if exists, err := gcp.SubnetworkExists(ctx, subnetwork); !err.IsOk() {
		return err
	} else if exists {
		return errors.Conflict.WithMessage(fmt.Sprintf("subnetwork %s already exists", subnetwork.IdentifierName.Subnetwork))
	}
	gcpSubnetwork := gcp.toGCPSubnetwork(ctx, subnetwork)
	if err := kubernetes.Client().Client.Create(ctx, gcpSubnetwork); err != nil {
		if k8serrors.IsAlreadyExists(err) {
			return errors.Conflict.WithMessage(fmt.Sprintf("subnetwork %s already exists", subnetwork.IdentifierName.Subnetwork))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to create subnetwork %s", subnetwork.IdentifierName.Subnetwork))
	}
	return errors.Created
}

func (gcp *gcpRepository) UpdateSubnetwork(ctx context.Context, subnetwork *network.Subnetwork) errors.Error {
	existingSubnet := &v1beta1.Subnetwork{}
	if err := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: subnetwork.IdentifierID.Subnetwork}, existingSubnet); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("subnetwork %s not found", subnetwork.IdentifierID.Subnetwork))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to get subnetwork %s", subnetwork.IdentifierID.Subnetwork))
	}
	gcpSubnetwork := gcp.toGCPSubnetwork(ctx, subnetwork)
	existingSubnet.Spec = gcpSubnetwork.Spec
	existingSubnet.Labels = gcpSubnetwork.Labels

	if err := kubernetes.Client().Client.Update(ctx, existingSubnet); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("subnetwork %s not found in namespace %s", subnetwork.IdentifierName.Subnetwork))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to update subnetwork %s in namespace %s", subnetwork.IdentifierName.Subnetwork))
	}
	return errors.NoContent
}

func (gcp *gcpRepository) DeleteSubnetwork(ctx context.Context, subnetwork *network.Subnetwork) errors.Error {
	// gcpSubnetwork := gcp.toGCPSubnetwork(ctx, subnetwork)
	gcpSubnetwork := &v1beta1.Subnetwork{}
	if err := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: subnetwork.IdentifierID.Subnetwork}, existingSubnet); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("subnetwork %s not found", subnetwork.IdentifierID.Subnetwork))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to get subnetwork %s", subnetwork.IdentifierID.Subnetwork))
	}

	searchLabels := lo.Assign(map[string]string{model.ProjectIDLabelKey: ctx.Value(context.ProjectIDKey).(string)}, subnetwork.IdentifierID.ToIDLabels())
	vms, errVMs := gcp.FindAllVMs(ctx, option.Option{Value: resourceRepo.FindAllResourceOption{Labels: searchLabels}})
	if !errVMs.IsOk() {
		return errVMs
	}
	if len(*vms) > 0 {
		return errors.Conflict.WithMessage(fmt.Sprintf("subnetwork %s still has VMs", subnetwork.IdentifierName.Subnetwork))
	}
	if err := kubernetes.Client().Client.Delete(ctx, gcpSubnetwork); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("subnetwork %s not found in namespace %s", subnetwork.IdentifierName.Subnetwork))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to delete subnetwork %s in namespace %s", subnetwork.IdentifierName.Subnetwork))
	}

	return errors.NoContent
}

func (gcp *gcpRepository) DeleteSubnetworkCascade(ctx context.Context, subnetwork *network.Subnetwork) errors.Error {
	searchLabels := lo.Assign(map[string]string{model.ProjectIDLabelKey: ctx.Value(context.ProjectIDKey).(string)}, subnetwork.IdentifierID.ToIDLabels())
	if vms, err := gcp.FindAllVMs(ctx, option.Option{Value: resourceRepo.FindAllResourceOption{Labels: searchLabels}}); !err.IsOk() {
		return err
	} else {
		for _, vm := range *vms {
			if vmErr := gcp.DeleteVM(ctx, &vm); !err.IsOk() {
				return vmErr
			}
		}
		gcpSubnetwork := &v1beta1.Subnetwork{}
		if err := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: subnetwork.IdentifierID.Subnetwork}, existingSubnet); err != nil {
			if k8serrors.IsNotFound(err) {
				return errors.NotFound.WithMessage(fmt.Sprintf("subnetwork %s not found", subnetwork.IdentifierID.Subnetwork))
			}
			return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to get subnetwork %s", subnetwork.IdentifierID.Subnetwork))
		}
		if err := kubernetes.Client().Client.Delete(ctx, gcpSubnetwork); err != nil {
			if k8serrors.IsNotFound(err) {
				return errors.NotFound.WithMessage(fmt.Sprintf("subnetwork %s not found in namespace %s", subnetwork.IdentifierName.Subnetwork))
			}
			return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to delete subnetwork %s in namespace %s", subnetwork.IdentifierName.Subnetwork))
		}
		return errors.NoContent
	}
}

func (gcp *gcpRepository) SubnetworkExists(ctx context.Context, subnetwork *network.Subnetwork) (bool, errors.Error) {
	gcpSubnetwork := &v1beta1.SubnetworkList{}
	searchLabels := lo.Assign(map[string]string{
		model.ProjectIDLabelKey:          ctx.Value(context.ProjectIDKey).(string),
		identifier.ProviderIdentifierKey: subnetwork.IdentifierID.Provider,
		identifier.NetworkIdentifierKey:  subnetwork.IdentifierID.Network,
		identifier.SubnetworkNameKey:     subnetwork.IdentifierName.Subnetwork,
	})
	listOpt := &client.ListOptions{
		LabelSelector: client.MatchingLabelsSelector{Selector: labels.SelectorFromSet(searchLabels)},
	}
	if err := kubernetes.Client().Client.List(ctx, gcpSubnetwork, listOpt); err != nil {
		return false, errors.KubernetesError.WithMessage(fmt.Sprintf("unable to get subnetworks"))
	}
	return len(gcpSubnetwork.Items) > 0, errors.OK
}

func (gcp *gcpRepository) toModelSubnetwork(subnet *v1beta1.Subnetwork) (*network.Subnetwork, errors.Error) {
	id := identifier.Subnetwork{}
	name := identifier.Subnetwork{}
	if err := id.IDFromLabels(subnet.Labels); !err.IsOk() {
		return nil, err
	}
	if err := name.NameFromLabels(subnet.Labels); !err.IsOk() {
		return nil, err
	}
	return &network.Subnetwork{
		Metadata: metadata.Metadata{
			Status:  metadata.StatusFromKubernetesStatus(subnet.Status.Conditions),
			Managed: subnet.Spec.ResourceSpec.DeletionPolicy == v1.DeletionDelete,
		},
		IdentifierID:   id,
		IdentifierName: name,
		Region:         *subnet.Spec.ForProvider.Region,
		IPCIDRRange:    *subnet.Spec.ForProvider.IPCidrRange,
	}, errors.OK
}

func (gcp *gcpRepository) toGCPSubnetwork(ctx context.Context, subnet *network.Subnetwork) *v1beta1.Subnetwork {
	resLabels := lo.Assign(crossplane.GetBaseLabels(ctx.Value(context.ProjectIDKey).(string)), subnet.IdentifierID.ToIDLabels(), subnet.IdentifierName.ToNameLabels())
	return &v1beta1.Subnetwork{
		ObjectMeta: metav1.ObjectMeta{
			Name:        subnet.IdentifierID.Subnetwork,
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
				Project:     &subnet.IdentifierID.VPC,
				Region:      &subnet.Region,
				IPCidrRange: &subnet.IPCIDRRange,
			},
		},
	}
}

func (gcp *gcpRepository) toModelSubnetworkCollection(subnetworkList *v1beta1.SubnetworkList) (*network.SubnetworkCollection, errors.Error) {
	items := network.SubnetworkCollection{}
	for _, item := range subnetworkList.Items {
		subnet, err := gcp.toModelSubnetwork(&item)
		if !err.IsOk() {
			return nil, err
		}
		items[item.ObjectMeta.Annotations[crossplane.ExternalNameAnnotationKey]] = *subnet
	}
	return &items, errors.OK
}
