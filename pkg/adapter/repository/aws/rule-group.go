package aws

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	model "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/upbound/provider-aws/apis/networkfirewall/v1beta1"
)

func (aws *awsRepository) _createRuleGroup(ctx context.Context, vm *model.Firewall) errors.Error {

	return errors.Created
}

func (aws *awsRepository) _getRuleGroup(ctx context.Context, vm *model.Firewall) errors.Error {

	return errors.Created
}

func (aws *awsRepository) _updateRuleGroup(ctx context.Context, vm *model.Firewall) errors.Error {

	return errors.Created
}

func (aws *awsRepository) _deleteRuleGroup(ctx context.Context, vm *model.Firewall) errors.Error {

	return errors.Created
}

func (aws *awsRepository) _getAwsRuleGroup(vm *model.Firewall) *v1beta1.RuleGroup {
	//model.FirewallRule{}
	return &v1beta1.RuleGroup{}
}
