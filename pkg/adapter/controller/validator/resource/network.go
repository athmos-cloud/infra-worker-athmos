package resourceValidator

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

func GetNetwork(ctx context.Context) errors.Error {
	if _, ok := ctx.Value(context.RequestKey).(dto.GetNetworkRequest); !ok {
		return errors.BadRequest.WithMessage(fmt.Sprintf("expected %+v, got %+v", dto.GetNetworkRequest{}, ctx.Value(context.RequestKey)))
	}
	return errors.OK
}

func CreateNetwork(ctx context.Context) errors.Error {
	if _, ok := ctx.Value(context.RequestKey).(dto.CreateNetworkRequest); !ok {
		return errors.BadRequest.WithMessage(fmt.Sprintf("expected %+v, got %+v", dto.CreateNetworkRequest{}, ctx.Value(context.RequestKey)))
	}
	return errors.OK
}

func UpdateNetwork(ctx context.Context) errors.Error {
	if _, ok := ctx.Value(context.RequestKey).(dto.UpdateNetworkRequest); !ok {
		return errors.BadRequest.WithMessage(fmt.Sprintf("expected %+v, got %+v", dto.UpdateNetworkRequest{}, ctx.Value(context.RequestKey)))
	}
	return errors.OK
}

func DeleteNetwork(ctx context.Context) errors.Error {
	if _, ok := ctx.Value(context.RequestKey).(dto.DeleteNetworkRequest); !ok {
		return errors.BadRequest.WithMessage(fmt.Sprintf("expected %+v, got %+v", dto.DeleteNetworkRequest{}, ctx.Value(context.RequestKey)))
	}
	return errors.OK
}
