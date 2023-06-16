package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/error"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/validator/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/network"
	output "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/output/resource"
	resource2 "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/resource"
)

type Network interface {
	GetNetwork(context.Context)
	CreateNetwork(context.Context)
	UpdateNetwork(context.Context)
	DeleteNetwork(context.Context)
}

type networkController struct {
	networkUseCase resource2.Network
	networkOutput  output.NetworkPort
}

func NewNetworkController(networkUseCase resource2.Network, networkOutput output.NetworkPort) Network {
	return &networkController{networkUseCase: networkUseCase, networkOutput: networkOutput}
}

func (nc *networkController) GetNetwork(ctx context.Context) {
	if err := resourceValidator.GetNetwork(ctx); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	}
	net := &network.Network{}
	if err := nc.networkUseCase.Get(ctx, net); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	} else {
		nc.networkOutput.Render(ctx, net)
	}
}

func (nc *networkController) CreateNetwork(ctx context.Context) {
	if err := resourceValidator.CreateNetwork(ctx); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	}
	net := &network.Network{}
	if err := nc.networkUseCase.Create(ctx, net); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	} else {
		nc.networkOutput.RenderCreate(ctx, net)
	}
}

func (nc *networkController) UpdateNetwork(ctx context.Context) {
	if err := resourceValidator.UpdateNetwork(ctx); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	}
	network := &network.Network{}
	if err := nc.networkUseCase.Update(ctx, network); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	} else {
		nc.networkOutput.RenderUpdate(ctx, network)
	}
}

func (nc *networkController) DeleteNetwork(ctx context.Context) {
	if err := resourceValidator.DeleteNetwork(ctx); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	}
	network := &network.Network{}
	if err := nc.networkUseCase.Delete(ctx, network); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	} else {
		nc.networkOutput.Render(ctx, network)
	}
}
