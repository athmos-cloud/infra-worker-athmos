package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/network"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	"sync"
)

type FirewallChannel struct {
	WaitGroup    *sync.WaitGroup
	Channel      chan *network.FirewallCollection
	ErrorChannel chan errors.Error
}

type Firewall interface {
	FindFirewall(context.Context, option.Option) (*network.Firewall, errors.Error)
	FindAllFirewalls(context.Context, option.Option) (*network.FirewallCollection, errors.Error)
	FindAllRecursiveFirewalls(context.Context, option.Option, *FirewallChannel)
	CreateFirewall(context.Context, *network.Firewall) errors.Error
	UpdateFirewall(context.Context, *network.Firewall) errors.Error
	DeleteFirewall(context.Context, *network.Firewall) errors.Error
	FirewallExists(context.Context, *network.Firewall) (bool, errors.Error)
}
