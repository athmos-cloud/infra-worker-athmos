package resource

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	network2 "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/network"
	output "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/output/resource"
	"github.com/gin-gonic/gin"
)

type firewall struct{}

func NewFirewallPresenter() output.FirewallPort {
	return &firewall{}
}

func (n *firewall) Render(ctx context.Context, firewall *network2.Firewall) {
	resp := dto.GetFirewallResponse{
		ProjectID: ctx.Value(context.ProjectIDKey).(string),
		Payload:   *firewall,
	}
	ctx.JSON(200, gin.H{"payload": resp})
}

func (n *firewall) RenderCreate(ctx context.Context, firewall *network2.Firewall) {
	resp := dto.CreateFirewallResponse{
		ProjectID: ctx.Value(context.ProjectIDKey).(string),
		Payload:   *firewall,
	}
	ctx.JSON(201, gin.H{"payload": resp})
}

func (n *firewall) RenderUpdate(ctx context.Context, firewall *network2.Firewall) {
	ctx.JSON(204, gin.H{"message": fmt.Sprintf("firewall %s updated", firewall.IdentifierID.Firewall)})
}

func (n *firewall) RenderDelete(ctx context.Context, firewall *network2.Firewall) {
	ctx.JSON(204, gin.H{"message": fmt.Sprintf("firewall %s deleted", firewall.IdentifierID.Firewall)})
}
