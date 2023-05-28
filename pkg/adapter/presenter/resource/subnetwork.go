package resource

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	output "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/output/resource"
	"github.com/gin-gonic/gin"
)

type subnetwork struct{}

func NewSubnetworkPresenter() output.SubnetworkPort {
	return &subnetwork{}
}

func (s *subnetwork) Render(ctx context.Context, subnetwork *resource.Subnetwork) {
	resp := &dto.GetSubnetworkResponse{
		ProjectID: ctx.Value(context.ProjectIDKey).(string),
		Payload:   *subnetwork,
	}
	ctx.JSON(200, gin.H{"payload": resp})
}

func (s *subnetwork) RenderCreate(ctx context.Context, subnetwork *resource.Subnetwork) {
	resp := &dto.CreateSubnetworkResponse{
		ProjectID: ctx.Value(context.ProjectIDKey).(string),
		Payload:   *subnetwork,
	}
	ctx.JSON(201, gin.H{"payload": resp})
}

func (s *subnetwork) RenderUpdate(ctx context.Context, subnetwork *resource.Subnetwork) {
	ctx.JSON(204, gin.H{"message": fmt.Sprintf("subnetwork %s updated", subnetwork.IdentifierID.Subnetwork)})
}

func (s *subnetwork) RenderDelete(ctx context.Context, subnetwork *resource.Subnetwork) {
	ctx.JSON(204, gin.H{"message": fmt.Sprintf("subnetwork %s deleted", subnetwork.IdentifierID.Subnetwork)})
}
