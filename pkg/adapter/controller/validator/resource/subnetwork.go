package resourceValidator

import (
	"encoding/json"
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

func CreateSubnetwork(ctx context.Context) errors.Error {
	req := ctx.Value(context.RequestKey).(string)
	dtoRequest := dto.CreateSubnetworkRequest{}
	if errUnmarshall := json.Unmarshal([]byte(req), &dtoRequest); errUnmarshall != nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("Expected request %+v, got %v", dto.CreateSubnetworkRequest{}, req))
	}
	ctx.Set(context.RequestKey, dtoRequest)
	return errors.OK
}

func GetSubnetwork(ctx context.Context) errors.Error {
	return errors.OK
}

func UpdateSubnetwork(ctx context.Context) errors.Error {
	req := ctx.Value(context.RequestKey).(string)
	dtoRequest := dto.UpdateSubnetworkRequest{}
	if errUnmarshall := json.Unmarshal([]byte(req), &dtoRequest); errUnmarshall != nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("Expected request %+v, got %v", dto.UpdateSubnetworkRequest{}, req))
	}
	ctx.Set(context.RequestKey, dtoRequest)
	return errors.OK
}

func DeleteSubnetwork(ctx context.Context) errors.Error {
	req := ctx.Value(context.RequestKey).(string)
	dtoRequest := dto.DeleteSubnetworkRequest{}
	if errUnmarshall := json.Unmarshal([]byte(req), &dtoRequest); errUnmarshall != nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("Expected request %+v, got %v", dto.DeleteSubnetworkRequest{}, req))
	}
	ctx.Set(context.RequestKey, dtoRequest)

	return errors.OK
}
