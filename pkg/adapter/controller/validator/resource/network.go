package resourceValidator

import (
	"encoding/json"
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

func GetNetwork(ctx context.Context) errors.Error {
	if _, ok := ctx.Value(context.RequestKey).(map[string]interface{}); !ok {
		return errors.BadRequest.WithMessage(fmt.Sprintf("expected a map got %v", ctx.Value(context.RequestKey)))
	}
	req := ctx.Value(context.RequestKey).(map[string]interface{})
	jsonbody, errMarshall := json.Marshal(req)
	if errMarshall != nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("Invalid JSON : %v", req))
	}
	dtoRequest := dto.GetResourceRequest{}
	if errUnmarshall := json.Unmarshal(jsonbody, &dtoRequest); errUnmarshall != nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("Expected request %+v, got %v", dto.GetResourceRequest{}, req))
	}
	ctx.Set(context.RequestKey, dtoRequest)
	return errors.OK
}

func CreateNetwork(ctx context.Context) errors.Error {
	req := ctx.Value(context.RequestKey).(string)
	dtoRequest := dto.CreateNetworkRequest{}
	if errUnmarshall := json.Unmarshal([]byte(req), &dtoRequest); errUnmarshall != nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("Expected request %+v, got %v", dto.CreateNetworkRequest{}, req))
	}
	ctx.Set(context.RequestKey, dtoRequest)
	return errors.OK
}

func UpdateNetwork(ctx context.Context) errors.Error {
	req := ctx.Value(context.RequestKey).(string)
	dtoRequest := dto.UpdateNetworkRequest{}
	if errUnmarshall := json.Unmarshal([]byte(req), &dtoRequest); errUnmarshall != nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("Expected request %+v, got %v", dto.UpdateNetworkRequest{}, req))
	}
	ctx.Set(context.RequestKey, dtoRequest)
	return errors.OK
}

func DeleteNetwork(ctx context.Context) errors.Error {
	req := ctx.Value(context.RequestKey).(string)
	dtoRequest := dto.DeleteNetworkRequest{}
	if errUnmarshall := json.Unmarshal([]byte(req), &dtoRequest); errUnmarshall != nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("Expected request %+v, got %v", dto.DeleteNetworkRequest{}, req))
	}
	ctx.Set(context.RequestKey, dtoRequest)
	return errors.OK
}
