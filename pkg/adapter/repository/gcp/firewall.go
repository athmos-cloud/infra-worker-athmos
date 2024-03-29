package gcp

import (
	"fmt"
	"reflect"

	v1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/samber/lo"
	"github.com/upbound/provider-gcp/apis/compute/v1beta1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

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
)

func (gcp *gcpRepository) FindFirewall(ctx context.Context, opt option.Option) (*network.Firewall, errors.Error) {
	if !opt.SetType(reflect.TypeOf(resourceRepo.FindResourceOption{}).String()).Validate() {
		return nil, errors.InvalidOption.WithMessage(fmt.Sprintf("invalid option : want %s, got %+v", reflect.TypeOf(resourceRepo.FindResourceOption{}).String(), opt.Get()))
	}
	req := opt.Get().(resourceRepo.FindResourceOption)
	gcpFirewall := &v1beta1.Firewall{}
	if err := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: req.Name}, gcpFirewall); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, errors.NotFound.WithMessage(fmt.Sprintf("firewall %s not found", req.Name))
		}
		return nil, errors.KubernetesError.WithMessage(fmt.Sprintf("unable to get firewall %s", req.Name))
	}
	mod, err := gcp.toModelFirewall(gcpFirewall)
	if !err.IsOk() {
		return nil, err
	}
	return mod, errors.OK
}

func (gcp *gcpRepository) FindAllFirewalls(ctx context.Context, opt option.Option) (*network.FirewallCollection, errors.Error) {
	if !opt.SetType(reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String()).Validate() {
		return nil, errors.InvalidOption.WithMessage(fmt.Sprintf("invalid option : want %s, got %+v", reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String(), opt.Get()))
	}
	req := opt.Get().(resourceRepo.FindAllResourceOption)
	gcpFirewallList := &v1beta1.FirewallList{}
	listOpt := &client.ListOptions{
		LabelSelector: client.MatchingLabelsSelector{Selector: labels.SelectorFromSet(req.Labels)},
	}
	if err := kubernetes.Client().Client.List(ctx, gcpFirewallList, listOpt); err != nil {
		return nil, errors.KubernetesError.WithMessage(fmt.Sprintf("unable to get firewalls"))
	}
	firewallCollection, err := gcp.toModelFirewallCollection(gcpFirewallList)
	if !err.IsOk() {
		return nil, err
	}
	return firewallCollection, errors.OK
}

func (gcp *gcpRepository) FindAllRecursiveFirewalls(ctx context.Context, opt option.Option, ch *resourceRepo.FirewallChannel) {
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
		ch.ErrorChannel <- errors.KubernetesError.WithMessage(fmt.Sprintf("unable to list firewalls"))
		return
	}
	if firewalls, err := gcp.toModelFirewallCollection(gcpFirewallList); !err.IsOk() {
		ch.ErrorChannel <- err
	} else {
		ch.Channel <- firewalls
	}
}

func (gcp *gcpRepository) CreateFirewall(ctx context.Context, firewall *network.Firewall) errors.Error {
	if exists, err := gcp.FirewallExists(ctx, firewall); !err.IsOk() {
		return err
	} else if exists {
		return errors.Conflict.WithMessage(fmt.Sprintf("firewall %s already exists", firewall.IdentifierName.Firewall))
	}
	gcpFirewall := gcp.toGCPFirewall(ctx, firewall)
	if err := kubernetes.Client().Client.Create(ctx, gcpFirewall); err != nil {
		if k8serrors.IsAlreadyExists(err) {
			return errors.Conflict.WithMessage(fmt.Sprintf("firewall %s already exists", firewall.IdentifierName.Firewall))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to create subnetwork %s", firewall.IdentifierName.Firewall))
	}
	return errors.Created
}

func (gcp *gcpRepository) UpdateFirewall(ctx context.Context, firewall *network.Firewall) errors.Error {
	existingFirewall := &v1beta1.Firewall{}
	if err := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: firewall.IdentifierID.Firewall}, existingFirewall); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("firewall %s not found", firewall.IdentifierID.Firewall))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to get subnetwork %s", firewall.IdentifierID.Firewall))
	}
	gcpFirewall := gcp.toGCPFirewall(ctx, firewall)
	existingFirewall.Spec = gcpFirewall.Spec
	existingFirewall.Labels = gcpFirewall.Labels
	if err := kubernetes.Client().Client.Update(ctx, existingFirewall); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("firewall %s not found", firewall.IdentifierName.Firewall))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to update firewall %s", firewall.IdentifierName.Firewall))
	}
	return errors.NoContent
}

func (gcp *gcpRepository) DeleteFirewall(ctx context.Context, firewall *network.Firewall) errors.Error {
	gcpFirewall := gcp.toGCPFirewall(ctx, firewall)
	if err := kubernetes.Client().Client.Delete(ctx, gcpFirewall); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("firewall %s not found", firewall.IdentifierName.Firewall))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to delete firewall %s", firewall.IdentifierName.Firewall))
	}
	return errors.NoContent
}

func (gcp *gcpRepository) FirewallExists(ctx context.Context, firewall *network.Firewall) (bool, errors.Error) {
	searchLabels := map[string]string{
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
	return len(gcpFirewalls.Items) > 0, errors.OK
}

func (gcp *gcpRepository) toModelFirewall(firewall *v1beta1.Firewall) (*network.Firewall, errors.Error) {
	id := identifier.Firewall{}
	name := identifier.Firewall{}
	if err := id.IDFromLabels(firewall.Labels); !err.IsOk() {
		return nil, err
	}
	if err := name.NameFromLabels(firewall.Labels); !err.IsOk() {
		return nil, err
	}
	allow := network.FirewallRuleList{}
	mapAllow := make(map[string][]string)
	for _, a := range firewall.Spec.ForProvider.Allow {
		var portsA []string
		for _, p := range a.Ports {
			portsA = append(portsA, *p)
		}
		mapAllow[*a.Protocol] = append(mapAllow[*a.Protocol], portsA...)
	}
	for k, v := range mapAllow {
		rule := network.FirewallRule{
			Protocol: k,
			Ports:    v,
		}
		allow = append(allow, rule)
	}

	deny := network.FirewallRuleList{}
	mapDeny := make(map[string][]string)
	for _, d := range firewall.Spec.ForProvider.Deny {
		var portsD []string
		for _, p := range d.Ports {
			portsD = append(portsD, *p)
		}
		mapDeny[*d.Protocol] = append(mapDeny[*d.Protocol], portsD...)
	}
	for k, v := range mapDeny {
		rule := network.FirewallRule{
			Protocol: k,
			Ports:    v,
		}
		deny = append(deny, rule)
	}

	return &network.Firewall{
		Metadata: metadata.Metadata{
			Status:  metadata.StatusFromKubernetesStatus(firewall.Status.Conditions),
			Managed: firewall.Spec.ResourceSpec.DeletionPolicy == v1.DeletionDelete,
		},
		IdentifierID:   id,
		IdentifierName: name,
		Allow:          allow,
		Deny:           deny,
	}, errors.OK
}

func (gcp *gcpRepository) toGCPFirewall(ctx context.Context, firewall *network.Firewall) *v1beta1.Firewall {
	resLabels := lo.Assign(crossplane.GetBaseLabels(ctx.Value(context.ProjectIDKey).(string)), firewall.IdentifierID.ToIDLabels(), firewall.IdentifierName.ToNameLabels())
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
	defaultSourceRanges := "0.0.0.0/0"
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
				Network:      &firewall.IdentifierID.Network,
				Project:      &firewall.IdentifierID.VPC,
				Allow:        allow,
				Deny:         deny,
				SourceRanges: []*string{&defaultSourceRanges},
			},
		},
	}

}

func (gcp *gcpRepository) toModelFirewallCollection(firewallList *v1beta1.FirewallList) (*network.FirewallCollection, errors.Error) {
	items := network.FirewallCollection{}
	for _, item := range firewallList.Items {
		firewall, err := gcp.toModelFirewall(&item)
		if !err.IsOk() {
			return nil, err
		}
		items[firewall.IdentifierName.Firewall] = *firewall
	}
	return &items, errors.OK
}
