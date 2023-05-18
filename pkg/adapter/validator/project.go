package validator

import (
	"context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/share"
)

func GetProject(ctx context.Context) {
	if ctx.Value(share.ProjectIDKey) == nil {
		panic(errors.BadRequest.WithMessage("ProjectID is mandatory"))
	}
	if _, ok := ctx.Value(share.ProjectIDKey).(string); !ok {
		panic(errors.BadRequest.WithMessage("ProjectID must be a string"))
	}
}

func ListProjectByOwner(ctx context.Context) {
	if ctx.Value(share.OwnerIDKey) == nil {
		panic(errors.BadRequest.WithMessage("OwnerID is mandatory"))
	}
	if _, ok := ctx.Value(share.OwnerIDKey).(string); !ok {
		panic(errors.BadRequest.WithMessage("OwnerID must be a string"))
	}
}

func CreateProject(ctx context.Context) {
	if ctx.Value(share.RequestContextKey) == nil {
		panic(errors.BadRequest.WithMessage("Request is mandatory"))
	}
	if _, ok := ctx.Value(share.RequestContextKey).(dto.CreateProjectRequest); !ok {
		panic(errors.BadRequest.WithMessage("Request must be a Project DTO"))
	}
}

func UpdateProject(ctx context.Context) {
	if ctx.Value(share.RequestContextKey) == nil {
		panic(errors.BadRequest.WithMessage("Request is mandatory"))
	}
	if _, ok := ctx.Value(share.RequestContextKey).(dto.UpdateProjectRequest); !ok {
		panic(errors.BadRequest.WithMessage("Request must be a Project DTO"))
	}
}

func DeleteProject(ctx context.Context) {
	if ctx.Value(share.ProjectIDKey) == nil {
		panic(errors.BadRequest.WithMessage("ProjectID is mandatory"))
	}
	if _, ok := ctx.Value(share.ProjectIDKey).(string); !ok {
		panic(errors.BadRequest.WithMessage("ProjectID must be a string"))
	}
}
