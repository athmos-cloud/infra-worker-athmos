package validator

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/share"
)

func GetSecret(ctx context.Context) errors.Error {
	if ctx.Value(share.RequestContextKey) == nil {
		return errors.BadRequest.WithMessage("Request is mandatory")
	}
	if _, ok := ctx.Value(share.RequestContextKey).(dto.GetSecretRequest); !ok {
		return errors.BadRequest.WithMessage(fmt.Sprintf("Request must be a Secret DTO : %+v, got %+v", dto.GetSecretRequest{}, ctx.Value(share.RequestContextKey)))
	}
	return errors.OK
}

func ListProjectSecret(ctx context.Context) errors.Error {
	if ctx.Value(share.ProjectIDKey) == nil {
		return errors.BadRequest.WithMessage("ProjectID is mandatory")
	}
	if _, ok := ctx.Value(share.ProjectIDKey).(string); !ok {
		return errors.BadRequest.WithMessage("ProjectID must be a string")
	}
	return errors.OK
}

func CreateSecret(ctx context.Context) errors.Error {
	if ctx.Value(share.RequestContextKey) == nil {
		panic(errors.BadRequest.WithMessage("Request is mandatory"))
	}
	if _, ok := ctx.Value(share.RequestContextKey).(dto.CreateSecretRequest); !ok {
		return errors.BadRequest.WithMessage(fmt.Sprintf("Request must be a Secret DTO : %+v, got %+v", dto.CreateProjectRequest{}, ctx.Value(share.RequestContextKey)))
	}
	if ctx.Value(share.ProjectIDKey) == nil {
		return errors.BadRequest.WithMessage("ProjectID is mandatory")
	}
	if _, ok := ctx.Value(share.ProjectIDKey).(string); !ok {
		return errors.BadRequest.WithMessage("ProjectID must be a string")
	}
	return errors.OK
}

func UpdateSecret(ctx context.Context) errors.Error {
	if ctx.Value(share.RequestContextKey) == nil {
		return errors.BadRequest.WithMessage("Request is mandatory")
	}
	if _, ok := ctx.Value(share.RequestContextKey).(dto.UpdateSecretRequest); !ok {
		return errors.BadRequest.WithMessage(fmt.Sprintf("Request must be a Secret DTO : %+v, got %+v", dto.UpdateProjectRequest{}, ctx.Value(share.RequestContextKey)))
	}
	if ctx.Value(share.ProjectIDKey) == nil {
		return errors.BadRequest.WithMessage("ProjectID is mandatory")
	}
	if _, ok := ctx.Value(share.ProjectIDKey).(string); !ok {
		return errors.BadRequest.WithMessage("ProjectID must be a string")
	}
	return errors.OK
}

func DeleteSecret(ctx context.Context) errors.Error {
	if ctx.Value(share.RequestContextKey) == nil {
		return errors.BadRequest.WithMessage("Request is mandatory")
	}
	if _, ok := ctx.Value(share.RequestContextKey).(dto.DeleteSecretRequest); !ok {
		return errors.BadRequest.WithMessage(fmt.Sprintf("Request must be a Secret DTO : %+v, got %+v", dto.DeleteSecretRequest{}, ctx.Value(share.RequestContextKey)))
	}
	if ctx.Value(share.ProjectIDKey) == nil {
		return errors.BadRequest.WithMessage("ProjectID is mandatory")
	}
	if _, ok := ctx.Value(share.ProjectIDKey).(string); !ok {
		return errors.BadRequest.WithMessage("ProjectID must be a string")
	}
	return errors.OK
}
