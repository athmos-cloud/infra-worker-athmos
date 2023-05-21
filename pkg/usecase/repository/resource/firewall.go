package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
)

type Firewall interface {
	FindFirewall(context.Context, option.Option) (*resource.Firewall, errors.Error)
	FindAllFirewalls(context.Context, option.Option) (*resource.FirewallCollection, errors.Error)
	FindAllRecursiveFirewalls(context.Context, option.Option) (*resource.FirewallCollection, errors.Error)
	CreateFirewall(context.Context, *resource.Firewall) errors.Error
	UpdateFirewall(context.Context, *resource.Firewall) errors.Error
	DeleteFirewall(context.Context, *resource.Firewall) errors.Error
}
