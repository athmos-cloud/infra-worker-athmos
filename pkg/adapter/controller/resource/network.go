package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase"
)

type Network interface {
	GetNetwork(context.Context)
	ListNetworks(context.Context)
	CreateNetwork(context.Context)
	UpdateNetwork(context.Context)
	DeleteNetwork(context.Context)
}

type networkController struct {
	networkUseCase usecase.Network
}

func (nc *networkController) GetNetwork(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

func (nc *networkController) ListNetworks(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

func (nc *networkController) CreateNetwork(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

func (nc *networkController) UpdateNetwork(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

func (nc *networkController) DeleteNetwork(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

func NewNetworkController(networkUseCase usecase.Network) Network {
	return &networkController{networkUseCase: networkUseCase}
}
