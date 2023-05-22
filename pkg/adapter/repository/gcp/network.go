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
	if err := kubernetes.Client().Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, gcpNetwork); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, errors.NotFound.WithMessage(fmt.Sprintf("network %s not found in namespace %s", req.Name, req.Namespace))
		}
		return nil, errors.KubernetesError.WithMessage(fmt.Sprintf("unable to get network %s in namespace %s", req.Name, req.Namespace))
	}
	mod, err := gcp.toModelNetwork(gcpNetwork)
	if !err.IsOk() {
		return nil, err
	}
	return mod, errors.OK
}

func (gcp *gcpRepository) FindAllNetworks(ctx context.Context, opt option.Option) (*resource.NetworkCollection, errors.Error) {
	//TODO implement me
	panic("implement me")
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
	if err := kubernetes.Client().List(ctx, gcpNetworkList, kubeOptions); err != nil {
		return nil, errors.KubernetesError.WithMessage(fmt.Sprintf("unable to list networks in namespace %s", req.Namespace))
	}
	modNetworks, err := gcp.toModelNetworkCollection(gcpNetworkList)
	if !err.IsOk() {
		return nil, err
	}
	wg := &sync.WaitGroup{}

	subnetChannels := make([]resourceRepo.SubnetworkChannel, 0)
	subnetCollection := resource.SubnetworkCollection{}
	firewallChannels := make([]resourceRepo.FirewallChannel, 0)
	firewallCollections := resource.FirewallCollection{}

	for _, network := range *modNetworks {
		subnetOpt := resourceRepo.FindAllResourceOption{
			Namespace: req.Namespace,
			Labels: map[string]string{
				identifier.ProviderLabelKey:     network.IdentifierID.Provider,
				identifier.NetworkIdentifierKey: network.IdentifierName.Network,
			},
		}
		chFirewall := &resourceRepo.FirewallChannel{
			WaitGroup:    wg,
			Channel:      make(chan *resource.Firewall),
			ErrorChannel: make(chan errors.Error),
		}
		chSubnet := &resourceRepo.SubnetworkChannel{
			WaitGroup:    wg,
			Channel:      make(chan *resource.Subnetwork),
			ErrorChannel: make(chan errors.Error),
		}
		subnetChannels = append(subnetChannels, *chSubnet)
		firewallChannels = append(firewallChannels, *chFirewall)

		wg.Add(2)
		go gcp.FindAllRecursiveFirewalls(ctx, option.Option{Value: subnetOpt}, chFirewall)
		go gcp.FindAllRecursiveSubnetworks(ctx, option.Option{Value: subnetOpt}, chSubnet)

		select {
		case firewall := <-chFirewall.Channel:
			firewallCollections[firewall.IdentifierName.Firewall] = *firewall
		case errChFirewall := <-chFirewall.ErrorChannel:
			logger.Error.Println("error while listing firewalls", errChFirewall)
		case subnetwork := <-chSubnet.Channel:
			subnetCollection[subnetwork.IdentifierName.Subnetwork] = *subnetwork
		case errCh := <-chSubnet.ErrorChannel:
			logger.Error.Println("error while listing subnetworks", errCh)
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
	gcpNetwork := gcp.toGCPNetwork(ctx, network)
	if err := kubernetes.Client().Create(ctx, gcpNetwork); err != nil {
		if k8serrors.IsAlreadyExists(err) {
			return errors.Conflict.WithMessage(fmt.Sprintf("subnetwork %s already exists in namespace %s", network.IdentifierName.Network, network.Metadata.Namespace))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to create subnetwork %s in namespace %s", network.IdentifierName.Network, network.Metadata.Namespace))
	}
	return errors.Created
}

func (gcp *gcpRepository) UpdateNetwork(ctx context.Context, network *resource.Network) errors.Error {
	gcpNetwork := gcp.toGCPNetwork(ctx, network)
	if err := kubernetes.Client().Update(ctx, gcpNetwork); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("subnetwork %s not found in namespace %s", network.IdentifierName.Network, network.Metadata.Namespace))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to update subnetwork %s in namespace %s", network.IdentifierName.Network, network.Metadata.Namespace))
	}
	return errors.NoContent
}

func (gcp *gcpRepository) DeleteNetwork(ctx context.Context, network *resource.Network) errors.Error {
	gcpSubnetwork := gcp.toGCPNetwork(ctx, network)
	if err := kubernetes.Client().Delete(ctx, gcpSubnetwork); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("subnetwork %s not found in namespace %s", network.IdentifierName.Network, network.Metadata.Namespace))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to delete subnetwork %s in namespace %s", network.IdentifierName.Network, network.Metadata.Namespace))
	}
	return errors.NoContent
}

func (gcp *gcpRepository) DeleteNetworkCascade(ctx context.Context, network *resource.Network) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (gcp *gcpRepository) toModelNetwork(network *v1beta1.Network) (*resource.Network, errors.Error) {
	id := identifier.Network{}
	if err := id.FromLabels(network.Labels); !err.IsOk() {
		return nil, err
	}
	name, ok := network.Annotations[crossplane.ExternalNameAnnotationKey]
	if !ok {
		return nil, errors.InternalError.WithMessage("cannot find external name annotation")
	}
	return &resource.Network{
		Metadata: metadata.Metadata{
			Managed:   network.Spec.ResourceSpec.DeletionPolicy == v1.DeletionDelete,
			Namespace: network.ObjectMeta.Namespace,
		},
		IdentifierID: id,
		IdentifierName: identifier.Network{
			Network:  name,
			VPC:      *network.Spec.ForProvider.Project,
			Provider: network.Spec.ProviderConfigReference.Name,
		},
	}, errors.OK
}

func (gcp *gcpRepository) toGCPNetwork(ctx context.Context, network *resource.Network) *v1beta1.Network {
	return &v1beta1.Network{
		ObjectMeta: metav1.ObjectMeta{
			Name:        network.IdentifierID.Network,
			Namespace:   network.Metadata.Namespace,
			Labels:      lo.Assign(crossplane.GetBaseLabels(ctx.Value(context.ProjectIDKey).(string)), network.IdentifierID.ToLabels()),
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
