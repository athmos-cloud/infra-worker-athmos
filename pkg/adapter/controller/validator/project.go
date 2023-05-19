package validator

import (
	"context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/share"
)

func GetProject(ctx context.Context) errors.Error {
	if ctx.Value(share.ProjectIDKey) == nil {
		panic(errors.BadRequest.WithMessage("ProjectID is mandatory"))
	}
	if _, ok := ctx.Value(share.ProjectIDKey).(string); !ok {
		panic(errors.BadRequest.WithMessage("ProjectID must be a string"))
	}
	return errors.OK
}

func ListProjectByOwner(ctx context.Context) errors.Error {
	if ctx.Value(share.OwnerIDKey) == nil {
		return errors.BadRequest.WithMessage("OwnerID is mandatory")
	}
	if _, ok := ctx.Value(share.OwnerIDKey).(string); !ok {
		return errors.BadRequest.WithMessage("OwnerID must be a string")
	}
	return errors.OK
}

func CreateProject(ctx context.Context) errors.Error {
	if ctx.Value(share.RequestContextKey) == nil {
		panic(errors.BadRequest.WithMessage("Request is mandatory"))
	}
	if _, ok := ctx.Value(share.RequestContextKey).(dto.CreateProjectRequest); !ok {
		panic(errors.BadRequest.WithMessage("Request must be a Project DTO"))
	}
	return errors.OK
}

func UpdateProject(ctx context.Context) errors.Error {
	if ctx.Value(share.RequestContextKey) == nil {
		return errors.BadRequest.WithMessage("Request is mandatory")
	}
	if _, ok := ctx.Value(share.RequestContextKey).(dto.UpdateProjectRequest); !ok {
		return errors.BadRequest.WithMessage("Request must be a Project DTO")
	}
	if ctx.Value(share.ProjectIDKey) == nil {
		return errors.BadRequest.WithMessage("ProjectID is mandatory")
	}
	if _, ok := ctx.Value(share.ProjectIDKey).(string); !ok {
		return errors.BadRequest.WithMessage("ProjectID must be a string")
	}
	return errors.OK
}

func DeleteProject(ctx context.Context) errors.Error {
	if ctx.Value(share.ProjectIDKey) == nil {
		panic(errors.BadRequest.WithMessage("ProjectID is mandatory"))
	}
	if _, ok := ctx.Value(share.ProjectIDKey).(string); !ok {
		panic(errors.BadRequest.WithMessage("ProjectID must be a string"))
	}
	return errors.OK
}
