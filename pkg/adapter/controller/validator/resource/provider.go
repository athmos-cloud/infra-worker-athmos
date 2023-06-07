package resourceValidator

import (
	"encoding/json"
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

func GetProvider(ctx context.Context) errors.Error {
	if _, ok := ctx.Value(context.RequestKey).(map[string]interface{}); !ok {
		return errors.BadRequest.WithMessage(fmt.Sprintf("expected a map got %v", ctx.Value(context.RequestKey)))
	}
	req := ctx.Value(context.RequestKey).(map[string]interface{})
	jsonbody, errMarshall := json.Marshal(req)
	if errMarshall != nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("Invalid JSON : %v", req))
	}
	dtoRequest := dto.GetProviderRequest{}
	if errUnmarshall := json.Unmarshal(jsonbody, &dtoRequest); errUnmarshall != nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("Expected request %+v, got %v", dto.GetProviderRequest{}, req))
	}
	ctx.Set(context.RequestKey, dtoRequest)

	return errors.OK
}

func Stack(ctx context.Context) errors.Error {
	if _, ok := ctx.Value(context.RequestKey).(map[string]interface{}); !ok {
		return errors.BadRequest.WithMessage(fmt.Sprintf("expected a map got %v", ctx.Value(context.RequestKey)))
	}
	req := ctx.Value(context.RequestKey).(map[string]interface{})
	jsonbody, errMarshall := json.Marshal(req)
	if errMarshall != nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("Invalid JSON : %v", req))
	}
	dtoRequest := dto.GetProviderStackRequest{}
	if errUnmarshall := json.Unmarshal(jsonbody, &dtoRequest); errUnmarshall != nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("Expected request %+v, got %v", dto.GetProviderStackRequest{}, req))
	}
	ctx.Set(context.RequestKey, dtoRequest)

	return errors.OK
}

func CreateProvider(ctx context.Context) errors.Error {
	if _, ok := ctx.Value(context.RequestKey).(map[string]interface{}); !ok {
		return errors.BadRequest.WithMessage(fmt.Sprintf("expected a map got %v", ctx.Value(context.RequestKey)))
	}
	req := ctx.Value(context.RequestKey).(map[string]interface{})
	jsonbody, errMarshall := json.Marshal(req)
	if errMarshall != nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("Invalid JSON : %v", req))
	}
	dtoRequest := dto.CreateProviderRequest{}
	if errUnmarshall := json.Unmarshal(jsonbody, &dtoRequest); errUnmarshall != nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("Expected request %+v, got %v", dto.CreateProviderRequest{}, req))
	}
	ctx.Set(context.RequestKey, dtoRequest)

	return errors.OK
}

func UpdateProvider(ctx context.Context) errors.Error {
	if _, ok := ctx.Value(context.RequestKey).(map[string]interface{}); !ok {
		return errors.BadRequest.WithMessage(fmt.Sprintf("expected a map got %v", ctx.Value(context.RequestKey)))
	}
	req := ctx.Value(context.RequestKey).(map[string]interface{})
	jsonbody, errMarshall := json.Marshal(req)
	if errMarshall != nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("Invalid JSON : %v", req))
	}
	dtoRequest := dto.UpdateProviderRequest{}
	if errUnmarshall := json.Unmarshal(jsonbody, &dtoRequest); errUnmarshall != nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("Expected request %+v, got %v", dto.UpdateProviderRequest{}, req))
	}
	ctx.Set(context.RequestKey, dtoRequest)

	return errors.OK
}

func DeleteProvider(ctx context.Context) errors.Error {
	if _, ok := ctx.Value(context.RequestKey).(map[string]interface{}); !ok {
		return errors.BadRequest.WithMessage(fmt.Sprintf("expected a map got %v", ctx.Value(context.RequestKey)))
	}
	req := ctx.Value(context.RequestKey).(map[string]interface{})
	jsonbody, errMarshall := json.Marshal(req)
	if errMarshall != nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("Invalid JSON : %v", req))
	}
	dtoRequest := dto.DeleteProviderRequest{}
	if errUnmarshall := json.Unmarshal(jsonbody, &dtoRequest); errUnmarshall != nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("Expected request %+v, got %v", dto.DeleteProviderRequest{}, req))
	}
	ctx.Set(context.RequestKey, dtoRequest)
	return errors.OK
}
