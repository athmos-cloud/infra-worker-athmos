package resourceOutput

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	model "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/instance"
)

type SqlDBPort interface {
	Render(context.Context, *model.SqlDB)
	RenderCreate(context.Context, *model.SqlDB)
	RenderUpdate(context.Context, *model.SqlDB)
	RenderDelete(context.Context, *model.SqlDB)
}
