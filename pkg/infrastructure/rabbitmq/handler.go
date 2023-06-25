package rabbitmq

import (
	"encoding/json"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/logger"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

func (rq *RabbitMQ) handleMessage(ctx context.Context, msg amqp.Delivery, err error) {
	if err != nil {
		logger.Error.Fatalf("Error occurred in RMQ consumer: %v", err)
	}
	message := messageReceived{}
	err = json.Unmarshal(msg.Body, &message)
	if err != nil {
		logger.Error.Printf("Wrong message format: %s", err)
		return
	}
	errorType := func(err errors.Error) {
		ctx.Set(context.ResponseCodeKey, 400)
		ctx.Set(context.ResponseKey, gin.H{"error": err.ToString()})
		rq.handleResponse(ctx, metadata.StatusTypeError)
	}
	providerType, errProvider := types.ProviderFromString(message.Data.ProviderType)
	if !errProvider.IsOk() {
		errorType(errProvider)
		return
	}
	resourceType, errResource := types.ResourceFromString(message.Data.ResourceType)
	if !errResource.IsOk() {
		errorType(errResource)
		return
	}
	ctx.Set(context.ProjectIDKey, message.Data.ProjectID)
	ctx.Set(context.ProviderTypeKey, providerType)
	ctx.Set(context.ResourceTypeKey, resourceType)
	ctx.Set(context.RequestKey, message.Data.Payload)

	switch message.Data.Verb {
	case CREATE:
		rq.ResourceController.CreateResource(ctx)
		rq.handleResponse(ctx, metadata.StatusTypeCreateRequestSent)
	case UPDATE:
		rq.ResourceController.UpdateResource(ctx)
		rq.handleResponse(ctx, metadata.StatusTypeUpdateRequestSent)
	case DELETE:
		rq.ResourceController.DeleteResource(ctx)
		rq.handleResponse(ctx, metadata.StatusTypeDeleteRequestSent)
	default:
		return
	}
}

func (rq *RabbitMQ) handleResponse(ctx context.Context, statusType metadata.StatusType) {
	code := ctx.Value(context.ResponseCodeKey).(int)
	if code%100 == 2 {
		msg := MessageSend{
			ProjectID: ctx.Value(context.ProjectIDKey).(string),
			Type:      statusType,
			Code:      code,
			Payload:   ctx.Value(context.ResponseKey),
		}
		rq.MessageHandler(rq.Channel, rq.ReceiveQueue, msg)
	} else {
		rq.handleError(ctx)
	}
	clearContext(ctx)
}

func (rq *RabbitMQ) handleError(ctx context.Context) {
	msg := MessageSend{
		ProjectID: ctx.Value(context.ProjectIDKey).(string),
		Type:      metadata.StatusTypeError,
		Code:      ctx.Value(context.ResponseCodeKey).(int),
		Payload:   ctx.Value(context.ResponseKey),
	}
	rq.MessageHandler(rq.Channel, rq.ReceiveQueue, msg)
}
