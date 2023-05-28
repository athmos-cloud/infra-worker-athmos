package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	errorCtrl "github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/error"
	resourceValidator "github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/validator/resource"
	model "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	output "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/output/resource"
	usecase "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/resource"
)

type Provider interface {
	GetProvider(context.Context)
	GetStack(context.Context)
	ListProviders(context.Context)
	CreateProvider(context.Context)
	UpdateProvider(context.Context)
	DeleteProvider(context.Context)
}

type providerController struct {
	providerUseCase usecase.Provider
	providerOutput  output.ProviderPort
}

func NewProviderController(providerUseCase usecase.Provider, providerOutput output.ProviderPort) Provider {
	return &providerController{providerUseCase: providerUseCase, providerOutput: providerOutput}
}

func (pc *providerController) GetProvider(ctx context.Context) {
	if err := resourceValidator.GetProvider(ctx); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	}
	provider := &model.Provider{}
	if err := pc.providerUseCase.Get(ctx, provider); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	} else {
		pc.providerOutput.Render(ctx, provider)
	}
}

func (pc *providerController) GetStack(ctx context.Context) {
	if err := resourceValidator.Stack(ctx); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	}
	stack := &model.Provider{}
	if err := pc.providerUseCase.GetRecursively(ctx, stack); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	} else {
		pc.providerOutput.Render(ctx, stack)
	}
}

func (pc *providerController) ListProviders(ctx context.Context) {
	if err := resourceValidator.ListProviders(ctx); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	}
	providers := &model.ProviderCollection{}
	if err := pc.providerUseCase.List(ctx, providers); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	} else {
		pc.providerOutput.RenderAll(ctx, providers)
	}
}

func (pc *providerController) CreateProvider(ctx context.Context) {
	if err := resourceValidator.CreateProvider(ctx); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	}
	provider := &model.Provider{}
	if err := pc.providerUseCase.Create(ctx, provider); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	} else {
		pc.providerOutput.RenderCreate(ctx, provider)
	}
}

func (pc *providerController) UpdateProvider(ctx context.Context) {
	if err := resourceValidator.UpdateProvider(ctx); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	}
	provider := &model.Provider{}
	if err := pc.providerUseCase.Update(ctx, provider); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	} else {
		pc.providerOutput.RenderUpdate(ctx, provider)
	}
}

func (pc *providerController) DeleteProvider(ctx context.Context) {
	if err := resourceValidator.DeleteProvider(ctx); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	}
	provider := &model.Provider{}
	if err := pc.providerUseCase.Delete(ctx, provider); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	} else {
		pc.providerOutput.RenderDelete(ctx, provider)
	}
}
