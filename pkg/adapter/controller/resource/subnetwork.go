package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	errorCtrl "github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/error"
	resourceValidator "github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/validator/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	output "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/output/resource"
	usecase "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/resource"
)

type Subnetwork interface {
	GetSubnetwork(context.Context)
	CreateSubnetwork(context.Context)
	UpdateSubnetwork(context.Context)
	DeleteSubnetwork(context.Context)
}

type subnetworkController struct {
	subnetworkUseCase usecase.Subnetwork
	subnetworkOutput  output.SubnetworkPort
}

func NewSubnetworkController(networkUseCase usecase.Subnetwork, subnetworkOutput output.SubnetworkPort) Subnetwork {
	return &subnetworkController{subnetworkUseCase: networkUseCase, subnetworkOutput: subnetworkOutput}
}

func (sc *subnetworkController) GetSubnetwork(ctx context.Context) {
	resourceValidator.GetSubnetwork(ctx)
	subnetwork := &resource.Subnetwork{}
	if err := sc.subnetworkUseCase.Get(ctx, subnetwork); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	} else {
		sc.subnetworkOutput.Render(ctx, subnetwork)
	}
}

func (sc *subnetworkController) CreateSubnetwork(ctx context.Context) {
	resourceValidator.CreateSubnetwork(ctx)
	subnetwork := &resource.Subnetwork{}
	if err := sc.subnetworkUseCase.Create(ctx, subnetwork); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	} else {
		sc.subnetworkOutput.RenderCreate(ctx, subnetwork)
	}
}

func (sc *subnetworkController) UpdateSubnetwork(ctx context.Context) {
	resourceValidator.UpdateSubnetwork(ctx)
	subnetwork := &resource.Subnetwork{}
	if err := sc.subnetworkUseCase.Update(ctx, subnetwork); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	} else {
		sc.subnetworkOutput.RenderUpdate(ctx, subnetwork)
	}
}

func (sc *subnetworkController) DeleteSubnetwork(ctx context.Context) {
	resourceValidator.DeleteSubnetwork(ctx)
	subnetwork := &resource.Subnetwork{}
	if err := sc.subnetworkUseCase.Delete(ctx, subnetwork); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	} else {
		sc.subnetworkOutput.RenderDelete(ctx, subnetwork)
	}
}
