package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	errorCtrl "github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/error"
	resourceValidator "github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/validator/resource"
	model "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/instance"
	output "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/output/resource"
	usecase "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/resource"
)

type VM interface {
	GetVM(context.Context)
	CreateVM(context.Context)
	UpdateVM(context.Context)
	DeleteVM(context.Context)
}

type vmController struct {
	vmUseCase usecase.VM
	vmOutput  output.VMPort
}

func NewVMController(vmUseCase usecase.VM, vmOutput output.VMPort) VM {
	return &vmController{vmUseCase: vmUseCase, vmOutput: vmOutput}
}

func (vc *vmController) GetVM(ctx context.Context) {
	resourceValidator.GetVM(ctx)
	vm := &model.VM{}
	if err := vc.vmUseCase.Get(ctx, vm); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	} else {
		vc.vmOutput.Render(ctx, vm)
	}
}

func (vc *vmController) CreateVM(ctx context.Context) {
	resourceValidator.CreateVM(ctx)
	vm := &model.VM{}
	if err := vc.vmUseCase.Create(ctx, vm); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	} else {
		vc.vmOutput.RenderCreate(ctx, vm)
	}
}

func (vc *vmController) UpdateVM(ctx context.Context) {
	resourceValidator.UpdateVM(ctx)
	vm := &model.VM{}
	if err := vc.vmUseCase.Update(ctx, vm); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	} else {
		vc.vmOutput.RenderUpdate(ctx, vm)
	}
}

func (vc *vmController) DeleteVM(ctx context.Context) {
	resourceValidator.DeleteVM(ctx)
	vm := &model.VM{}
	if err := vc.vmUseCase.Delete(ctx, vm); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	} else {
		vc.vmOutput.RenderDelete(ctx, vm)
	}
}
