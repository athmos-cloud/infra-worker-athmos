package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase"
)

type Subnetwork interface {
	GetSubnetwork(context.Context)
	ListSubnetworks(context.Context)
	CreateSubnetwork(context.Context)
	UpdateSubnetwork(context.Context)
	DeleteSubnetwork(context.Context)
}

type subnetworkController struct {
	networkUseCase usecase.Subnetwork
}

func NewSubnetworkController(networkUseCase usecase.Subnetwork) Subnetwork {
	return &subnetworkController{networkUseCase: networkUseCase}
}

func (sc *subnetworkController) GetSubnetwork(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

func (sc *subnetworkController) ListSubnetworks(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

func (sc *subnetworkController) CreateSubnetwork(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

func (sc *subnetworkController) UpdateSubnetwork(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

func (sc *subnetworkController) DeleteSubnetwork(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}
