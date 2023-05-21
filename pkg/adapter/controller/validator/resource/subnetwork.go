package resourceValidator

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

func CreateSubnetwork(ctx context.Context) errors.Error {
	if _, ok := ctx.Value(context.RequestKey).(dto.CreateSubnetworkRequest); !ok {
		return errors.BadRequest.WithMessage(fmt.Sprintf("expected %+v, got %+v", dto.CreateSubnetworkRequest{}, ctx.Value(context.RequestKey)))
	}
	return errors.OK
}

func GetSubnetwork(ctx context.Context) errors.Error {
	if _, ok := ctx.Value(context.RequestKey).(dto.GetSubnetworkRequest); !ok {
		return errors.BadRequest.WithMessage(fmt.Sprintf("expected %+v, got %+v", dto.GetSubnetworkRequest{}, ctx.Value(context.RequestKey)))
	}
	return errors.OK
}

func UpdateSubnetwork(ctx context.Context) errors.Error {
	if _, ok := ctx.Value(context.RequestKey).(dto.UpdateSubnetworkRequest); !ok {
		return errors.BadRequest.WithMessage(fmt.Sprintf("expected %+v, got %+v", dto.UpdateSubnetworkRequest{}, ctx.Value(context.RequestKey)))
	}
	return errors.OK
}

func DeleteSubnetwork(ctx context.Context) errors.Error {
	if _, ok := ctx.Value(context.RequestKey).(dto.DeleteSubnetworkRequest); !ok {
		return errors.BadRequest.WithMessage(fmt.Sprintf("expected %+v, got %+v", dto.DeleteSubnetworkRequest{}, ctx.Value(context.RequestKey)))
	}
	return errors.OK
}
