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
)

func (gcp *gcpRepository) FindFirewall(ctx context.Context, opt option.Option) (*resource.Firewall, errors.Error) {
	if !opt.SetType(reflect.TypeOf(resourceRepo.FindResourceOption{}).String()).Validate() {
		return nil, errors.InvalidOption.WithMessage(fmt.Sprintf("invalid option : want %s, got %+v", reflect.TypeOf(resourceRepo.FindResourceOption{}).String(), opt.Get()))
	}
	req := opt.Get().(resourceRepo.FindResourceOption)
	gcpFirewall := &v1beta1.Firewall{}
	if err := kubernetes.Client().Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, gcpFirewall); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, errors.NotFound.WithMessage(fmt.Sprintf("firewall %s not found in namespace %s", req.Name, req.Namespace))
		}
		return nil, errors.KubernetesError.WithMessage(fmt.Sprintf("unable to get firewall %s in namespace %s", req.Name, req.Namespace))
	}
	mod, err := gcp.toModelFirewall(gcpFirewall)
	if !err.IsOk() {
		return nil, err
	}
	return mod, errors.OK
}

func (gcp *gcpRepository) FindAllFirewalls(ctx context.Context, opt option.Option) (*resource.FirewallCollection, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (gcp *gcpRepository) FindAllRecursiveFirewalls(ctx context.Context, opt option.Option, ch *resourceRepo.FirewallChannel) {
	if !opt.SetType(reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String()).Validate() {
		ch.ErrorChannel <- errors.InvalidOption.WithMessage(fmt.Sprintf("invalid option : want %s, got %+v", reflect.TypeOf(resourceRepo.FindAllResourceOption{}).String(), opt.Get()))
		return
	}
	req := opt.Get().(resourceRepo.FindAllResourceOption)
	gcpFirewallList := &v1beta1.FirewallList{}
	listOpt := &client.ListOptions{
		Namespace:     req.Namespace,
		LabelSelector: client.MatchingLabelsSelector{Selector: labels.SelectorFromSet(req.Labels)},
	}
	if err := kubernetes.Client().List(ctx, gcpFirewallList, listOpt); err != nil {
		ch.ErrorChannel <- errors.KubernetesError.WithMessage(fmt.Sprintf("unable to list firewalls in namespace %s", req.Namespace))
		return
	}
	if firewalls, err := gcp.toModelFirewallCollection(gcpFirewallList); !err.IsOk() {
		ch.ErrorChannel <- err
	} else {
		ch.Channel <- firewalls
	}
}

func (gcp *gcpRepository) CreateFirewall(ctx context.Context, firewall *resource.Firewall) errors.Error {
	gcpFirewall := gcp.toGCPFirewall(ctx, firewall)
	if err := kubernetes.Client().Create(ctx, gcpFirewall); err != nil {
		if k8serrors.IsAlreadyExists(err) {
			return errors.Conflict.WithMessage(fmt.Sprintf("firewall %s already exists in namespace %s", firewall.IdentifierName.Firewall, firewall.Metadata.Namespace))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to create subnetwork %s in namespace %s", firewall.IdentifierName.Firewall, firewall.Metadata.Namespace))
	}
	return errors.Created
}

func (gcp *gcpRepository) UpdateFirewall(ctx context.Context, firewall *resource.Firewall) errors.Error {
	gcpFirewall := gcp.toGCPFirewall(ctx, firewall)
	if err := kubernetes.Client().Update(ctx, gcpFirewall); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("firewall %s not found in namespace %s", firewall.IdentifierName.Firewall, firewall.Metadata.Namespace))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to update firewall %s in namespace %s", firewall.IdentifierName.Firewall, firewall.Metadata.Namespace))
	}
	return errors.NoContent
}

func (gcp *gcpRepository) DeleteFirewall(ctx context.Context, firewall *resource.Firewall) errors.Error {
	gcpFirewall := gcp.toGCPFirewall(ctx, firewall)
	if err := kubernetes.Client().Delete(ctx, gcpFirewall); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("firewall %s not found in namespace %s", firewall.IdentifierName.Firewall, firewall.Metadata.Namespace))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to delete firewall %s in namespace %s", firewall.IdentifierName.Firewall, firewall.Metadata.Namespace))
	}
	return errors.NoContent
}

func (gcp *gcpRepository) toModelFirewall(firewall *v1beta1.Firewall) (*resource.Firewall, errors.Error) {
	id := identifier.Firewall{}
	if err := id.FromLabels(firewall.Labels); !err.IsOk() {
		return nil, err
	}
	allow := resource.FirewallRuleList{}
	for _, a := range firewall.Spec.ForProvider.Allow {
		rule := resource.FirewallRule{
			Protocol: *a.Protocol,
		}
		var ports []string
		for _, p := range a.Ports {
			ports = append(ports, *p)
		}
		allow = append(allow, rule)
	}
	deny := resource.FirewallRuleList{}
	for _, a := range firewall.Spec.ForProvider.Deny {
		rule := resource.FirewallRule{
			Protocol: *a.Protocol,
		}
		var ports []string
		for _, p := range a.Ports {
			ports = append(ports, *p)
		}
		deny = append(deny, rule)
	}

	return &resource.Firewall{
		Metadata: metadata.Metadata{
			Managed:   firewall.Spec.ResourceSpec.DeletionPolicy == v1.DeletionDelete,
			Namespace: firewall.ObjectMeta.Namespace,
		},
		IdentifierID: id,
		IdentifierName: identifier.Firewall{
			Network:  *firewall.Spec.ForProvider.Network,
			VPC:      *firewall.Spec.ForProvider.Project,
			Provider: firewall.Spec.ResourceSpec.ProviderConfigReference.Name,
			Firewall: firewall.ObjectMeta.Annotations[crossplane.ExternalNameAnnotationKey],
		},

		Allow: allow,
		Deny:  deny,
	}, errors.OK
}

func (gcp *gcpRepository) toGCPFirewall(ctx context.Context, firewall *resource.Firewall) *v1beta1.Firewall {
	allow := &[]v1beta1.AllowParameters{}
	for _, a := range firewall.Allow {
		rule := v1beta1.AllowParameters{
			Protocol: &a.Protocol,
		}
		var ports []*string
		for _, p := range a.Ports {
			ports = append(ports, &p)
		}
		*allow = append(*allow, rule)
	}
	deny := &[]v1beta1.DenyParameters{}
	for _, a := range firewall.Deny {
		rule := v1beta1.DenyParameters{
			Protocol: &a.Protocol,
		}
		var ports []*string
		for _, p := range a.Ports {
			ports = append(ports, &p)
		}
		*deny = append(*deny, rule)
	}
	return &v1beta1.Firewall{
		ObjectMeta: metav1.ObjectMeta{
			Name:        firewall.IdentifierID.Firewall,
			Namespace:   firewall.Metadata.Namespace,
			Labels:      lo.Assign(crossplane.GetBaseLabels(ctx.Value(context.ProjectIDKey).(string)), firewall.IdentifierID.ToLabels()),
			Annotations: crossplane.GetAnnotations(firewall.Metadata.Managed, firewall.IdentifierName.Firewall),
		},
		Spec: v1beta1.FirewallSpec{
			ForProvider: v1beta1.FirewallParameters{
				Network: &firewall.IdentifierName.Network,
				Project: &firewall.IdentifierName.VPC,
				Allow:   *allow,
				Deny:    *deny,
			},
		},
	}

}

func (gcp *gcpRepository) toModelFirewallCollection(firewallList *v1beta1.FirewallList) (*resource.FirewallCollection, errors.Error) {
	var items resource.FirewallCollection
	for _, item := range firewallList.Items {
		firewall, err := gcp.toModelFirewall(&item)
		if !err.IsOk() {
			return nil, err
		}
		items[item.ObjectMeta.Annotations[crossplane.ExternalNameAnnotationKey]] = *firewall
	}
	return &items, errors.OK
}
