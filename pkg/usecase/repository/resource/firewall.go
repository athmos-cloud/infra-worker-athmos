package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	"sync"
)

type FirewallChannel struct {
	WaitGroup    *sync.WaitGroup
	Channel      chan *resource.FirewallCollection
	ErrorChannel chan errors.Error
}

type Firewall interface {
	FindFirewall(context.Context, option.Option) (*resource.Firewall, errors.Error)
	FindAllFirewalls(context.Context, option.Option) (*resource.FirewallCollection, errors.Error)
	FindAllRecursiveFirewalls(context.Context, option.Option, *FirewallChannel)
	CreateFirewall(context.Context, *resource.Firewall) errors.Error
	UpdateFirewall(context.Context, *resource.Firewall) errors.Error
	DeleteFirewall(context.Context, *resource.Firewall) errors.Error
	FirewallExists(context.Context, option.Option) (bool, errors.Error)
}
