package resourceValidator

import (
	"encoding/json"
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

func GetFirewall(ctx context.Context) errors.Error {
	return errors.OK
}

func CreateFirewall(ctx context.Context) errors.Error {
	req := ctx.Value(context.RequestKey).(string)
	dtoRequest := dto.CreateFirewallRequest{}
	if errUnmarshall := json.Unmarshal([]byte(req), &dtoRequest); errUnmarshall != nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("Expected request %+v, got %v", dto.CreateFirewallRequest{}, req))
	}
	ctx.Set(context.RequestKey, dtoRequest)

	return errors.OK
}

func UpdateFirewall(ctx context.Context) errors.Error {
	req := ctx.Value(context.RequestKey).(string)
	dtoRequest := dto.UpdateFirewallRequest{}
	if errUnmarshall := json.Unmarshal([]byte(req), &dtoRequest); errUnmarshall != nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("Expected request %+v, got %v", dto.UpdateFirewallRequest{}, req))
	}
	ctx.Set(context.RequestKey, dtoRequest)

	return errors.OK
}

func DeleteFirewall(ctx context.Context) errors.Error {
	req := ctx.Value(context.RequestKey).(string)
	dtoRequest := dto.DeleteFirewallRequest{}
	if errUnmarshall := json.Unmarshal([]byte(req), &dtoRequest); errUnmarshall != nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("Expected request %+v, got %v", dto.DeleteFirewallRequest{}, req))
	}
	ctx.Set(context.RequestKey, dtoRequest)
	return errors.OK
}
