package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/network"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	"sync"
)

type NetworkChannel struct {
	WaitGroup    *sync.WaitGroup
	Channel      chan *network.NetworkCollection
	ErrorChannel chan errors.Error
}

type Network interface {
	FindNetwork(context.Context, option.Option) (*network.Network, errors.Error)
	FindAllNetworks(context.Context, option.Option) (*network.NetworkCollection, errors.Error)
	FindAllRecursiveNetworks(context.Context, option.Option, *NetworkChannel) (*network.NetworkCollection, errors.Error)
	CreateNetwork(context.Context, *network.Network) errors.Error
	UpdateNetwork(context.Context, *network.Network) errors.Error
	DeleteNetwork(context.Context, *network.Network) errors.Error
	DeleteNetworkCascade(context.Context, *network.Network) errors.Error
	NetworkExists(context.Context, *network.Network) (bool, errors.Error)
}
