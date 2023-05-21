package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/error"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/validator/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	output "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/output/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase"
)

type Network interface {
	GetNetwork(context.Context)
	CreateNetwork(context.Context)
	UpdateNetwork(context.Context)
	DeleteNetwork(context.Context)
}

type networkController struct {
	networkUseCase usecase.Network
	networkOutput  output.NetworkPort
}

func NewNetworkController(networkUseCase usecase.Network, networkOutput output.NetworkPort) Network {
	return &networkController{networkUseCase: networkUseCase, networkOutput: networkOutput}
}

func (nc *networkController) GetNetwork(ctx context.Context) {
	resourceValidator.GetNetwork(ctx)
	network := &resource.Network{}
	if err := nc.networkUseCase.Get(ctx, network); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	} else {
		nc.networkOutput.Render(ctx, network)
	}
}

func (nc *networkController) CreateNetwork(ctx context.Context) {
	resourceValidator.CreateNetwork(ctx)
	network := &resource.Network{}
	if err := nc.networkUseCase.Create(ctx, network); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	} else {
		nc.networkOutput.RenderCreate(ctx, network)
	}
}

func (nc *networkController) UpdateNetwork(ctx context.Context) {
	resourceValidator.UpdateNetwork(ctx)
	network := &resource.Network{}
	if err := nc.networkUseCase.Update(ctx, network); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	} else {
		nc.networkOutput.RenderUpdate(ctx, network)
	}
}

func (nc *networkController) DeleteNetwork(ctx context.Context) {
	resourceValidator.DeleteNetwork(ctx)
	network := &resource.Network{}
	if err := nc.networkUseCase.Delete(ctx, network); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	} else {
		nc.networkOutput.Render(ctx, network)
	}
}
