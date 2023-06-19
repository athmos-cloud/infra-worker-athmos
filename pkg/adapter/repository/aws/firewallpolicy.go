package aws

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/crossplane"
	model "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/network"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	v1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/samber/lo"
	"github.com/upbound/provider-aws/apis/networkfirewall/v1beta1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (aws *awsRepository) _createFirewallPolicy(ctx context.Context, firewall *model.Firewall, region *string) errors.Error {
	awsFirewallPolicy, fpErr := aws._toAWSFirewallPolicy(ctx, firewall, region)
	if !fpErr.IsOk() {
		return fpErr
	}

	if err := kubernetes.Client().Client.Create(ctx, awsFirewallPolicy); err != nil {
		if k8serrors.IsAlreadyExists(err) {
			return errors.Conflict.WithMessage(fmt.Sprintf("firewall policy %s already exists", awsFirewallPolicy.Name))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to create firewall policy %s", awsFirewallPolicy.Name))
	}

	return errors.Created
}

func (aws *awsRepository) _deleteFirewallPolicy(ctx context.Context, firewall *model.Firewall, region *string) errors.Error {
	awsFirewallPolicy, fpErr := aws._toAWSFirewallPolicy(ctx, firewall, region)
	if !fpErr.IsOk() {
		return fpErr
	}

	if err := kubernetes.Client().Client.Delete(ctx, awsFirewallPolicy); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("firewall policy %s not found", awsFirewallPolicy.Name))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to delete firewall policy %s", awsFirewallPolicy.Name))
	}

	return errors.NoContent
}

func (aws *awsRepository) _toAWSFirewallPolicy(ctx context.Context, firewall *model.Firewall, region *string) (*v1beta1.FirewallPolicy, errors.Error) {
	resLabels := lo.Assign(crossplane.GetBaseLabels(ctx.Value(context.ProjectIDKey).(string)), firewall.IdentifierID.ToIDLabels(), firewall.IdentifierName.ToNameLabels())

	priority := float64(1)
	var srgRef []v1beta1.StatelessRuleGroupReferenceParameters
	srgRef = append(srgRef, v1beta1.StatelessRuleGroupReferenceParameters{
		Priority: &priority,
		ResourceArnSelector: &v1.Selector{
			MatchLabels: resLabels,
		},
	})

	defaultAction := "aws:drop"
	var fpfPolicy []v1beta1.FirewallPolicyFirewallPolicyParameters
	fpfPolicy = append(fpfPolicy, v1beta1.FirewallPolicyFirewallPolicyParameters{
		StatelessDefaultActions:         []*string{&defaultAction},
		StatelessFragmentDefaultActions: []*string{&defaultAction},
		StatelessRuleGroupReference:     srgRef,
	})

	return &v1beta1.FirewallPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf("%s-policy", firewall.IdentifierID.Firewall),
			Labels:      resLabels,
			Annotations: crossplane.GetAnnotations(firewall.Metadata.Managed, firewall.IdentifierID.Firewall),
		},
		Spec: v1beta1.FirewallPolicySpec{
			ResourceSpec: v1.ResourceSpec{
				DeletionPolicy: crossplane.GetDeletionPolicy(firewall.Metadata.Managed),
				ProviderConfigReference: &v1.Reference{
					Name: firewall.IdentifierID.Provider,
				},
			},
			ForProvider: v1beta1.FirewallPolicyParameters{
				FirewallPolicy: fpfPolicy,
				Region:         region,
			},
		},
	}, errors.OK
}
