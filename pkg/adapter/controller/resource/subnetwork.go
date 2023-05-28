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

func NewSubnetworkController(subnetworkUseCase usecase.Subnetwork, subnetworkOutput output.SubnetworkPort) Subnetwork {
	return &subnetworkController{subnetworkUseCase: subnetworkUseCase, subnetworkOutput: subnetworkOutput}
}

func (sc *subnetworkController) GetSubnetwork(ctx context.Context) {
	if err := resourceValidator.GetSubnetwork(ctx); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	}
	subnetwork := &resource.Subnetwork{}
	if err := sc.subnetworkUseCase.Get(ctx, subnetwork); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	} else {
		sc.subnetworkOutput.Render(ctx, subnetwork)
	}
}

func (sc *subnetworkController) CreateSubnetwork(ctx context.Context) {
	if err := resourceValidator.CreateSubnetwork(ctx); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	}
	subnetwork := &resource.Subnetwork{}
	if err := sc.subnetworkUseCase.Create(ctx, subnetwork); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	} else {
		sc.subnetworkOutput.RenderCreate(ctx, subnetwork)
	}
}

func (sc *subnetworkController) UpdateSubnetwork(ctx context.Context) {
	if err := resourceValidator.UpdateSubnetwork(ctx); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	}
	subnetwork := &resource.Subnetwork{}
	if err := sc.subnetworkUseCase.Update(ctx, subnetwork); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	} else {
		sc.subnetworkOutput.RenderUpdate(ctx, subnetwork)
	}
}

func (sc *subnetworkController) DeleteSubnetwork(ctx context.Context) {
	if err := resourceValidator.DeleteSubnetwork(ctx); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	}
	subnetwork := &resource.Subnetwork{}
	if err := sc.subnetworkUseCase.Delete(ctx, subnetwork); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	} else {
		sc.subnetworkOutput.RenderDelete(ctx, subnetwork)
	}
}
