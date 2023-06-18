package rabbitmq

import (
	"encoding/json"
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/config"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/logger"
	"github.com/streadway/amqp"
	"log"
)

type RabbitMQ struct {
	Connection         *amqp.Connection
	Channel            *amqp.Channel
	ReceiveQueue       string
	ResourceController controller.Resource
	MessageHandler     func(*amqp.Channel, string, MessageSend)
}

func New(uri string, receiveQueue string, sendQueue string, resourceController controller.Resource) *RabbitMQ {
	conn, err := amqp.Dial(uri)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		logger.Error.Fatalf("Failed to open a channel: %v", err)
	}
	return &RabbitMQ{
		Connection:         conn,
		Channel:            ch,
		ResourceController: resourceController,
		ReceiveQueue:       receiveQueue,
		MessageHandler: func(channel *amqp.Channel, queue string, msg MessageSend) {
			publish(channel, sendQueue, msg.WithNestWrapper())
		},
	}
}
func (rq *RabbitMQ) Close() {
	errCon := rq.Connection.Close()
	if errCon != nil {
		logger.Error.Fatalf("Failed to close a connection: %v", errCon)
	}
	errCh := rq.Channel.Close()
	if errCh != nil {
		logger.Error.Fatalf("Failed to close a channel: %v", errCh)
	}
}

func (rq *RabbitMQ) StartConsumer(ctx context.Context) {
	queueName := config.Current.Queue.IncomingQueue
	q, err := rq.Channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		panic(errors.InternalError.WithMessage("Failed to declare a rabbitmq"))
	}

	msgs, err := rq.Channel.Consume(
		q.Name, // rabbitmq
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		panic(errors.InternalError.WithMessage(fmt.Sprintf("Failed to register a consumer: %v", err)))
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			rq.handleMessage(ctx, d, nil)
		}
	}()
	logger.Info.Printf("Started listening for messages on '%s' rabbitmq", queueName)
	<-forever
}

func publish(channel *amqp.Channel, queue string, payload any) {
	messageBody, _ := json.Marshal(payload)
	message := amqp.Publishing{
		ContentType: "application/json",
		Body:        messageBody,
	}
	// Attempt to publish a message to the rabbitmq.
	logger.Info.Printf("Message to %s", queue)
	if err := channel.Publish(
		"",      // exchange
		queue,   // rabbitmq name
		false,   // mandatory
		false,   // immediate
		message, // message to publish
	); err != nil {
		logger.Error.Printf("Error publishing a message to the rabbitmq: %s", err)
	}
}
