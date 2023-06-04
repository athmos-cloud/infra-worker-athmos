package gcp

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/crossplane"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/logger"
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

var (
	autoCreateSubnetworks = false
)

func (gcp *gcpRepository) FindNetwork(ctx context.Context, opt option.Option) (*resource.Network, errors.Error) {
	if !opt.SetType(reflect.TypeOf(resourceRepo.FindResourceOption{}).String()).Validate() {
		return nil, errors.InvalidOption.WithMessage(fmt.Sprintf("invalid option : want %s, got %+v", reflect.TypeOf(resourceRepo.FindResourceOption{}).String(), opt.Get()))
	}
	req := opt.Get().(resourceRepo.FindResourceOption)
	gcpNetwork := &v1beta1.Network{}
	if err := kubernetes.Client().Client.Get(ctx, client.ObjectKey{Name: req.Name, Namespace: req.Namespace}, gcpNetwork); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, errors.NotFound.WithMessage(fmt.Sprintf("network %s not found in namespace %s", req.Name, req.Namespace))
		}
		return nil, errors.KubernetesError.WithMessage(fmt.Sprintf("unable to get network %s in namespace %s", req.Name, req.Namespace))
	}
	gcpNetwork.ObjectMeta.Namespace = req.Namespace // hack to avoid kubernetes client bug
	mod, err := gcp.toModelNetwork(gcpNetwork)
	if !err.IsOk() {
		return nil, err
	}
	return mod, errors.OK
}

func (gcp *gcpRepository) FindAllNetworks(ctx context.Context, opt option.Option) (*resource.NetworkCollection, errors.Error) {
	if !opt.SetType(reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String()).Validate() {
		return nil, errors.InvalidOption.WithMessage(fmt.Sprintf("invalid option : want %s, got %+v", reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String(), opt.Get()))
	}
	req := opt.Get().(resourceRepo.FindAllResourceOption)
	gcpNetworkList := &v1beta1.NetworkList{}
	kubeOptions := &client.ListOptions{
		Namespace:     req.Namespace,
		LabelSelector: client.MatchingLabelsSelector{Selector: labels.SelectorFromSet(req.Labels)},
	}
	if err := kubernetes.Client().Client.List(ctx, gcpNetworkList, kubeOptions); err != nil {
		return nil, errors.KubernetesError.WithMessage(fmt.Sprintf("unable to list networks in namespace %s", req.Namespace))
	}
	modNetworks, err := gcp.toModelNetworkCollection(gcpNetworkList)
	if !err.IsOk() {
		return nil, err
	}
	return modNetworks, errors.OK

}

func (gcp *gcpRepository) FindAllRecursiveNetworks(ctx context.Context, opt option.Option, _ *resourceRepo.NetworkChannel) (*resource.NetworkCollection, errors.Error) {
	if !opt.SetType(reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String()).Validate() {
		return nil, errors.InvalidOption.WithMessage(fmt.Sprintf("invalid option : want %s, got %+v", reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String(), opt.Get()))
	}
	req := opt.Get().(resourceRepo.FindAllResourceOption)
	gcpNetworkList := &v1beta1.NetworkList{}
	kubeOptions := &client.ListOptions{
		Namespace:     req.Namespace,
		LabelSelector: client.MatchingLabelsSelector{Selector: labels.SelectorFromSet(req.Labels)},
	}
	if err := kubernetes.Client().Client.List(ctx, gcpNetworkList, kubeOptions); err != nil {
		return nil, errors.KubernetesError.WithMessage(fmt.Sprintf("unable to list networks in namespace %s", req.Namespace))
	}
	modNetworks, err := gcp.toModelNetworkCollection(gcpNetworkList)
	if !err.IsOk() {
		return nil, err
	}
	wg := &sync.WaitGroup{}

	subnetChannels := make([]resourceRepo.SubnetworkChannel, 0)
	firewallChannels := make([]resourceRepo.FirewallChannel, 0)

	for _, network := range *modNetworks {
		chFirewall := &resourceRepo.FirewallChannel{
			WaitGroup:    wg,
			Channel:      make(chan *resource.FirewallCollection),
			ErrorChannel: make(chan errors.Error),
		}
		chSubnet := &resourceRepo.SubnetworkChannel{
			WaitGroup:    wg,
			Channel:      make(chan *resource.SubnetworkCollection),
			ErrorChannel: make(chan errors.Error),
		}
		subnetChannels = append(subnetChannels, *chSubnet)
		firewallChannels = append(firewallChannels, *chFirewall)

		wg.Add(2)
		go gcp.FindAllRecursiveFirewalls(ctx, option.Option{Value: resourceRepo.FindAllResourceOption{Labels: network.IdentifierID.ToIDLabels()}}, chFirewall)
		go gcp.FindAllRecursiveSubnetworks(ctx, option.Option{Value: resourceRepo.FindAllResourceOption{Labels: network.IdentifierID.ToIDLabels()}}, chSubnet)
		gotFirewalls := false
		gotSubnets := false
		for {
			select {
			case firewalls := <-chFirewall.Channel:
				network.Firewalls = *firewalls
				if gotSubnets {
					wg.Done()
					return modNetworks, errors.OK
				}
				gotFirewalls = true
			case errChFirewall := <-chFirewall.ErrorChannel:
				logger.Error.Println("error while listing firewalls", errChFirewall)
				if gotSubnets {
					wg.Done()
					return modNetworks, errors.OK
				}
				gotFirewalls = true
			case subnetworks := <-chSubnet.Channel:
				network.Subnetworks = *subnetworks
				if gotFirewalls {
					wg.Done()
					return modNetworks, errors.OK
				}
				gotSubnets = true
			case errCh := <-chSubnet.ErrorChannel:
				logger.Error.Println("error while listing subnetworks", errCh)
				if gotFirewalls {
					wg.Done()
					return modNetworks, errors.OK
				}
				gotFirewalls = true
			}
		}
	}
	go func() {
		wg.Wait()
		for _, ch := range subnetChannels {
			close(ch.Channel)
			close(ch.ErrorChannel)
		}
		for _, ch := range firewallChannels {
			close(ch.Channel)
			close(ch.ErrorChannel)
		}
	}()

	return nil, errors.OK
}

func (gcp *gcpRepository) CreateNetwork(ctx context.Context, network *resource.Network) errors.Error {
	if exist, errExists := gcp.NetworkExists(ctx, network); !errExists.IsOk() {
		return errExists
	} else if exist {
		return errors.Conflict.WithMessage(fmt.Sprintf("network %s already exists in namespace %s", network.IdentifierName.Network, network.Metadata.Namespace))
	}
	gcpNetwork := gcp.toGCPNetwork(ctx, network)
	if err := kubernetes.Client().Client.Create(ctx, gcpNetwork); err != nil {
		if k8serrors.IsAlreadyExists(err) {
			return errors.Conflict.WithMessage(fmt.Sprintf("subnetwork %s already exists in namespace %s", network.IdentifierName.Network, network.Metadata.Namespace))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to create subnetwork %s in namespace %s", network.IdentifierName.Network, network.Metadata.Namespace))
	}
	return errors.Created
}

func (gcp *gcpRepository) UpdateNetwork(ctx context.Context, network *resource.Network) errors.Error {
	existingNetwork := &v1beta1.Network{}
	if err := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: network.IdentifierID.Network, Namespace: network.Metadata.Namespace}, existingNetwork); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("network %s not found", network.IdentifierID.Network))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to get network %s", network.IdentifierID.Network))
	}
	gcpNetwork := gcp.toGCPNetwork(ctx, network)
	existingNetwork.Spec = gcpNetwork.Spec
	existingNetwork.Labels = gcpNetwork.Labels

	if err := kubernetes.Client().Client.Update(ctx, existingNetwork); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("subnetwork %s not found in namespace %s", network.IdentifierName.Network, network.Metadata.Namespace))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to update subnetwork %s", network.IdentifierName.Network))
	}
	return errors.NoContent
}

func (gcp *gcpRepository) DeleteNetwork(ctx context.Context, network *resource.Network) errors.Error {
	searchLabels := lo.Assign(map[string]string{model.ProjectIDLabelKey: ctx.Value(context.ProjectIDKey).(string)}, network.IdentifierID.ToIDLabels())
	if subnets, err := gcp.FindAllSubnetworks(ctx, option.Option{Value: resourceRepo.FindAllResourceOption{Labels: searchLabels}}); !err.IsOk() {
		return err
	} else if len(*subnets) > 0 {
		return errors.BadRequest.WithMessage(fmt.Sprintf("can't delete network %s without cascade option", network.IdentifierName.Network))
	}
	gcpSubnetwork := gcp.toGCPNetwork(ctx, network)
	if err := kubernetes.Client().Client.Delete(ctx, gcpSubnetwork); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("subnetwork %s not found in namespace %s", network.IdentifierName.Network, network.Metadata.Namespace))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to delete subnetwork %s in namespace %s", network.IdentifierName.Network, network.Metadata.Namespace))
	}
	return errors.NoContent
}

func (gcp *gcpRepository) DeleteNetworkCascade(ctx context.Context, network *resource.Network) errors.Error {
	searchLabels := lo.Assign(map[string]string{model.ProjectIDLabelKey: ctx.Value(context.ProjectIDKey).(string)}, network.IdentifierID.ToIDLabels())
	subnets, subnetsErr := gcp.FindAllSubnetworks(ctx, option.Option{Value: resourceRepo.FindAllResourceOption{Labels: searchLabels}})
	firewalls, firewallsErr := gcp.FindAllFirewalls(ctx, option.Option{Value: resourceRepo.FindAllResourceOption{Labels: searchLabels}})

	if !subnetsErr.IsOk() {
		return subnetsErr
	}
	if !firewallsErr.IsOk() {
		return firewallsErr
	}

	for _, subnet := range *subnets {
		if subnetErr := gcp.DeleteSubnetworkCascade(ctx, &subnet); !subnetErr.IsOk() {
			return subnetErr
		}
	}

	for _, firewall := range *firewalls {
		if firewallErr := gcp.DeleteFirewall(ctx, &firewall); !firewallErr.IsOk() {
			return firewallErr
		}
	}
	return gcp.DeleteNetwork(ctx, network)

}

func (gcp *gcpRepository) NetworkExists(ctx context.Context, network *resource.Network) (bool, errors.Error) {
	gcpNetwork := &v1beta1.NetworkList{}
	searchLabels := lo.Assign(map[string]string{
		model.ProjectIDLabelKey:          ctx.Value(context.ProjectIDKey).(string),
		identifier.ProviderIdentifierKey: network.IdentifierID.Provider,
		identifier.NetworkNameKey:        network.IdentifierName.Network,
	})
	kubeOptions := &client.ListOptions{
		LabelSelector: client.MatchingLabelsSelector{Selector: labels.SelectorFromSet(searchLabels)},
	}
	if err := kubernetes.Client().Client.List(ctx, gcpNetwork, kubeOptions); err != nil {
		return false, errors.KubernetesError.WithMessage("unable to list subnetworks in namespace")
	}
	return len(gcpNetwork.Items) > 0, errors.OK
}

func (gcp *gcpRepository) toModelNetwork(network *v1beta1.Network) (*resource.Network, errors.Error) {
	id := identifier.Network{}
	name := identifier.Network{}
	if err := id.IDFromLabels(network.Labels); !err.IsOk() {
		return nil, err
	}
	if err := name.NameFromLabels(network.Labels); !err.IsOk() {
		return nil, err
	}
	return &resource.Network{
		Metadata: metadata.Metadata{
			Managed:   network.Spec.ResourceSpec.DeletionPolicy == v1.DeletionDelete,
			Namespace: network.ObjectMeta.Namespace,
		},
		IdentifierID:   id,
		IdentifierName: name,
	}, errors.OK
}

func (gcp *gcpRepository) toGCPNetwork(ctx context.Context, network *resource.Network) *v1beta1.Network {
	resLabels := lo.Assign(crossplane.GetBaseLabels(ctx.Value(context.ProjectIDKey).(string)), network.IdentifierID.ToIDLabels(), network.IdentifierName.ToNameLabels())

	return &v1beta1.Network{
		ObjectMeta: metav1.ObjectMeta{
			Name:        network.IdentifierID.Network,
			Labels:      resLabels,
			Annotations: crossplane.GetAnnotations(network.Metadata.Managed, network.IdentifierName.Network),
		},
		Spec: v1beta1.NetworkSpec{
			ResourceSpec: v1.ResourceSpec{
				DeletionPolicy: crossplane.GetDeletionPolicy(network.Metadata.Managed),
				ProviderConfigReference: &v1.Reference{
					Name: network.IdentifierID.Provider,
				},
			},
			ForProvider: v1beta1.NetworkParameters{
				Project:               &network.IdentifierName.VPC,
				AutoCreateSubnetworks: &autoCreateSubnetworks,
			},
		},
	}
}

func (gcp *gcpRepository) toModelNetworkCollection(list *v1beta1.NetworkList) (*resource.NetworkCollection, errors.Error) {
	res := resource.NetworkCollection{}
	for _, item := range list.Items {
		network, err := gcp.toModelNetwork(&item)
		if !err.IsOk() {
			return &res, err
		}
		res[network.IdentifierName.Network] = *network
	}
	return &res, errors.OK
}
