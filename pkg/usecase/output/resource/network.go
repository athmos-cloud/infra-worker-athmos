package resourceOutput

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	model "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/network"
)

type NetworkPort interface {
	Render(context.Context, *model.Network)
	RenderCreate(context.Context, *model.Network)
	RenderUpdate(context.Context, *model.Network)
	RenderDelete(context.Context, *model.Network)
}
