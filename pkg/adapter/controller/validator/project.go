package validator

import (
	"context"
	context2 "github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

func GetProject(ctx context.Context) errors.Error {
	if ctx.Value(context2.ProjectIDKey) == nil {
		panic(errors.BadRequest.WithMessage("ProjectID is mandatory"))
	}
	if _, ok := ctx.Value(context2.ProjectIDKey).(string); !ok {
		panic(errors.BadRequest.WithMessage("ProjectID must be a string"))
	}
	return errors.OK
}

func ListProjectByOwner(ctx context.Context) errors.Error {
	if ctx.Value(context2.OwnerIDKey) == nil {
		return errors.BadRequest.WithMessage("OwnerID is mandatory")
	}
	if _, ok := ctx.Value(context2.OwnerIDKey).(string); !ok {
		return errors.BadRequest.WithMessage("OwnerID must be a string")
	}
	return errors.OK
}

func CreateProject(ctx context.Context) errors.Error {
	if ctx.Value(context2.RequestKey) == nil {
		panic(errors.BadRequest.WithMessage("Request is mandatory"))
	}
	if _, ok := ctx.Value(context2.RequestKey).(dto.CreateProjectRequest); !ok {
		panic(errors.BadRequest.WithMessage("Request must be a Project DTO"))
	}
	return errors.OK
}

func UpdateProject(ctx context.Context) errors.Error {
	if ctx.Value(context2.RequestKey) == nil {
		return errors.BadRequest.WithMessage("Request is mandatory")
	}
	if _, ok := ctx.Value(context2.RequestKey).(dto.UpdateProjectRequest); !ok {
		return errors.BadRequest.WithMessage("Request must be a Project DTO")
	}
	if ctx.Value(context2.ProjectIDKey) == nil {
		return errors.BadRequest.WithMessage("ProjectID is mandatory")
	}
	if _, ok := ctx.Value(context2.ProjectIDKey).(string); !ok {
		return errors.BadRequest.WithMessage("ProjectID must be a string")
	}
	return errors.OK
}

func DeleteProject(ctx context.Context) errors.Error {
	if ctx.Value(context2.ProjectIDKey) == nil {
		panic(errors.BadRequest.WithMessage("ProjectID is mandatory"))
	}
	if _, ok := ctx.Value(context2.ProjectIDKey).(string); !ok {
		panic(errors.BadRequest.WithMessage("ProjectID must be a string"))
	}
	return errors.OK
}
