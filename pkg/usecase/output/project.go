package output

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model"
)

type ProjectPort interface {
	Render(context.Context, *model.Project)
	RenderCreate(context.Context, *model.Project)
	RenderAll(context.Context, []*model.Project)
}
