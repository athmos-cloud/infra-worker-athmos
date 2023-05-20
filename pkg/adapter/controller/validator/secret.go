package validator

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

func GetSecret(ctx context.Context) errors.Error {
	if ctx.Value(context.RequestKey) == nil {
		return errors.BadRequest.WithMessage("Request is mandatory")
	}
	if _, ok := ctx.Value(context.RequestKey).(dto.GetSecretRequest); !ok {
		return errors.BadRequest.WithMessage(fmt.Sprintf("Request must be a Secret DTO : %+v, got %+v", dto.GetSecretRequest{}, ctx.Value(context.RequestKey)))
	}
	return errors.OK
}

func ListProjectSecret(ctx context.Context) errors.Error {
	if ctx.Value(context.ProjectIDKey) == nil {
		return errors.BadRequest.WithMessage("ProjectID is mandatory")
	}
	if _, ok := ctx.Value(context.ProjectIDKey).(string); !ok {
		return errors.BadRequest.WithMessage("ProjectID must be a string")
	}
	return errors.OK
}

func CreateSecret(ctx context.Context) errors.Error {
	if ctx.Value(context.RequestKey) == nil {
		panic(errors.BadRequest.WithMessage("Request is mandatory"))
	}
	if _, ok := ctx.Value(context.RequestKey).(dto.CreateSecretRequest); !ok {
		return errors.BadRequest.WithMessage(fmt.Sprintf("Request must be a Secret DTO : %+v, got %+v", dto.CreateProjectRequest{}, ctx.Value(context.RequestKey)))
	}
	if ctx.Value(context.ProjectIDKey) == nil {
		return errors.BadRequest.WithMessage("ProjectID is mandatory")
	}
	if _, ok := ctx.Value(context.ProjectIDKey).(string); !ok {
		return errors.BadRequest.WithMessage("ProjectID must be a string")
	}
	return errors.OK
}

func UpdateSecret(ctx context.Context) errors.Error {
	if ctx.Value(context.RequestKey) == nil {
		return errors.BadRequest.WithMessage("Request is mandatory")
	}
	if _, ok := ctx.Value(context.RequestKey).(dto.UpdateSecretRequest); !ok {
		return errors.BadRequest.WithMessage(fmt.Sprintf("Request must be a Secret DTO : %+v, got %+v", dto.UpdateProjectRequest{}, ctx.Value(context.RequestKey)))
	}
	if ctx.Value(context.ProjectIDKey) == nil {
		return errors.BadRequest.WithMessage("ProjectID is mandatory")
	}
	if _, ok := ctx.Value(context.ProjectIDKey).(string); !ok {
		return errors.BadRequest.WithMessage("ProjectID must be a string")
	}
	return errors.OK
}

func DeleteSecret(ctx context.Context) errors.Error {
	if ctx.Value(context.RequestKey) == nil {
		return errors.BadRequest.WithMessage("Request is mandatory")
	}
	if _, ok := ctx.Value(context.RequestKey).(dto.DeleteSecretRequest); !ok {
		return errors.BadRequest.WithMessage(fmt.Sprintf("Request must be a Secret DTO : %+v, got %+v", dto.DeleteSecretRequest{}, ctx.Value(context.RequestKey)))
	}
	if ctx.Value(context.ProjectIDKey) == nil {
		return errors.BadRequest.WithMessage("ProjectID is mandatory")
	}
	if _, ok := ctx.Value(context.ProjectIDKey).(string); !ok {
		return errors.BadRequest.WithMessage("ProjectID must be a string")
	}
	return errors.OK
}
