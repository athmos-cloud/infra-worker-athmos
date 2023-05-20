package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
)

type Provider interface {
	FindProvider(context.Context, option.Option) (*resource.Provider, errors.Error)
	FindAllProviders(context.Context, option.Option) (*resource.ProviderCollection, errors.Error)
	CreateProvider(context.Context, *resource.Provider) errors.Error
	UpdateProvider(context.Context, *resource.Provider) errors.Error
	DeleteProvider(context.Context, *resource.Provider) errors.Error
}
