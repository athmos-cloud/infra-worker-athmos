package rabbitmq

import (
	"encoding/json"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/logger"
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
	}
	ctx.Set(context.ProjectIDKey, message.Data.ProjectID)
	ctx.Set(context.ProviderTypeKey, message.Data.ProviderType)
	ctx.Set(context.ResourceTypeKey, message.Data.ResourceType)
	ctx.Set(context.RequestKey, message.Data.Payload)

	switch message.Data.Verb {
	case CREATE:
		rq.ResourceController.CreateResource(ctx)
		rq.handleResponse(ctx, eventTypeCreateRequestSent)
	case UPDATE:
		rq.ResourceController.UpdateResource(ctx)
		rq.handleResponse(ctx, eventTypeUpdateRequestSent)
	case DELETE:
		rq.ResourceController.DeleteResource(ctx)
		rq.handleResponse(ctx, eventTypeDeleteRequestSent)
	default:
		return
	}
}

func (rq *RabbitMQ) handleResponse(ctx context.Context, eventType eventType) {
	code := ctx.Value(context.ResponseCodeKey).(int)
	if code%100 == 2 {
		msg := messageSend{
			ProjectID: ctx.Value(context.ProjectIDKey).(string),
			Type:      eventType,
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
	msg := messageSend{
		ProjectID: ctx.Value(context.ProjectIDKey).(string),
		Type:      Error,
		Code:      ctx.Value(context.ResponseCodeKey).(int),
		Payload:   ctx.Value(context.ResponseKey),
	}
	rq.MessageHandler(rq.Channel, rq.ReceiveQueue, msg)
}
