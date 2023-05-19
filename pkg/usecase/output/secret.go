package output

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/secret"
)

type SecretPort interface {
	Render(context.Context, *secret.Secret)
	RenderCreate(context.Context, *secret.Secret)
	RenderUpdate(context.Context, *secret.Secret)
	RenderAll(context.Context, *[]secret.Secret)
	RenderDelete(context.Context)
}
