package validator

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

func Resource(ctx context.Context) errors.Error {
	if ctx.Value(context.ProjectIDKey) == nil {
		return errors.BadRequest.WithMessage("project id is required")
	}
	if ctx.Value(context.ProjectIDKey).(string) == "" {
		return errors.BadRequest.WithMessage("project id is required")
	}
	if ctx.Value(context.ProviderTypeKey) == nil {
		return errors.BadRequest.WithMessage("provider type is required")
	}
	if _, err := types.StringToProvider(ctx.Value(context.ProviderTypeKey).(string)); !err.IsOk() {
		return err
	}
	if ctx.Value(context.ResourceTypeKey) == nil {
		return errors.BadRequest.WithMessage("resource type is required")
	}
	if _, err := types.StringToResource(ctx.Value(context.ResourceTypeKey).(string)); !err.IsOk() {
		return err
	}
	if ctx.Value(context.RequestKey) == nil {
		return errors.BadRequest.WithMessage("request is required")
	}

	return errors.OK
}
