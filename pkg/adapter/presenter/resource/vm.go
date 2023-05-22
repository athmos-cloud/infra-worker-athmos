package resource

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	output "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/output/resource"
	"github.com/gin-gonic/gin"
)

type virtualMachine struct{}

func (v *virtualMachine) Render(ctx context.Context, vm *resource.VM) {
	resp := dto.GetVMResponse{
		ProjectID: ctx.Value(context.ProjectIDKey).(string),
		Payload:   *vm,
	}
	ctx.JSON(200, gin.H{"payload": resp})
}

func (v *virtualMachine) RenderCreate(ctx context.Context, vm *resource.VM) {
	resp := dto.CreateVMResponse{
		ProjectID: ctx.Value(context.ProjectIDKey).(string),
		Payload:   *vm,
	}
	ctx.JSON(201, gin.H{"payload": resp})
}

func (v *virtualMachine) RenderUpdate(ctx context.Context, vm *resource.VM) {
	ctx.JSON(204, gin.H{"message": fmt.Sprintf("vm %s updated", vm.IdentifierID.VM)})
}

func (v *virtualMachine) RenderDelete(ctx context.Context, vm *resource.VM) {
	ctx.JSON(204, gin.H{"message": fmt.Sprintf("vm %s deleted", vm.IdentifierID.VM)})
}

func NewVMPresenter() output.VMPort {
	return &virtualMachine{}
}
