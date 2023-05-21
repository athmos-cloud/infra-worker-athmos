package rabbitmq

import (
	"encoding/json"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/config"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/logger"
	"github.com/streadway/amqp"
)

type Verb string

const (
	CREATE Verb = "create"
	UPDATE Verb = "update"
	DELETE Verb = "delete"
)

type Message struct {
	Verb    Verb        `json:"verb"`
	Payload interface{} `json:"payload"`
}

func Publish(payload Event) {
	messageBody, _ := json.Marshal(payload)
	message := amqp.Publishing{
		ContentType: "application/json",
		Body:        messageBody,
	}
	logger.Info.Printf("Publishing message : %s", messageBody)
	// Attempt to publish a message to the rabbitmq.
	if err := Queue.Channel.Publish(
		"",                            // exchange
		config.Current.Queue.Exchange, // rabbitmq name
		false,                         // mandatory
		false,                         // immediate
		message,                       // message to publish
	); err != nil {
		logger.Error.Printf("Error publishing a message to the rabbitmq: %s", err)
	}

}
