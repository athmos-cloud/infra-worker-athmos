package aws

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	resourceRepo "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository/resource"
	"github.com/upbound/provider-aws/apis/networkfirewall/v1beta1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (aws *awsRepository) FindFirewall(ctx context.Context, opt option.Option) (*resource.Firewall, errors.Error) {
	if !opt.SetType(reflect.TypeOf(resourceRepo.FindResourceOption{}).String()).Validate() {
		return nil, errors.InvalidOption.WithMessage(fmt.Sprintf("invalid option : want %s, got %+v", reflect.TypeOf(resourceRepo.FindResourceOption{}).String(), opt.Get()))
	}
	req := opt.Get().(resourceRepo.FindResourceOption)
	awsFirewall := &v1beta1.Firewall{}
	if err := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, awsFirewall); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, errors.NotFound.WithMessage(fmt.Sprintf("firewall %s not found in namespace %s", req.Name, req.Namespace))
		}
		return nil, errors.KubernetesError.WithMessage(fmt.Sprintf("unable to get firewall %s in namespace %s", req.Name, req.Namespace))
	}
	mod, err := aws.toModelFirewall(awsFirewall)
	if !err.IsOk() {
		return nil, err
	}
	return mod, errors.OK
}

func (aws *awsRepository) FindAllFirewalls(ctx context.Context, opt option.Option) (*resource.FirewallCollection, errors.Error) {
	//TODO implement me
	//panic("implement me")
	return &resource.FirewallCollection{}, errors.OK
}

func (aws *awsRepository) FindAllRecursiveFirewalls(ctx context.Context, opt option.Option, ch *resourceRepo.FirewallChannel) {
	if !opt.SetType(reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String()).Validate() {
		ch.ErrorChannel <- errors.InvalidOption.WithMessage(fmt.Sprintf("invalid option : want %s, got %+v", reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String(), opt.Get()))
		return
	}
	req := opt.Get().(resourceRepo.FindAllResourceOption)
	gcpFirewallList := &v1beta1.FirewallList{}
	listOpt := &client.ListOptions{
		LabelSelector: client.MatchingLabelsSelector{Selector: labels.SelectorFromSet(req.Labels)},
	}
	if err := kubernetes.Client().Client.List(ctx, gcpFirewallList, listOpt); err != nil {
		ch.ErrorChannel <- errors.KubernetesError.WithMessage(fmt.Sprintf("unable to list firewalls in namespace %s", req.Namespace))
		return
	}
	if firewalls, err := aws.toModelFirewallCollection(gcpFirewallList); !err.IsOk() {
		ch.ErrorChannel <- err
	} else {
		ch.Channel <- firewalls
	}
}

func (aws *awsRepository) CreateFirewall(ctx context.Context, firewall *resource.Firewall) errors.Error {
	if exists, err := aws.FirewallExists(ctx, firewall); !err.IsOk() {
		return err
	} else if exists {
		return errors.Conflict.WithMessage(fmt.Sprintf("firewall %s already exists in namespace %s", firewall.IdentifierName.Firewall, firewall.Metadata.Namespace))
	}
	gcpFirewall := aws.toAWSFirewall(ctx, firewall)
	if err := kubernetes.Client().Client.Create(ctx, gcpFirewall); err != nil {
		if k8serrors.IsAlreadyExists(err) {
			return errors.Conflict.WithMessage(fmt.Sprintf("firewall %s already exists in namespace %s", firewall.IdentifierName.Firewall, firewall.Metadata.Namespace))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to create subnetwork %s in namespace %s", firewall.IdentifierName.Firewall, firewall.Metadata.Namespace))
	}
	return errors.Created
}

func (aws *awsRepository) UpdateFirewall(ctx context.Context, firewall *resource.Firewall) errors.Error {
	existingFirewall := &v1beta1.Firewall{}
	if err := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: firewall.IdentifierID.Firewall, Namespace: firewall.Metadata.Namespace}, existingFirewall); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("firewall %s not found", firewall.IdentifierID.Firewall))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to get subnetwork %s", firewall.IdentifierID.Firewall))
	}
	gcpFirewall := aws.toAWSFirewall(ctx, firewall)
	existingFirewall.Spec = gcpFirewall.Spec
	existingFirewall.Labels = gcpFirewall.Labels
	if err := kubernetes.Client().Client.Update(ctx, existingFirewall); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("firewall %s not found in namespace %s", firewall.IdentifierName.Firewall, firewall.Metadata.Namespace))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to update firewall %s in namespace %s", firewall.IdentifierName.Firewall, firewall.Metadata.Namespace))
	}
	return errors.NoContent
}

func (aws *awsRepository) DeleteFirewall(ctx context.Context, firewall *resource.Firewall) errors.Error {
	gcpFirewall := aws.toAWSFirewall(ctx, firewall)
	if err := kubernetes.Client().Client.Delete(ctx, gcpFirewall); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("firewall %s not found in namespace %s", firewall.IdentifierName.Firewall, firewall.Metadata.Namespace))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to delete firewall %s in namespace %s", firewall.IdentifierName.Firewall, firewall.Metadata.Namespace))
	}
	return errors.NoContent
}

func (aws *awsRepository) FirewallExists(ctx context.Context, firewall *resource.Firewall) (bool, errors.Error) {
	/*	searchLabels := map[string]string{
			model.ProjectIDLabelKey:          ctx.Value(context.ProjectIDKey).(string),
			identifier.ProviderIdentifierKey: firewall.IdentifierID.Provider,
			identifier.NetworkIdentifierKey:  firewall.IdentifierID.Network,
			identifier.FirewallNameKey:       firewall.IdentifierName.Firewall,
		}
		gcpFirewalls := &v1beta1.FirewallList{}
		listOpt := &client.ListOptions{
			LabelSelector: client.MatchingLabelsSelector{Selector: labels.SelectorFromSet(searchLabels)},
		}
		if err := kubernetes.Client().Client.List(ctx, gcpFirewalls, listOpt); err != nil {
			return false, errors.KubernetesError.WithMessage("unable to list firewalls")
		}
		return len(gcpFirewalls.Items) > 0, errors.OK*/
	return false, errors.OK
}

func (aws *awsRepository) toModelFirewall(firewall *v1beta1.Firewall) (*resource.Firewall, errors.Error) {
	/*	id := identifier.Firewall{}
		name := identifier.Firewall{}
		if err := id.IDFromLabels(firewall.Labels); !err.IsOk() {
			return nil, err
		}
		if err := name.NameFromLabels(firewall.Labels); !err.IsOk() {
			return nil, err
		}
		allow := resource.FirewallRuleList{}
		mapAllow := make(map[string][]string)
		for _, a := range firewall.Spec.ForProvider.Allow {
			var portsA []string
			for _, p := range a.Ports {
				portsA = append(portsA, *p)
			}
			mapAllow[*a.Protocol] = append(mapAllow[*a.Protocol], portsA...)
		}
		for k, v := range mapAllow {
			rule := resource.FirewallRule{
				Protocol: k,
				Ports:    v,
			}
			allow = append(allow, rule)
		}

		deny := resource.FirewallRuleList{}
		mapDeny := make(map[string][]string)
		for _, d := range firewall.Spec.ForProvider.Deny {
			var portsD []string
			for _, p := range d.Ports {
				portsD = append(portsD, *p)
			}
			mapDeny[*d.Protocol] = append(mapDeny[*d.Protocol], portsD...)
		}
		for k, v := range mapDeny {
			rule := resource.FirewallRule{
				Protocol: k,
				Ports:    v,
			}
			deny = append(deny, rule)
		}

		return &resource.Firewall{
			Metadata: metadata.Metadata{
				Managed:   firewall.Spec.ResourceSpec.DeletionPolicy == v1.DeletionDelete,
				Namespace: firewall.ObjectMeta.Namespace,
			},
			IdentifierID:   id,
			IdentifierName: name,
			Allow:          allow,
			Deny:           deny,
		}, errors.OK*/
	return &resource.Firewall{}, errors.OK
}

func (aws *awsRepository) toAWSFirewall(ctx context.Context, firewall *resource.Firewall) *v1beta1.Firewall {
	/*	resLabels := lo.Assign(crossplane.GetBaseLabels(ctx.Value(context.ProjectIDKey).(string)), firewall.IdentifierID.ToIDLabels(), firewall.IdentifierName.ToNameLabels())
		var allow []v1beta1.AllowParameters
		for _, a := range firewall.Allow {
			for _, p := range a.Ports {
				curP := p
				allow = append(allow, v1beta1.AllowParameters{Protocol: &a.Protocol, Ports: []*string{&curP}})
			}
		}
		var deny []v1beta1.DenyParameters
		for _, d := range firewall.Deny {
			for _, p := range d.Ports {
				curP := p
				deny = append(deny, v1beta1.DenyParameters{Protocol: &d.Protocol, Ports: []*string{&curP}})
			}
		}
		return &v1beta1.Firewall{
			ObjectMeta: metav1.ObjectMeta{
				Name:        firewall.IdentifierID.Firewall,
				Namespace:   firewall.Metadata.Namespace,
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
					VPCID: &firewall.IdentifierID.Network,
					Allow: allow,
					Deny:  deny,
				},
			},
		}*/
	return &v1beta1.Firewall{}

}

func (aws *awsRepository) toModelFirewallCollection(firewallList *v1beta1.FirewallList) (*resource.FirewallCollection, errors.Error) {
	items := resource.FirewallCollection{}
	for _, item := range firewallList.Items {
		firewall, err := aws.toModelFirewall(&item)
		if !err.IsOk() {
			return nil, err
		}
		items[firewall.IdentifierName.Firewall] = *firewall
	}
	return &items, errors.OK
}
