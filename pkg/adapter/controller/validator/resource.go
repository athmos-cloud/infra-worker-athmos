package validator

import (
	"fmt"
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
	if _, ok := ctx.Value(context.ProviderTypeKey).(types.Provider); !ok {
		return errors.BadRequest.WithMessage(fmt.Sprintf("Invalid provider %v", ctx.Value(context.ProviderTypeKey)))
	}
	if ctx.Value(context.ResourceTypeKey) == nil {
		return errors.BadRequest.WithMessage("resource type is required")
	}
	if _, ok := ctx.Value(context.ResourceTypeKey).(types.Resource); !ok {
		return errors.BadRequest.WithMessage(fmt.Sprintf("Invalid resource type %v", ctx.Value(context.ResourceTypeKey)))
	}
	if ctx.Value(context.RequestKey) == nil {
		return errors.BadRequest.WithMessage("request is required")
	}

	return errors.OK
}
