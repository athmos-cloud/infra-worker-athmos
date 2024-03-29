package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
)

type Provider interface {
	FindProvider(context.Context, option.Option) (*resource.Provider, errors.Error)
	FindAllProviders(context.Context, option.Option) (*resource.ProviderCollection, errors.Error)
	FindProviderStack(context.Context, option.Option) (*resource.Provider, errors.Error)
	CreateProvider(context.Context, *resource.Provider) errors.Error
	UpdateProvider(context.Context, *resource.Provider) errors.Error
	DeleteProvider(context.Context, *resource.Provider) errors.Error
	DeleteProviderCascade(context.Context, *resource.Provider) errors.Error
	ProviderExists(context.Context, identifier.Provider) (bool, errors.Error)
}
