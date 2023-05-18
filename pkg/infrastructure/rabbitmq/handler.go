package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/config"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/logger"
	"github.com/streadway/amqp"
)

func (queue *RabbitMQ) HandleMessage(ctx context.Context, msg amqp.Delivery, err error) {
	if err != nil {
		logger.Error.Fatalf("Error occurred in RMQ consumer: %v", err)
	}
	message := Message{}
	err = json.Unmarshal(msg.Body, &message)
	if err != nil {
		logger.Error.Printf("Wrong message format: %s", err)
	}
	logger.Info.Printf("Message received : %s", message.Verb)
	switch message.Verb {
	case CREATE:
		svcErr := errors.OK
		defer func() {
			if r := recover(); r != nil {
				logger.Info.Printf("Error occurred in RMQ consumer: %v", r)
				svcErr = r.(errors.Error)
				logger.Error.Printf("Error occurred in RMQ consumer: %v", svcErr)
			}
		}()
		resp := queue.ResourceService.CreateResource(ctx, mapToCreateResourceRequest(message.Payload.(map[string]interface{})))
		if !svcErr.IsOk() {
			Publish(Event{
				ProjectID: resp.ProjectID,
				Code:      svcErr.Code,
				Type:      CreateError,
				Payload:   fmt.Sprintf("Error occurred in RMQ consumer: %v", svcErr),
			})
			return
		}
		Publish(Event{
			ProjectID: resp.ProjectID,
			Code:      200,
			Type:      CreateRequestTreated,
			Payload:   resp.Resource,
		})
	case UPDATE:
		var payload resource.UpdateResourceRequest
		svcErr := errors.NoContent
		jsonData, errMarshal := json.Marshal(message.Payload)
		if errMarshal != nil {
			logger.Error.Printf("Can't marshall payload : %v", errMarshal)
		}
		errUnmarshall := json.Unmarshal(jsonData, &payload)
		if errUnmarshall != nil {
			return
		}
		defer func() {
			if r := recover(); r != nil {
				svcErr = r.(errors.Error)
			}
		}()
		queue.ResourceService.UpdateResource(ctx, payload)
		if !svcErr.IsOk() {
			Publish(Event{
				ProjectID: payload.ProjectID,
				Code:      svcErr.Code,
				Type:      CreateError,
				Payload:   fmt.Sprintf("Error occurred in RMQ consumer: %v", svcErr),
			})
			return
		}
	case DELETE:
		var payload resource.DeleteResourceRequest
		svcErr := errors.NoContent

		jsonData, errMarshal := json.Marshal(message.Payload)
		if errMarshal != nil {
			logger.Error.Printf("Can't marshall payload : %v", errMarshal)
		}
		errUnmarshall := json.Unmarshal(jsonData, &payload)
		if errUnmarshall != nil {
			return
		}
		defer func() {
			if r := recover(); r != nil {
				svcErr = r.(errors.Error)
			}
		}()
		queue.ResourceService.DeleteResource(ctx, payload)
		if !svcErr.IsOk() {
			logger.Error.Printf("Error occurred in RMQ consumer: %v", svcErr)
		}
	}
}

func (queue *RabbitMQ) OnError(err errors.Error) {
	if !err.IsOk() {
		queue.MessageHandler(config.Current.Queue.Queue, amqp.Delivery{}, err)
	}
}
