package resourceValidator

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

func GetFirewall(ctx context.Context) errors.Error {
	if _, ok := ctx.Value(context.RequestKey).(dto.GetFirewallRequest); !ok {
		return errors.BadRequest.WithMessage(fmt.Sprintf("expected %+v, got %+v", dto.GetFirewallRequest{}, ctx.Value(context.RequestKey)))
	}
	return errors.OK
}

func CreateFirewall(ctx context.Context) errors.Error {
	if _, ok := ctx.Value(context.RequestKey).(dto.CreateFirewallRequest); !ok {
		return errors.BadRequest.WithMessage(fmt.Sprintf("expected %+v, got %+v", dto.CreateFirewallRequest{}, ctx.Value(context.RequestKey)))
	}
	return errors.OK
}

func UpdateFirewall(ctx context.Context) errors.Error {
	if _, ok := ctx.Value(context.RequestKey).(dto.UpdateFirewallRequest); !ok {
		return errors.BadRequest.WithMessage(fmt.Sprintf("expected %+v, got %+v", dto.UpdateFirewallRequest{}, ctx.Value(context.RequestKey)))
	}
	return errors.OK
}

func DeleteFirewall(ctx context.Context) errors.Error {
	if _, ok := ctx.Value(context.RequestKey).(dto.DeleteFirewallRequest); !ok {
		return errors.BadRequest.WithMessage(fmt.Sprintf("expected %+v, got %+v", dto.DeleteFirewallRequest{}, ctx.Value(context.RequestKey)))
	}
	return errors.OK
}
