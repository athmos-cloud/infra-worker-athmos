package resourceValidator

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

func GetProvider(ctx context.Context) errors.Error {
	if _, ok := ctx.Value(context.RequestKey).(dto.GetProviderRequest); !ok {
		return errors.BadRequest.WithMessage(fmt.Sprintf("expected %+v, got %+v", dto.GetProviderRequest{}, ctx.Value(context.RequestKey)))
	}
	return errors.OK
}

func ListProviders(ctx context.Context) errors.Error {
	if _, ok := ctx.Value(context.RequestKey).(dto.ListProvidersRequest); !ok {
		return errors.BadRequest.WithMessage(fmt.Sprintf("expected %+v, got %+v", dto.ListProvidersRequest{}, ctx.Value(context.RequestKey)))
	}
	return errors.OK
}

func CreateProvider(ctx context.Context) errors.Error {
	if _, ok := ctx.Value(context.RequestKey).(dto.CreateProviderRequest); !ok {
		return errors.BadRequest.WithMessage(fmt.Sprintf("expected %+v, got %+v", dto.CreateProviderRequest{}, ctx.Value(context.RequestKey)))
	}
	return errors.OK
}

func UpdateProvider(ctx context.Context) errors.Error {
	if _, ok := ctx.Value(context.RequestKey).(dto.UpdateProviderRequest); !ok {
		return errors.BadRequest.WithMessage(fmt.Sprintf("expected %+v, got %+v", dto.UpdateProviderRequest{}, ctx.Value(context.RequestKey)))
	}
	return errors.OK
}

func DeleteProvider(ctx context.Context) errors.Error {
	if _, ok := ctx.Value(context.RequestKey).(dto.DeleteProviderRequest); !ok {
		return errors.BadRequest.WithMessage(fmt.Sprintf("expected %+v, got %+v", dto.DeleteProviderRequest{}, ctx.Value(context.RequestKey)))
	}
	return errors.OK
}
