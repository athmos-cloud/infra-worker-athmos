package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	errorCtrl "github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/error"
	resourceValidator "github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/validator/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	output "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/output/resource"
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
	providerUseCase usecase.Provider
	providerOutput  output.ProviderPort
}

func NewProviderController(networkUseCase usecase.Provider, providerOutput output.ProviderPort) Provider {
	return &providerController{providerUseCase: networkUseCase, providerOutput: providerOutput}
}

func (pc *providerController) GetProvider(ctx context.Context) {
	resourceValidator.GetProvider(ctx)
	provider := &resource.Provider{}
	if err := pc.providerUseCase.Get(ctx, provider); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	} else {
		pc.providerOutput.Render(ctx, provider)
	}
}

func (pc *providerController) ListProviders(ctx context.Context) {
	resourceValidator.ListProviders(ctx)
	providers := &resource.ProviderCollection{}
	if err := pc.providerUseCase.List(ctx, providers); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	} else {
		pc.providerOutput.RenderAll(ctx, providers)
	}
}

func (pc *providerController) CreateProvider(ctx context.Context) {
	resourceValidator.CreateProvider(ctx)
	provider := &resource.Provider{}
	if err := pc.providerUseCase.Create(ctx, provider); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	} else {
		pc.providerOutput.RenderCreate(ctx, provider)
	}
}

func (pc *providerController) UpdateProvider(ctx context.Context) {
	resourceValidator.UpdateProvider(ctx)
	provider := &resource.Provider{}
	if err := pc.providerUseCase.Update(ctx, provider); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	} else {
		pc.providerOutput.RenderUpdate(ctx, provider)
	}
}

func (pc *providerController) DeleteProvider(ctx context.Context) {
	resourceValidator.DeleteProvider(ctx)
	provider := &resource.Provider{}
	if err := pc.providerUseCase.Delete(ctx, provider); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	} else {
		pc.providerOutput.RenderDelete(ctx, provider)
	}
}
