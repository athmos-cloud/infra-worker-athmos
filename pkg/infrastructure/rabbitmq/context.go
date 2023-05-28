package rabbitmq

import (
	goContext "context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
)

type rabbitContext struct {
	goContext.Context
}

func (rc *rabbitContext) JSON(code int, val any) {
	rc.Set(context.ResponseKey, val)
	rc.Set(context.ResponseCodeKey, code)
}

func clearContext(ctx context.Context) {
	ctx.Set(context.ResponseKey, nil)
	ctx.Set(context.ResponseCodeKey, nil)
	ctx.Set(context.RequestKey, nil)
	ctx.Set(context.ProjectIDKey, nil)
	ctx.Set(context.OwnerIDKey, nil)
	ctx.Set(context.ProviderTypeKey, nil)
	ctx.Set(context.ResourceTypeKey, nil)
}
func (rc *rabbitContext) Set(key string, val any) {
	rc.Context = goContext.WithValue(rc.Context, key, val)
}

func NewContext() context.Context {
	return &rabbitContext{goContext.Background()}
}
