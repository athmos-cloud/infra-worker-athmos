package gcpRepo

import (
	"context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
)

type Network interface {
	Find(context.Context, option.Option) *resource.Network
	FindAll(context.Context, option.Option) []*resource.Network
	Create(context.Context, *resource.Network) *resource.Network
	Update(context.Context, *resource.Network) *resource.Network
	Delete(context.Context, *resource.Network)
}
