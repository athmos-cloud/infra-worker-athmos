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
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
)

func (aws *awsRepository) _createRuleGroup(ctx context.Context, firewall *model.Firewall, region *string) errors.Error {
	awsRuleGroup, rgErr := aws._toAwsRuleGroup(ctx, firewall, region)
	if !rgErr.IsOk() {
		return rgErr
	}
	if err := kubernetes.Client().Client.Create(ctx, awsRuleGroup); err != nil {
		if k8serrors.IsAlreadyExists(err) {
			return errors.Conflict.WithMessage(fmt.Sprintf("rule group %s already exists", awsRuleGroup.Name))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to create rule group %s", awsRuleGroup.Name))
	}

	return errors.Created
}

func (aws *awsRepository) _getRuleGroup(ctx context.Context, firewall *model.Firewall) (*v1beta1.RuleGroup, errors.Error) {
	name := fmt.Sprintf("%s-rule-group", firewall.IdentifierID.Firewall)

	awsRuleGroup := &v1beta1.RuleGroup{}
	if err := kubernetes.Client().Client.Get(ctx, client.ObjectKey{Name: name}, awsRuleGroup); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, errors.NotFound.WithMessage(fmt.Sprintf("rule group %s not found", name))
		}
		return nil, errors.KubernetesError.WithMessage(fmt.Sprintf("unable to get rule group %s", name))
	}
	return awsRuleGroup, errors.OK
}

func (aws *awsRepository) _updateRuleGroup(ctx context.Context, firewall *model.Firewall, region *string) errors.Error {
	awsRuleGroup, rgErr := aws._toAwsRuleGroup(ctx, firewall, region)
	if !rgErr.IsOk() {
		return rgErr
	}

	if err := kubernetes.Client().Client.Update(ctx, awsRuleGroup); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("rule group %s not found", awsRuleGroup.Name))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to update rule group %s", awsRuleGroup.Name))
	}
	return errors.NoContent
}

func (aws *awsRepository) _deleteRuleGroup(ctx context.Context, firewall *model.Firewall, region *string) errors.Error {
	awsRuleGroup, rgErr := aws._toAwsRuleGroup(ctx, firewall, region)
	if !rgErr.IsOk() {
		return rgErr
	}

	if err := kubernetes.Client().Client.Delete(ctx, awsRuleGroup); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("rule group %s not found", awsRuleGroup.Name))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to delete rule group %s", awsRuleGroup.Name))
	}

	return errors.NoContent
}

func (aws *awsRepository) _toAwsRuleGroup(ctx context.Context, firewall *model.Firewall, region *string) (*v1beta1.RuleGroup, errors.Error) {
	resLabels := lo.Assign(crossplane.GetBaseLabels(
		ctx.Value(context.ProjectIDKey).(string)),
		firewall.IdentifierID.ToIDLabels(),
		firewall.IdentifierName.ToNameLabels())

	allowAction := "aws:pass"
	denyAction := "aws:drop"

	allowParameters, errA := aws._toAwsRuleList(firewall.Allow)
	if !errA.IsOk() {
		return nil, errA
	}

	denyParameters, errD := aws._toAwsRuleList(firewall.Deny)
	if !errD.IsOk() {
		return nil, errD
	}

	return &v1beta1.RuleGroup{
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf("%s-rule-group", firewall.IdentifierID.Firewall),
			Labels:      resLabels,
			Annotations: crossplane.GetAnnotations(firewall.Metadata.Managed, firewall.IdentifierID.Firewall),
		},
		Spec: v1beta1.RuleGroupSpec{
			ResourceSpec: v1.ResourceSpec{
				DeletionPolicy: crossplane.GetDeletionPolicy(firewall.Metadata.Managed),
				ProviderConfigReference: &v1.Reference{
					Name: firewall.IdentifierID.Provider,
				},
			},
			ForProvider: v1beta1.RuleGroupParameters{
				Region: region,
				RuleGroup: []v1beta1.RuleGroupRuleGroupParameters{
					{
						RulesSource: []v1beta1.RulesSourceParameters{
							{
								StatelessRulesAndCustomActions: []v1beta1.StatelessRulesAndCustomActionsParameters{
									{
										StatelessRule: []v1beta1.StatelessRuleParameters{
											{
												RuleDefinition: []v1beta1.RuleDefinitionParameters{
													{
														Actions:         []*string{&allowAction},
														MatchAttributes: *allowParameters,
													},
													{
														Actions:         []*string{&denyAction},
														MatchAttributes: *denyParameters,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}, errors.OK
}

func (aws *awsRepository) _toAwsRuleList(ruleList model.FirewallRuleList) (*[]v1beta1.MatchAttributesParameters, errors.Error) {
	var res []v1beta1.MatchAttributesParameters
	toFloatTable := func(entry []string) ([]*float64, errors.Error) {
		var res []*float64
		for _, port := range entry {
			porti, err := strconv.Atoi(port)
			if err != nil {
				return nil, errors.BadRequest.WithMessage(fmt.Sprintf("%s is not a valid port", port))
			}
			portf := float64(porti)
			res = append(res, &portf)
		}
		return res, errors.OK
	}
	for _, rule := range ruleList {
		protocolf, errP := _protocolIANAFromString(rule.Protocol)
		if !errP.IsOk() {
			return nil, errP
		}

		portsf, errPorts := toFloatTable(rule.Ports)
		if !errPorts.IsOk() {
			return nil, errPorts
		}

		newPort := v1beta1.MatchAttributesParameters{
			Protocols: []*float64{protocolf},
		}
		var ports []v1beta1.DestinationPortParameters
		for _, port := range portsf {
			ports = append(ports, v1beta1.DestinationPortParameters{
				FromPort: port,
				ToPort:   port,
			})
		}
		newPort.DestinationPort = ports
		res = append(res, newPort)
	}
	return &res, errors.OK
}

func _protocolIANAFromString(protocol string) (*float64, errors.Error) {
	switch protocol {
	case "tcp":
		protocolf := float64(6)
		return &protocolf, errors.OK
	case "udp":
		protocolf := float64(17)
		return &protocolf, errors.OK
	case "icmp":
		protocolf := float64(1)
		return &protocolf, errors.OK
	case "esp":
		protocolf := float64(50)
		return &protocolf, errors.OK
	case "ah":
		protocolf := float64(51)
		return &protocolf, errors.OK
	case "sctp":
		protocolf := float64(132)
		return &protocolf, errors.OK
	case "ipip":
		protocolf := float64(94)
		return &protocolf, errors.OK
	}
	return nil, errors.InvalidArgument.WithMessage(fmt.Sprintf("protocol %s is not handled", protocol))
}
