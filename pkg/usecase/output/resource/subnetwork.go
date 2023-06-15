package resourceOutput

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	model "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/network"
)

type SubnetworkPort interface {
	Render(context.Context, *model.Subnetwork)
	RenderCreate(context.Context, *model.Subnetwork)
	RenderUpdate(context.Context, *model.Subnetwork)
	RenderDelete(context.Context, *model.Subnetwork)
}
