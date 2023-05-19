package output

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model"
)

type ProjectPort interface {
	Render(context.Context, *model.Project)
	RenderCreate(context.Context, *model.Project)
	RenderUpdate(context.Context, *model.Project)
	RenderAll(context.Context, *[]model.Project)
}
