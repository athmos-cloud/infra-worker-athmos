package resourceValidator

import (
	"encoding/json"
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

func GetSqlDB(ctx context.Context) errors.Error {
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

func CreateSqlDB(ctx context.Context) errors.Error {
	if _, ok := ctx.Value(context.RequestKey).(map[string]interface{}); !ok {
		return errors.BadRequest.WithMessage(fmt.Sprintf("expected a map got %v", ctx.Value(context.RequestKey)))
	}
	req := ctx.Value(context.RequestKey).(map[string]interface{})
	jsonbody, errMarshall := json.Marshal(req)
	if errMarshall != nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("Invalid JSON : %v", req))
	}
	dtoRequest := dto.CreateSqlDBRequest{}
	if errUnmarshall := json.Unmarshal(jsonbody, &dtoRequest); errUnmarshall != nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("Expected request %+v, got %v", dto.CreateSqlDBRequest{}, req))
	}
	ctx.Set(context.RequestKey, dtoRequest)

	return errors.OK
}

func UpdateSqlDB(ctx context.Context) errors.Error {
	if _, ok := ctx.Value(context.RequestKey).(map[string]interface{}); !ok {
		return errors.BadRequest.WithMessage(fmt.Sprintf("expected a map got %v", ctx.Value(context.RequestKey)))
	}
	req := ctx.Value(context.RequestKey).(map[string]interface{})
	jsonbody, errMarshall := json.Marshal(req)
	if errMarshall != nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("Invalid JSON : %v", req))
	}
	dtoRequest := dto.UpdateSqlDBRequest{}
	if errUnmarshall := json.Unmarshal(jsonbody, &dtoRequest); errUnmarshall != nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("Expected request %+v, got %v", dto.UpdateSqlDBRequest{}, req))
	}
	ctx.Set(context.RequestKey, dtoRequest)

	return errors.OK
}

func DeleteSqlDB(ctx context.Context) errors.Error {
	if _, ok := ctx.Value(context.RequestKey).(map[string]interface{}); !ok {
		return errors.BadRequest.WithMessage(fmt.Sprintf("expected a map got %v", ctx.Value(context.RequestKey)))
	}
	req := ctx.Value(context.RequestKey).(map[string]interface{})
	jsonbody, errMarshall := json.Marshal(req)
	if errMarshall != nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("Invalid JSON : %v", req))
	}
	dtoRequest := dto.DeleteSqlDBRequest{}
	if errUnmarshall := json.Unmarshal(jsonbody, &dtoRequest); errUnmarshall != nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("Expected request %+v, got %v", dto.DeleteSqlDBRequest{}, req))
	}
	ctx.Set(context.RequestKey, dtoRequest)

	return errors.OK
}
