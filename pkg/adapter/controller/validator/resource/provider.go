package resourceValidator

import (
	"encoding/json"
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

func GetProvider(ctx context.Context) errors.Error {
	return errors.OK
}

func CreateProvider(ctx context.Context) errors.Error {
	req := ctx.Value(context.RequestKey).(string)
	dtoRequest := dto.CreateProviderRequest{}
	if errUnmarshall := json.Unmarshal([]byte(req), &dtoRequest); errUnmarshall != nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("Expected request %+v, got %v", dto.CreateProviderRequest{}, req))
	}
	ctx.Set(context.RequestKey, dtoRequest)

	return errors.OK
}

func UpdateProvider(ctx context.Context) errors.Error {
	req := ctx.Value(context.RequestKey).(string)
	dtoRequest := dto.UpdateProviderRequest{}
	if errUnmarshall := json.Unmarshal([]byte(req), &dtoRequest); errUnmarshall != nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("Expected request %+v, got %v", dto.UpdateProviderRequest{}, req))
	}
	ctx.Set(context.RequestKey, dtoRequest)

	return errors.OK
}

func DeleteProvider(ctx context.Context) errors.Error {
	req := ctx.Value(context.RequestKey).(string)
	dtoRequest := dto.DeleteProviderRequest{}
	if errUnmarshall := json.Unmarshal([]byte(req), &dtoRequest); errUnmarshall != nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("Expected request %+v, got %v", dto.DeleteProviderRequest{}, req))
	}
	ctx.Set(context.RequestKey, dtoRequest)

	return errors.OK
}
