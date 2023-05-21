package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
)

type Network interface {
	FindNetwork(context.Context, option.Option) (*resource.Network, errors.Error)
	FindAllRecursiveNetworks(context.Context, option.Option) (*resource.NetworkCollection, errors.Error)
	CreateNetwork(context.Context, *resource.Network) errors.Error
	UpdateNetwork(context.Context, *resource.Network) errors.Error
	DeleteNetwork(context.Context, *resource.Network) errors.Error
}
