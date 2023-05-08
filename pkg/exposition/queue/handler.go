package queue

import (
	"context"
	"encoding/json"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/application/resource"
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
		var payload resource.CreateResourceRequest
		svcErr := errors.OK

		jsonData, errMarshal := json.Marshal(message.Payload)
		if errMarshal != nil {
			logger.Error.Printf("Can't marshall payload : %v", errMarshal)
		}
		errUnmarshall := json.Unmarshal(jsonData, &payload)
		if errUnmarshall != nil {
			logger.Error.Printf("Can't unmarshall payload : %v", errUnmarshall)
		}
		defer func() {
			if r := recover(); r != nil {
				logger.Info.Printf("Error occurred in RMQ consumer: %v", r)
				svcErr = r.(errors.Error)
				logger.Error.Printf("Error occurred in RMQ consumer: %v", svcErr)
			}
		}()
		_ = queue.ResourceService.CreateResource(ctx, payload)
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
			logger.Error.Printf("Error occurred in RMQ consumer: %v", svcErr)
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

func (queue *RabbitMQ) OnError(err error, msg string) {
	if err != nil {
		queue.MessageHandler(config.Current.Queue.Queue, amqp.Delivery{}, err)
	}
}
