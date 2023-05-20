package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase"
)

type Provider interface {
	GetProvider(context.Context)
	ListProviders(context.Context)
	CreateProvider(context.Context)
	UpdateProvider(context.Context)
	DeleteProvider(context.Context)
}

type providerController struct {
	networkUseCase usecase.Provider
}

func NewProviderController(networkUseCase usecase.Provider) Provider {
	return &providerController{networkUseCase: networkUseCase}
}

func (pc *providerController) GetProvider(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

func (pc *providerController) ListProviders(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

func (pc *providerController) CreateProvider(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

func (pc *providerController) UpdateProvider(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

func (pc *providerController) DeleteProvider(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}
