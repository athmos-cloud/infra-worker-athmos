package queue

import (
	"context"
	"encoding/json"
	dto "github.com/PaulBarrie/infra-worker/pkg/common/dto/resource"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/config"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/logger"
	"github.com/streadway/amqp"
)

func (queue *RabbitMQ) HandleMessage(ctx context.Context, msg amqp.Delivery, err error) {
	if err != nil {
		logger.Error.Fatalf("Error occurred in RMQ consumer", err)
	}
	logger.Info.Printf("Message received : %s", string(msg.Body))
	message := Message{}
	err = json.Unmarshal(msg.Body, &message)
	if err != nil {
		logger.Error.Printf("Wrong message format: %s", err)
	}

	switch message.Verb {
	case CREATE:
		var payload dto.CreateResourceRequest

		jsonData, errMarshal := json.Marshal(message.Payload)
		if errMarshal != nil {
			logger.Error.Printf("Can't marshall payload : %v", errMarshal)
		}
		errUnmarshall := json.Unmarshal(jsonData, &payload)
		if errUnmarshall != nil {
			return
		}
		resp, svcErr := queue.ResourceService.CreateResource(ctx, payload)
		if !svcErr.IsOk() {
			logger.Error.Printf("Error occurred in RMQ consumer", svcErr)
		}
		logger.Info.Printf("Message response : %s", resp)
	case UPDATE:
		var payload dto.UpdateResourceRequest

		jsonData, errMarshal := json.Marshal(message.Payload)
		if errMarshal != nil {
			logger.Error.Printf("Can't marshall payload : %v", errMarshal)
		}
		errUnmarshall := json.Unmarshal(jsonData, &payload)
		if errUnmarshall != nil {
			return
		}
		svcErr := queue.ResourceService.UpdateResource(ctx, payload)
		if !svcErr.IsOk() {
			logger.Error.Printf("Error occurred in RMQ consumer", svcErr)
		}
	case DELETE:
		var payload dto.DeleteResourceRequest

		jsonData, errMarshal := json.Marshal(message.Payload)
		if errMarshal != nil {
			logger.Error.Printf("Can't marshall payload : %v", errMarshal)
		}
		errUnmarshall := json.Unmarshal(jsonData, &payload)
		if errUnmarshall != nil {
			return
		}
		svcErr := queue.ResourceService.DeleteResource(ctx, payload)
		if !svcErr.IsOk() {
			logger.Error.Printf("Error occurred in RMQ consumer", svcErr)
		}
	}
}

func (queue *RabbitMQ) OnError(err error, msg string) {
	if err != nil {
		queue.MsgHandler(config.Current.Queue.Queue, amqp.Delivery{}, err)
	}
}
