package gcp

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/crossplane"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/instance"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/metadata"
	networkModels "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/network"
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
)

var (
	autoCreateSubnetworks = false
)

func (gcp *gcpRepository) FindNetwork(ctx context.Context, opt option.Option) (*networkModels.Network, errors.Error) {
	if !opt.SetType(reflect.TypeOf(resourceRepo.FindResourceOption{}).String()).Validate() {
		return nil, errors.InvalidOption.WithMessage(fmt.Sprintf("invalid option : want %s, got %+v", reflect.TypeOf(resourceRepo.FindResourceOption{}).String(), opt.Get()))
	}
	req := opt.Get().(resourceRepo.FindResourceOption)
	gcpNetwork := &v1beta1.Network{}
	if err := kubernetes.Client().Client.Get(ctx, client.ObjectKey{Name: req.Name}, gcpNetwork); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, errors.NotFound.WithMessage(fmt.Sprintf("networkModels %s not found", req.Name))
		}
		return nil, errors.KubernetesError.WithMessage(fmt.Sprintf("unable to get networkModels %s in namespace %s", req.Name))
	}
	mod, err := gcp.toModelNetwork(gcpNetwork)
	if !err.IsOk() {
		return nil, err
	}
	return mod, errors.OK
}

func (gcp *gcpRepository) FindAllNetworks(ctx context.Context, opt option.Option) (*networkModels.Collection, errors.Error) {
	if !opt.SetType(reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String()).Validate() {
		return nil, errors.InvalidOption.WithMessage(fmt.Sprintf("invalid option : want %s, got %+v", reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String(), opt.Get()))
	}
	req := opt.Get().(resourceRepo.FindAllResourceOption)
	gcpNetworkList := &v1beta1.NetworkList{}
	kubeOptions := &client.ListOptions{
		LabelSelector: client.MatchingLabelsSelector{Selector: labels.SelectorFromSet(req.Labels)},
	}
	if err := kubernetes.Client().Client.List(ctx, gcpNetworkList, kubeOptions); err != nil {
		return nil, errors.KubernetesError.WithMessage(fmt.Sprintf("unable to list networks"))
	}
	modNetworks, err := gcp.toModelNetworkCollection(gcpNetworkList)
	if !err.IsOk() {
		return nil, err
	}
	return modNetworks, errors.OK

}

func (gcp *gcpRepository) FindAllRecursiveNetworks(ctx context.Context, opt option.Option, _ *resourceRepo.NetworkChannel) (*networkModels.Collection, errors.Error) {
	if !opt.SetType(reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String()).Validate() {
		return nil, errors.InvalidOption.WithMessage(fmt.Sprintf("invalid option : want %s, got %+v", reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String(), opt.Get()))
	}
	req := opt.Get().(resourceRepo.FindAllResourceOption)
	gcpNetworkList := &v1beta1.NetworkList{}
	kubeOptions := &client.ListOptions{
		LabelSelector: client.MatchingLabelsSelector{Selector: labels.SelectorFromSet(req.Labels)},
	}
	if err := kubernetes.Client().Client.List(ctx, gcpNetworkList, kubeOptions); err != nil {
		return nil, errors.KubernetesError.WithMessage(fmt.Sprintf("unable to list networks"))
	}
	modNetworks, err := gcp.toModelNetworkCollection(gcpNetworkList)
	if !err.IsOk() {
		return nil, err
	}

	subnetChannels := make([]resourceRepo.SubnetworkChannel, 0)
	firewallChannels := make([]resourceRepo.FirewallChannel, 0)
	dbChannels := make([]resourceRepo.SqlDBChannel, 0)

	getNested := func(network *networkModels.Network) {
		chFirewall := &resourceRepo.FirewallChannel{
			Channel:      make(chan *networkModels.FirewallCollection),
			ErrorChannel: make(chan errors.Error),
		}
		chSubnet := &resourceRepo.SubnetworkChannel{
			Channel:      make(chan *networkModels.SubnetworkCollection),
			ErrorChannel: make(chan errors.Error),
		}
		chDB := &resourceRepo.SqlDBChannel{
			Channel:      make(chan *instance.SqlDBCollection),
			ErrorChannel: make(chan errors.Error),
		}
		subnetChannels = append(subnetChannels, *chSubnet)
		firewallChannels = append(firewallChannels, *chFirewall)
		dbChannels = append(dbChannels, *chDB)

		go gcp.FindAllRecursiveFirewalls(ctx, option.Option{Value: resourceRepo.FindAllResourceOption{Labels: network.IdentifierID.ToIDLabels()}}, chFirewall)
		go gcp.FindAllRecursiveSubnetworks(ctx, option.Option{Value: resourceRepo.FindAllResourceOption{Labels: network.IdentifierID.ToIDLabels()}}, chSubnet)
		go gcp.FindAllRecursiveSqlDBs(ctx, option.Option{Value: resourceRepo.FindAllResourceOption{Labels: network.IdentifierID.ToIDLabels()}}, chDB)

		gotFirewalls := false
		gotSubnets := false
		gotDBs := false
		for {
			select {
			case firewalls := <-chFirewall.Channel:
				network.Firewalls = *firewalls
				if gotSubnets {
					return
				}
				gotFirewalls = true
			case errCh := <-chFirewall.ErrorChannel:
				logger.Error.Println("error while listing firewalls", errCh)
				if gotSubnets {
					return
				}
				gotFirewalls = true
			case subnetworks := <-chSubnet.Channel:
				network.Subnetworks = *subnetworks
				if gotFirewalls {
					return
				}
				gotSubnets = true
			case errCh := <-chSubnet.ErrorChannel:
				logger.Error.Println("error while listing subnetworks", errCh)
				if gotFirewalls {
					return
				}
				gotFirewalls = true
			case dbs := <-chDB.Channel:
				network.SqlDbs = *dbs
				if gotDBs {
					return
				}
				gotDBs = true
			case errCh := <-chSubnet.ErrorChannel:
				logger.Error.Println("error while listing dbs", errCh)
				if gotDBs {
					return
				}
				gotDBs = true
			}
		}
	}
	networks := &networkModels.Collection{}
	for _, network := range *modNetworks {
		getNested(&network)
		(*networks)[network.IdentifierName.Network] = network
	}
	for _, ch := range subnetChannels {
		close(ch.Channel)
		close(ch.ErrorChannel)
	}
	for _, ch := range firewallChannels {
		close(ch.Channel)
		close(ch.ErrorChannel)
	}

	return networks, errors.OK
}

func (gcp *gcpRepository) CreateNetwork(ctx context.Context, network *networkModels.Network) errors.Error {
	if exist, errExists := gcp.NetworkExists(ctx, network); !errExists.IsOk() {
		return errExists
	} else if exist {
		return errors.Conflict.WithMessage(fmt.Sprintf("networkModels %s already exists", network.IdentifierName.Network))
	}
	gcpNetwork := gcp.toGCPNetwork(ctx, network)
	if err := kubernetes.Client().Client.Create(ctx, gcpNetwork); err != nil {
		if k8serrors.IsAlreadyExists(err) {
			return errors.Conflict.WithMessage(fmt.Sprintf("subnetwork %s already exists", network.IdentifierName.Network))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to create subnetwork %s", network.IdentifierName.Network))
	}
	return errors.Created
}

func (gcp *gcpRepository) UpdateNetwork(ctx context.Context, network *networkModels.Network) errors.Error {
	existingNetwork := &v1beta1.Network{}
	if err := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: network.IdentifierID.Network}, existingNetwork); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("networkModels %s not found", network.IdentifierID.Network))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to get networkModels %s", network.IdentifierID.Network))
	}
	gcpNetwork := gcp.toGCPNetwork(ctx, network)
	existingNetwork.Spec = gcpNetwork.Spec
	existingNetwork.Labels = gcpNetwork.Labels

	if err := kubernetes.Client().Client.Update(ctx, existingNetwork); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("subnetwork %s not found", network.IdentifierName.Network))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to update subnetwork %s", network.IdentifierName.Network))
	}
	return errors.NoContent
}

func (gcp *gcpRepository) DeleteNetwork(ctx context.Context, network *networkModels.Network) errors.Error {
	searchLabels := lo.Assign(map[string]string{model.ProjectIDLabelKey: ctx.Value(context.ProjectIDKey).(string)}, network.IdentifierID.ToIDLabels())
	if subnets, err := gcp.FindAllSubnetworks(ctx, option.Option{Value: resourceRepo.FindAllResourceOption{Labels: searchLabels}}); !err.IsOk() {
		return err
	} else if len(*subnets) > 0 {
		return errors.BadRequest.WithMessage(fmt.Sprintf("network %s still has subnetworks", network.IdentifierID.Network))
	}
	if dbs, err := gcp.FindAllSqlDBs(ctx, option.Option{Value: resourceRepo.FindAllResourceOption{Labels: searchLabels}}); !err.IsOk() {
		return err
	} else if len(*dbs) > 0 {
		return errors.BadRequest.WithMessage(fmt.Sprintf("network %s still has dbs", network.IdentifierID.Network))
	}
	firewalls, err := gcp.FindAllFirewalls(ctx, option.Option{Value: resourceRepo.FindAllResourceOption{Labels: searchLabels}})
	if !err.IsOk() {
		return err
	}
	for _, firewall := range *firewalls {
		if errFirewall := gcp.DeleteFirewall(ctx, &firewall); !errFirewall.IsOk() {
			return errFirewall
		}
	}
	gcpSubnetwork := gcp.toGCPNetwork(ctx, network)
	if errSubnet := kubernetes.Client().Client.Delete(ctx, gcpSubnetwork); errSubnet != nil {
		if k8serrors.IsNotFound(errSubnet) {
			return errors.NotFound.WithMessage(fmt.Sprintf("subnetwork %s not found", network.IdentifierName.Network))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to delete subnetwork %s", network.IdentifierName.Network))
	}
	return errors.NoContent
}

func (gcp *gcpRepository) DeleteNetworkCascade(ctx context.Context, network *networkModels.Network) errors.Error {
	searchLabels := lo.Assign(map[string]string{model.ProjectIDLabelKey: ctx.Value(context.ProjectIDKey).(string)}, network.IdentifierID.ToIDLabels())
	subnets, subnetsErr := gcp.FindAllSubnetworks(ctx, option.Option{Value: resourceRepo.FindAllResourceOption{Labels: searchLabels}})
	firewalls, firewallsErr := gcp.FindAllFirewalls(ctx, option.Option{Value: resourceRepo.FindAllResourceOption{Labels: searchLabels}})
	sqldbs, sqldbsErr := gcp.FindAllSqlDBs(ctx, option.Option{Value: resourceRepo.FindAllResourceOption{Labels: searchLabels}})
	if !subnetsErr.IsOk() {
		return subnetsErr
	}
	if !firewallsErr.IsOk() {
		return firewallsErr
	}
	if !sqldbsErr.IsOk() {
		return sqldbsErr
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
	for _, sqlDB := range *sqldbs {
		if sqlDBErr := gcp.DeleteSqlDB(ctx, &sqlDB); !sqlDBErr.IsOk() {
			return sqlDBErr
		}
	}
	return gcp.DeleteNetwork(ctx, network)

}

func (gcp *gcpRepository) NetworkExists(ctx context.Context, network *networkModels.Network) (bool, errors.Error) {
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

func (gcp *gcpRepository) toModelNetwork(network *v1beta1.Network) (*networkModels.Network, errors.Error) {
	id := identifier.Network{}
	name := identifier.Network{}
	if err := id.IDFromLabels(network.Labels); !err.IsOk() {
		return nil, err
	}
	if err := name.NameFromLabels(network.Labels); !err.IsOk() {
		return nil, err
	}
	return &networkModels.Network{
		Metadata: metadata.Metadata{
			Managed: network.Spec.ResourceSpec.DeletionPolicy == v1.DeletionDelete,
		},
		IdentifierID:   id,
		IdentifierName: name,
	}, errors.OK
}

func (gcp *gcpRepository) toGCPNetwork(ctx context.Context, network *networkModels.Network) *v1beta1.Network {
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

func (gcp *gcpRepository) toModelNetworkCollection(list *v1beta1.NetworkList) (*networkModels.Collection, errors.Error) {
	res := networkModels.Collection{}
	for _, item := range list.Items {
		network, err := gcp.toModelNetwork(&item)
		if !err.IsOk() {
			return &res, err
		}
		res[network.IdentifierName.Network] = *network
	}
	return &res, errors.OK
}
