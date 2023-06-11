package resourceOutput

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	model "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/network"
)

type FirewallPort interface {
	Render(context.Context, *model.Firewall)
	RenderCreate(context.Context, *model.Firewall)
	RenderUpdate(context.Context, *model.Firewall)
	RenderDelete(context.Context, *model.Firewall)
}
