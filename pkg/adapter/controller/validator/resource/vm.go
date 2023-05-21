package resourceValidator

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

func GetVM(ctx context.Context) errors.Error {
	if _, ok := ctx.Value(context.RequestKey).(dto.GetVMRequest); !ok {
		return errors.BadRequest.WithMessage(fmt.Sprintf("expected %+v, got %+v", dto.GetVMRequest{}, ctx.Value(context.RequestKey)))
	}
	return errors.OK
}

func CreateVM(ctx context.Context) errors.Error {
	if _, ok := ctx.Value(context.RequestKey).(dto.CreateVMRequest); !ok {
		return errors.BadRequest.WithMessage(fmt.Sprintf("expected %+v, got %+v", dto.CreateVMRequest{}, ctx.Value(context.RequestKey)))
	}
	return errors.OK
}

func UpdateVM(ctx context.Context) errors.Error {
	if _, ok := ctx.Value(context.RequestKey).(dto.UpdateVMRequest); !ok {
		return errors.BadRequest.WithMessage(fmt.Sprintf("expected %+v, got %+v", dto.UpdateVMRequest{}, ctx.Value(context.RequestKey)))
	}
	return errors.OK
}

func DeleteVM(ctx context.Context) errors.Error {
	if _, ok := ctx.Value(context.RequestKey).(dto.DeleteVMRequest); !ok {
		return errors.BadRequest.WithMessage(fmt.Sprintf("expected %+v, got %+v", dto.DeleteVMRequest{}, ctx.Value(context.RequestKey)))
	}
	return errors.OK
}
