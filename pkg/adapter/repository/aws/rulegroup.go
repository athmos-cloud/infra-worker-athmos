package aws

import (
	"fmt"
	"strconv"

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

func (aws *awsRepository) _getRuleGroup(ctx context.Context, ruleGroupName *string) (*v1beta1.RuleGroup, errors.Error) {

	awsRuleGroup := &v1beta1.RuleGroup{}
	if err := kubernetes.Client().Client.Get(ctx, client.ObjectKey{Name: *ruleGroupName}, awsRuleGroup); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, errors.NotFound.WithMessage(fmt.Sprintf("rule group %s not found", *ruleGroupName))
		}
		return nil, errors.KubernetesError.WithMessage(fmt.Sprintf("unable to get rule group %s", *ruleGroupName))
	}
	return awsRuleGroup, errors.OK
}

func (aws *awsRepository) _updateRuleGroup(ctx context.Context, firewall *model.Firewall, region *string) errors.Error {
	name := fmt.Sprintf("%s-rule-group", firewall.IdentifierID.Firewall)
	existingAwsRuleGroup, err := aws._getRuleGroup(ctx, &name)
	if !err.IsOk() {
		return err
	}

	awsRuleGroup, err := aws._toAwsRuleGroup(ctx, firewall, region)
	if !err.IsOk() {
		return err
	}

	existingAwsRuleGroup.Labels = awsRuleGroup.Labels
	existingAwsRuleGroup.Spec = awsRuleGroup.Spec
	if err := kubernetes.Client().Client.Update(ctx, existingAwsRuleGroup); err != nil {
		if k8serrors.IsNotFound(err) {
			return errors.NotFound.WithMessage(fmt.Sprintf("rule group %s not found", awsRuleGroup.Name))
		}
		return errors.KubernetesError.WithMessage(fmt.Sprintf("unable to update rule group %s", awsRuleGroup.Name))
	}
	return errors.NoContent
}

func (aws *awsRepository) _deleteRuleGroup(ctx context.Context, firewall *model.Firewall, region *string) errors.Error {
	awsRuleGroup, err := aws._toAwsRuleGroup(ctx, firewall, region)
	if !err.IsOk() {
		return err
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

	priority := float64(1)
	ruleType := "STATELESS"
	ruleGroupCapacity := float64(5000)

	allowAction := "aws:pass"
	denyAction := "aws:drop"

	allowParameters, errA := _toAwsMatchAttributesParameters(firewall.Allow)
	if !errA.IsOk() {
		return nil, errA
	}

	denyParameters, errD := _toAwsMatchAttributesParameters(firewall.Deny)
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
				Capacity: &ruleGroupCapacity,
				Name:     &firewall.IdentifierID.Provider,
				Region:   region,
				RuleGroup: []v1beta1.RuleGroupRuleGroupParameters{
					{
						RulesSource: []v1beta1.RulesSourceParameters{
							{
								StatelessRulesAndCustomActions: []v1beta1.StatelessRulesAndCustomActionsParameters{
									{
										StatelessRule: []v1beta1.StatelessRuleParameters{
											{
												Priority: &priority,
												RuleDefinition: []v1beta1.RuleDefinitionParameters{
													{
														Actions:         []*string{&allowAction},
														MatchAttributes: allowParameters,
													},
													{
														Actions:         []*string{&denyAction},
														MatchAttributes: denyParameters,
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
				Type: &ruleType,
			},
		},
	}, errors.OK
}

func (aws *awsRepository) _toFirewallRuleList(ruleGroup *v1beta1.RuleGroup) (*model.FirewallRuleList, *model.FirewallRuleList, errors.Error) {
	ruleDefinitions, errRD := _extractRuleDefinitionParametersArray(ruleGroup)
	if !errRD.IsOk() {
		return nil, nil, errRD
	}

	return _toAllowAndDenyRuleList(ruleDefinitions)
}

func _extractRuleDefinitionParametersArray(ruleGroup *v1beta1.RuleGroup) (*[]v1beta1.RuleDefinitionParameters, errors.Error) {
	internalError := errors.InternalError.WithMessage("Invalid rule group parameters received from kubernetes.")

	params := (*ruleGroup).Spec.ForProvider.RuleGroup
	if len(params) == 0 {
		return nil, internalError
	}

	sources := params[0].RulesSource
	if len(sources) == 0 {
		return nil, internalError
	}

	statelessRCA := sources[0].StatelessRulesAndCustomActions
	if len(statelessRCA) == 0 {
		return nil, internalError
	}

	stateless := statelessRCA[0].StatelessRule
	if len(stateless) == 0 {
		return nil, internalError
	}

	ruleDefinitions := stateless[0].RuleDefinition
	if len(ruleDefinitions) != 2 {
		return nil, internalError
	}

	return &ruleDefinitions, errors.OK
}

func _toAllowAndDenyRuleList(ruleDefinitions *[]v1beta1.RuleDefinitionParameters) (*model.FirewallRuleList, *model.FirewallRuleList, errors.Error) {
	var allow model.FirewallRuleList
	var deny model.FirewallRuleList

	for _, ruleDefinition := range *ruleDefinitions {
		if len(ruleDefinition.Actions) == 0 {
			return nil, nil, errors.InternalError.WithMessage("Invalid rule definition received from kubernetes.")
		}

		action := *ruleDefinition.Actions[0]

		switch action {
		case "aws:pass":
			a, errA := _toModelRuleLists(&ruleDefinition.MatchAttributes)
			allow = *a
			if !errA.IsOk() {
				return nil, nil, errA
			}
			break
		case "aws:drop":
			d, errD := _toModelRuleLists(&ruleDefinition.MatchAttributes)
			deny = *d
			if !errD.IsOk() {
				return nil, nil, errD
			}
			break
		default:
			return nil, nil, errors.InternalError.WithMessage("Invalid rule definition received from kubernetes.")
		}
	}

	return &allow, &deny, errors.OK
}

func _toModelRuleLists(awsMatchAttributes *[]v1beta1.MatchAttributesParameters) (*model.FirewallRuleList, errors.Error) {
	var ruleList model.FirewallRuleList

	for _, match := range *awsMatchAttributes {
		if len(match.Protocols) == 0 {
			return nil, errors.InternalError.WithMessage("Invalid port rule received from kubernetes.")
		}
		strProtocol, errP := _protocolIANAToString(*match.Protocols[0])
		if !errP.IsOk() {
			return nil, errors.InternalError.WithMessage("Invalid port received from kubernetes.")
		}

		var strPorts []string
		for _, port := range match.DestinationPort {
			strPorts = append(strPorts, fmt.Sprintf("%.0f", *port.FromPort))
		}

		ruleList = append(ruleList, model.FirewallRule{
			Protocol: strProtocol,
			Ports:    strPorts,
		})
	}

	return &ruleList, errors.OK
}

func _toAwsMatchAttributesParameters(ruleList model.FirewallRuleList) ([]v1beta1.MatchAttributesParameters, errors.Error) {
	res := make([]v1beta1.MatchAttributesParameters, 0)

	for _, rule := range ruleList {
		protocolf, errP := _protocolIANAFromString(rule.Protocol)
		if !errP.IsOk() {
			return res, errP
		}

		portsf, errPorts := _floatSliceFromStringPorts(rule.Ports)
		if !errPorts.IsOk() {
			return res, errPorts
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
	return res, errors.OK
}

func _floatSliceFromStringPorts(ports []string) ([]*float64, errors.Error) {
	var res []*float64
	for _, port := range ports {
		portI, err := strconv.Atoi(port)
		if err != nil {
			return nil, errors.BadRequest.WithMessage(fmt.Sprintf("%s is not a valid port", port))
		}
		portF := float64(portI)
		res = append(res, &portF)
	}
	return res, errors.OK
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

func _protocolIANAToString(protocol float64) (string, errors.Error) {
	switch protocol {
	case float64(6):
		return "tcp", errors.OK
	case float64(17):
		return "udp", errors.OK
	case float64(1):
		return "icmp", errors.OK
	case float64(50):
		return "esp", errors.OK
	case float64(51):
		return "ah", errors.OK
	case float64(132):
		return "sctp", errors.OK
	case float64(94):
		return "ipip", errors.OK
	}
	return "", errors.InvalidArgument.WithMessage(fmt.Sprintf("protocol %.0f is not handled", protocol))
}
