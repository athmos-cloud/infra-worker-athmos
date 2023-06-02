package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	"sync"
)

type NetworkChannel struct {
	WaitGroup    *sync.WaitGroup
	Channel      chan *resource.NetworkCollection
	ErrorChannel chan errors.Error
}

type Network interface {
	FindNetwork(context.Context, option.Option) (*resource.Network, errors.Error)
	FindAllNetworks(context.Context, option.Option) (*resource.NetworkCollection, errors.Error)
	FindAllRecursiveNetworks(context.Context, option.Option, *NetworkChannel) (*resource.NetworkCollection, errors.Error)
	CreateNetwork(context.Context, *resource.Network) errors.Error
	UpdateNetwork(context.Context, *resource.Network) errors.Error
	DeleteNetwork(context.Context, *resource.Network) errors.Error
	DeleteNetworkCascade(context.Context, *resource.Network) errors.Error
	NetworkExists(context.Context, *resource.Network) (bool, errors.Error)
}
