package resourceOutput

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	model "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
)

type VMPort interface {
	Render(context.Context, *model.VM)
	RenderCreate(context.Context, *model.VM)
	RenderUpdate(context.Context, *model.VM)
	RenderDelete(context.Context, *model.VM)
}
