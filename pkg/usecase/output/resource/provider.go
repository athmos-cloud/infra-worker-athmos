package resourceOutput

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	model "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
)

type ProviderPort interface {
	Render(context.Context, *model.Provider)
	RenderCreate(context.Context, *model.Provider)
	RenderUpdate(context.Context, *model.Provider)
	RenderAll(context.Context, *[]model.Provider)
}
