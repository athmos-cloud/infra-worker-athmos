package aws

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/crossplane"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/network"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	resourceRepo "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository/resource"
	v1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/samber/lo"
	awsCompute "github.com/upbound/provider-aws/apis/ec2/v1beta1"
	"github.com/upbound/provider-aws/apis/networkfirewall/v1beta1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (aws *awsRepository) FindFirewall(ctx context.Context, opt option.Option) (*network.Firewall, errors.Error) {
	if !opt.SetType(reflect.TypeOf(resourceRepo.FindResourceOption{}).String()).Validate() {
		return nil, errors.InvalidOption.WithMessage(fmt.Sprintf("invalid option : want %s, got %+v", reflect.TypeOf(resourceRepo.FindResourceOption{}).String(), opt.Get()))
	}
	req := opt.Get().(resourceRepo.FindResourceOption)
	awsFirewall := &v1beta1.Firewall{}
	if err := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: req.Name}, awsFirewall); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, errors.NotFound.WithMessage(fmt.Sprintf("firewall %s not found", req.Name))
		}
		return nil, errors.KubernetesError.WithMessage(fmt.Sprintf("unable to get firewall %s", req.Name))
	}

	awsRuleGroupName := fmt.Sprintf("%s-rule-group", awsFirewall.Name)
	awsRuleGroup, err := aws._getRuleGroup(ctx, &awsRuleGroupName)
	if !err.IsOk() {
		return nil, err
	}

	mod, err := aws.toModelFirewall(awsFirewall, awsRuleGroup)
	if !err.IsOk() {
		return nil, err
	}
	return mod, errors.OK
}

func (aws *awsRepository) FindAllFirewalls(ctx context.Context, opt option.Option) (*network.FirewallCollection, errors.Error) {
	//TODO implement me
	//panic("implement me")
	return &network.FirewallCollection{}, errors.OK
}

func (aws *awsRepository) FindAllRecursiveFirewalls(ctx context.Context, opt option.Option, ch *resourceRepo.FirewallChannel) {
	if !opt.SetType(reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String()).Validate() {
		ch.ErrorChannel <- errors.InvalidOption.WithMessage(fmt.Sprintf("invalid option : want %s, got %+v", reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String(), opt.Get()))
		return
	}
	req := opt.Get().(resourceRepo.FindAllResourceOption)
	awsFirewallList := &v1beta1.FirewallList{}
	listOpt := &client.ListOptions{
		LabelSelector: client.MatchingLabelsSelector{Selector: labels.SelectorFromSet(req.Labels)},
	}
	if err := kubernetes.Client().Client.List(ctx, awsFirewallList, listOpt); err != nil {
		ch.ErrorChannel <- errors.KubernetesError.WithMessage("unable to list firewalls in namespace")
		return
	}
	if firewalls, err := aws.toModelFirewallCollection(ctx, awsFirewallList); !err.IsOk() {
		ch.ErrorChannel <- err
	} else {
		ch.Channel <- firewalls
	}
}

func (aws *awsRepository) CreateFirewall(ctx context.Context, firewall *network.Firewall) errors.Error {
	if exists, err := aws.FirewallExists(ctx, firewall); !err.IsOk() {
		return err
	} else if exists {
		return errors.Conflict.WithMessage(fmt.Sprintf("firewall %s already exists", firewall.IdentifierName.Firewall))
	}
	region, err := _regionForFirewall(ctx, firewall)
	if !err.IsOk() {
		return err
	}

	err = aws._createRuleGroup(ctx, firewall, region)
	if !err.IsOk() {
		return err
	}

	err = aws._createFirewallPolicy(ctx, firewall, region)
	if !err.IsOk() {
		return err
	}

	awsFirewall := aws.toAWSFirewall(ctx, firewall, region)

	if err := kubernetes.Client().Client.Create(ctx, awsFirewall); err != nil {
		if k8serrors.IsAlreadyExists(err) {
			return errors.Conflict.WithMessage(fmt.Sprintf("firewall %s already exists", firewall.IdentifierName.Firewall))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to create subnetwork %s", firewall.IdentifierName.Firewall))
	}
	return errors.Created
}

func (aws *awsRepository) UpdateFirewall(ctx context.Context, firewall *network.Firewall) errors.Error {
	existingFirewall := &v1beta1.Firewall{}
	if err := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: firewall.IdentifierID.Firewall}, existingFirewall); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("firewall %s not found", firewall.IdentifierID.Firewall))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to get subnetwork %s", firewall.IdentifierID.Firewall))
	}

	region, err := _regionForFirewall(ctx, firewall)
	if !err.IsOk() {
		return err
	}

	awsFirewall := aws.toAWSFirewall(ctx, firewall, region)
	existingFirewall.Spec = awsFirewall.Spec
	existingFirewall.Labels = awsFirewall.Labels
	if err := kubernetes.Client().Client.Update(ctx, existingFirewall); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("firewall %s not found", firewall.IdentifierName.Firewall))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to update firewall %s", firewall.IdentifierName.Firewall))
	}

	err = aws._updateRuleGroup(ctx, firewall, region)
	if !err.IsOk() {
		return err
	}

	return errors.NoContent
}

func (aws *awsRepository) DeleteFirewall(ctx context.Context, firewall *network.Firewall) errors.Error {
	region, err := _regionForFirewall(ctx, firewall)
	if !err.IsOk() {
		return err
	}

	awsFirewall := aws.toAWSFirewall(ctx, firewall, region)
	if err := kubernetes.Client().Client.Delete(ctx, awsFirewall); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("firewall %s not found", firewall.IdentifierName.Firewall))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to delete firewall %s", firewall.IdentifierName.Firewall))
	}

	err = aws._deleteFirewallPolicy(ctx, firewall, region)
	if !err.IsOk() {
		return err
	}

	err = aws._deleteRuleGroup(ctx, firewall, region)
	if !err.IsOk() {
		return err
	}

	return errors.NoContent
}

func (aws *awsRepository) FirewallExists(ctx context.Context, firewall *network.Firewall) (bool, errors.Error) {
	searchLabels := map[string]string{
		model.ProjectIDLabelKey:          ctx.Value(context.ProjectIDKey).(string),
		identifier.ProviderIdentifierKey: firewall.IdentifierID.Provider,
		identifier.NetworkIdentifierKey:  firewall.IdentifierID.Network,
		identifier.FirewallNameKey:       firewall.IdentifierName.Firewall,
	}
	awsFirewalls := &v1beta1.FirewallList{}
	listOpt := &client.ListOptions{
		LabelSelector: client.MatchingLabelsSelector{Selector: labels.SelectorFromSet(searchLabels)},
	}
	if err := kubernetes.Client().Client.List(ctx, awsFirewalls, listOpt); err != nil {
		return false, errors.KubernetesError.WithMessage("unable to list firewalls")
	}
	return len(awsFirewalls.Items) > 0, errors.OK
}

func (aws *awsRepository) toModelFirewall(firewall *v1beta1.Firewall, ruleGroup *v1beta1.RuleGroup) (*network.Firewall, errors.Error) {
	id := identifier.Firewall{}
	name := identifier.Firewall{}
	if err := id.IDFromLabels(firewall.Labels); !err.IsOk() {
		return nil, err
	}
	if err := name.NameFromLabels(firewall.Labels); !err.IsOk() {
		return nil, err
	}

	allow, deny, err := aws._toFirewallRuleList(ruleGroup)
	if !err.IsOk() {
		return nil, err
	}

	return &network.Firewall{
		Metadata: metadata.Metadata{
			Managed: firewall.Spec.ResourceSpec.DeletionPolicy == v1.DeletionDelete,
		},
		IdentifierID:   id,
		IdentifierName: name,
		Allow:          *allow,
		Deny:           *deny,
	}, errors.OK
}

func (aws *awsRepository) toAWSFirewall(ctx context.Context, firewall *network.Firewall, region *string) *v1beta1.Firewall {
	resLabels := lo.Assign(crossplane.GetBaseLabels(ctx.Value(context.ProjectIDKey).(string)), firewall.IdentifierID.ToIDLabels(), firewall.IdentifierName.ToNameLabels())

	return &v1beta1.Firewall{
		ObjectMeta: metav1.ObjectMeta{
			Name:        firewall.IdentifierID.Firewall,
			Labels:      resLabels,
			Annotations: crossplane.GetAnnotations(firewall.Metadata.Managed, firewall.IdentifierName.Firewall),
		},
		Spec: v1beta1.FirewallSpec{
			ResourceSpec: v1.ResourceSpec{
				DeletionPolicy: crossplane.GetDeletionPolicy(firewall.Metadata.Managed),
				ProviderConfigReference: &v1.Reference{
					Name: firewall.IdentifierID.Provider,
				},
			},
			ForProvider: v1beta1.FirewallParameters{
				VPCID:  &firewall.IdentifierID.Network,
				Region: region,
			},
		},
	}
	return &v1beta1.Firewall{}

}

func (aws *awsRepository) toModelFirewallCollection(ctx context.Context, firewallList *v1beta1.FirewallList) (*network.FirewallCollection, errors.Error) {
	items := network.FirewallCollection{}
	for _, item := range firewallList.Items {
		awsRuleGroupName := fmt.Sprintf("%s-rule-group", item.Name)
		awsRuleGroup, errRG := aws._getRuleGroup(ctx, &awsRuleGroupName)
		if !errRG.IsOk() {
			return nil, errRG
		}

		firewall, err := aws.toModelFirewall(&item, awsRuleGroup)
		if !err.IsOk() {
			return nil, err
		}

		items[firewall.IdentifierName.Firewall] = *firewall
	}
	return &items, errors.OK
}

func _regionForFirewall(ctx context.Context, firewall *network.Firewall) (*string, errors.Error) {
	awsVPC := &awsCompute.VPC{}
	if err := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: firewall.IdentifierID.Network}, awsVPC); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, errors.NotFound.WithMessage(fmt.Sprintf("Region for firewall %s not found", firewall.IdentifierID.Firewall))
		}
		return nil, errors.KubernetesError.WithMessage(fmt.Sprintf("Unable to get region for firewall %s", firewall.IdentifierID.Firewall))
	}
	return awsVPC.Spec.ForProvider.Region, errors.OK
}
