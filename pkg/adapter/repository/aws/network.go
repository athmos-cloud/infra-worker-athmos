package aws

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
	"github.com/upbound/provider-aws/apis/ec2/v1beta1"
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

func (aws *awsRepository) FindNetwork(ctx context.Context, opt option.Option) (*networkModels.Network, errors.Error) {
	if !opt.SetType(reflect.TypeOf(resourceRepo.FindResourceOption{}).String()).Validate() {
		return nil, errors.InvalidOption.WithMessage(fmt.Sprintf("invalid option : want %s, got %+v", reflect.TypeOf(resourceRepo.FindResourceOption{}).String(), opt.Get()))
	}
	req := opt.Get().(resourceRepo.FindResourceOption)
	awsNetwork := &v1beta1.VPC{}
	if err := kubernetes.Client().Client.Get(ctx, client.ObjectKey{Name: req.Name}, awsNetwork); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, errors.NotFound.WithMessage(fmt.Sprintf("networkModels %s not found", req.Name))
		}
		return nil, errors.KubernetesError.WithMessage(fmt.Sprintf("unable to get networkModels %s", req.Name))
	}

	mod, err := aws.toModelNetwork(awsNetwork)
	if !err.IsOk() {
		return nil, err
	}
	return mod, errors.OK
}

func (aws *awsRepository) FindAllNetworks(ctx context.Context, opt option.Option) (*networkModels.NetworkCollection, errors.Error) {
	if !opt.SetType(reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String()).Validate() {
		return nil, errors.InvalidOption.WithMessage(fmt.Sprintf("invalid option : want %s, got %+v", reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String(), opt.Get()))
	}
	req := opt.Get().(resourceRepo.FindAllResourceOption)
	awsNetworkList := &v1beta1.VPCList{}
	kubeOptions := &client.ListOptions{
		LabelSelector: client.MatchingLabelsSelector{Selector: labels.SelectorFromSet(req.Labels)},
	}
	if err := kubernetes.Client().Client.List(ctx, awsNetworkList, kubeOptions); err != nil {
		return nil, errors.KubernetesError.WithMessage("unable to list networks")
	}
	modNetworks, err := aws.toModelNetworkCollection(awsNetworkList)
	if !err.IsOk() {
		return nil, err
	}
	return modNetworks, errors.OK
}

func (aws *awsRepository) FindAllRecursiveNetworks(ctx context.Context, opt option.Option, ch *resourceRepo.NetworkChannel) (*networkModels.NetworkCollection, errors.Error) {
	if !opt.SetType(reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String()).Validate() {
		return nil, errors.InvalidOption.WithMessage(fmt.Sprintf("invalid option : want %s, got %+v", reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String(), opt.Get()))
	}
	req := opt.Get().(resourceRepo.FindAllResourceOption)
	awsVPCList := &v1beta1.VPCList{}
	kubeOptions := &client.ListOptions{
		LabelSelector: client.MatchingLabelsSelector{Selector: labels.SelectorFromSet(req.Labels)},
	}
	if err := kubernetes.Client().Client.List(ctx, awsVPCList, kubeOptions); err != nil {
		return nil, errors.KubernetesError.WithMessage(fmt.Sprintf("unable to list networks"))
	}

	modNetworks, err := aws.toModelNetworkCollection(awsVPCList)
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

		go aws.FindAllRecursiveFirewalls(ctx, option.Option{Value: resourceRepo.FindAllResourceOption{Labels: network.IdentifierID.ToIDLabels()}}, chFirewall)
		go aws.FindAllRecursiveSubnetworks(ctx, option.Option{Value: resourceRepo.FindAllResourceOption{Labels: network.IdentifierID.ToIDLabels()}}, chSubnet)
		go aws.FindAllRecursiveSqlDBs(ctx, option.Option{Value: resourceRepo.FindAllResourceOption{Labels: network.IdentifierID.ToIDLabels()}}, chDB)

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
	networks := &networkModels.NetworkCollection{}
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

func (aws *awsRepository) CreateNetwork(ctx context.Context, network *networkModels.Network) errors.Error {
	if exist, errExists := aws.NetworkExists(ctx, network); !errExists.IsOk() {
		return errExists
	} else if exist {
		return errors.Conflict.WithMessage(fmt.Sprintf("networkModels %s already exists", network.IdentifierName.Network))
	}
	awsNetwork, err := aws.toAWSNetwork(ctx, network)
	if !err.IsOk() {
		return err
	}

	if err := kubernetes.Client().Client.Create(ctx, awsNetwork); err != nil {
		if k8serrors.IsAlreadyExists(err) {
			return errors.Conflict.WithMessage(fmt.Sprintf("network %s already exists", network.IdentifierName.Network))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to create network %s", network.IdentifierName.Network))
	}
	return errors.Created
}

func (aws *awsRepository) UpdateNetwork(ctx context.Context, network *networkModels.Network) errors.Error {
	existingNetwork := &v1beta1.VPC{}
	if err := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: network.IdentifierID.Network}, existingNetwork); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("network %s not found", network.IdentifierID.Network))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to get network %s", network.IdentifierID.Network))
	}

	awsNetwork, err := aws.toAWSNetwork(ctx, network)
	if !err.IsOk() {
		return err
	}

	existingNetwork.Spec = awsNetwork.Spec
	existingNetwork.Labels = awsNetwork.Labels

	if err := kubernetes.Client().Client.Update(ctx, existingNetwork); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("network %s not found", network.IdentifierName.Network))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to update network %s", network.IdentifierName.Network))
	}
	return errors.NoContent
}

func (aws *awsRepository) DeleteNetwork(ctx context.Context, network *networkModels.Network) errors.Error {
	searchLabels := lo.Assign(map[string]string{model.ProjectIDLabelKey: ctx.Value(context.ProjectIDKey).(string)}, network.IdentifierID.ToIDLabels())
	if subnets, err := aws.FindAllSubnetworks(ctx, option.Option{Value: resourceRepo.FindAllResourceOption{Labels: searchLabels}}); !err.IsOk() {
		return err
	} else if len(*subnets) > 0 {
		return errors.BadRequest.WithMessage(fmt.Sprintf("can't delete networkModels %s without cascade option", network.IdentifierName.Network))
	}

	awsNetwork, err := aws.toAWSNetwork(ctx, network)
	if !err.IsOk() {
		return err
	}

	if err := kubernetes.Client().Client.Delete(ctx, awsNetwork); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("subnetwork %s not found", network.IdentifierName.Network))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to delete subnetwork %s", network.IdentifierName.Network))
	}
	return errors.NoContent
}

func (aws *awsRepository) DeleteNetworkCascade(ctx context.Context, network *networkModels.Network) errors.Error {
	searchLabels := lo.Assign(map[string]string{model.ProjectIDLabelKey: ctx.Value(context.ProjectIDKey).(string)}, network.IdentifierID.ToIDLabels())
	subnets, subnetsErr := aws.FindAllSubnetworks(ctx, option.Option{Value: resourceRepo.FindAllResourceOption{Labels: searchLabels}})
	firewalls, firewallsErr := aws.FindAllFirewalls(ctx, option.Option{Value: resourceRepo.FindAllResourceOption{Labels: searchLabels}})

	if !subnetsErr.IsOk() {
		return subnetsErr
	}
	if !firewallsErr.IsOk() {
		return firewallsErr
	}

	for _, subnet := range *subnets {
		if subnetErr := aws.DeleteSubnetworkCascade(ctx, &subnet); !subnetErr.IsOk() {
			return subnetErr
		}
	}

	for _, firewall := range *firewalls {
		if firewallErr := aws.DeleteFirewall(ctx, &firewall); !firewallErr.IsOk() {
			return firewallErr
		}
	}
	return aws.DeleteNetwork(ctx, network)

}

func (aws *awsRepository) NetworkExists(ctx context.Context, network *networkModels.Network) (bool, errors.Error) {
	awsNetwork := &v1beta1.VPCList{}
	searchLabels := lo.Assign(map[string]string{
		model.ProjectIDLabelKey:          ctx.Value(context.ProjectIDKey).(string),
		identifier.ProviderIdentifierKey: network.IdentifierID.Provider,
		identifier.NetworkNameKey:        network.IdentifierName.Network,
	})
	kubeOptions := &client.ListOptions{
		LabelSelector: client.MatchingLabelsSelector{Selector: labels.SelectorFromSet(searchLabels)},
	}
	if err := kubernetes.Client().Client.List(ctx, awsNetwork, kubeOptions); err != nil {
		return false, errors.KubernetesError.WithMessage("unable to list subnetworks in namespace")
	}
	return len(awsNetwork.Items) > 0, errors.OK
}

func (aws *awsRepository) toModelNetwork(network *v1beta1.VPC) (*networkModels.Network, errors.Error) {
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
			Status:  metadata.StatusFromKubernetesStatus(network.Status.Conditions),
			Managed: network.Spec.ResourceSpec.DeletionPolicy == v1.DeletionDelete,
		},
		IdentifierID:   id,
		IdentifierName: name,
		Region:         *network.Spec.ForProvider.Region,
	}, errors.OK
}

func (aws *awsRepository) toAWSNetwork(ctx context.Context, network *networkModels.Network) (*v1beta1.VPC, errors.Error) {
	resLabels := lo.Assign(crossplane.GetBaseLabels(ctx.Value(context.ProjectIDKey).(string)), network.IdentifierID.ToIDLabels(), network.IdentifierName.ToNameLabels())
	if network.Region == "" {
		return nil, errors.BadRequest.WithMessage("Region is a required field.")
	}

	return &v1beta1.VPC{
		ObjectMeta: metav1.ObjectMeta{
			Name:        network.IdentifierID.Network,
			Labels:      resLabels,
			Annotations: crossplane.GetAnnotations(network.Metadata.Managed, network.IdentifierName.Network),
		},
		Spec: v1beta1.VPCSpec{
			ResourceSpec: v1.ResourceSpec{
				DeletionPolicy: crossplane.GetDeletionPolicy(network.Metadata.Managed),
				ProviderConfigReference: &v1.Reference{
					Name: network.IdentifierID.Provider,
				},
			},
			ForProvider: v1beta1.VPCParameters_2{
				Region: &network.Region,
			},
		},
	}, errors.OK
}

func (aws *awsRepository) toModelNetworkCollection(list *v1beta1.VPCList) (*networkModels.NetworkCollection, errors.Error) {
	res := networkModels.NetworkCollection{}
	for _, item := range list.Items {
		network, err := aws.toModelNetwork(&item)
		if !err.IsOk() {
			return &res, err
		}
		res[network.IdentifierName.Network] = *network
	}
	return &res, errors.OK
}
