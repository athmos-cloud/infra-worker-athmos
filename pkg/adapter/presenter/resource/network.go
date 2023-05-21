package resource

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	output "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/output/resource"
	"github.com/gin-gonic/gin"
)

type network struct{}

func NewNetworkPresenter() output.NetworkPort {
	return &network{}
}

func (n *network) Render(ctx context.Context, network *resource.Network) {
	resp := dto.GetNetworkResponse{
		ProjectID: ctx.Value(context.ProjectIDKey).(string),
		Payload:   *network,
	}
	ctx.JSON(200, gin.H{"payload": resp})
}

func (n *network) RenderCreate(ctx context.Context, network *resource.Network) {
	resp := dto.CreateNetworkResponse{
		ProjectID: ctx.Value(context.ProjectIDKey).(string),
		Payload:   *network,
	}
	ctx.JSON(201, gin.H{"payload": resp})
}

func (n *network) RenderUpdate(ctx context.Context, network *resource.Network) {
	ctx.JSON(204, gin.H{"message": fmt.Sprintf("network %s updated", network.IdentifierID.Network)})
}

func (n *network) RenderDelete(ctx context.Context, network *resource.Network) {
	ctx.JSON(204, gin.H{"message": fmt.Sprintf("network %s deleted", network.IdentifierID.Network)})
}
