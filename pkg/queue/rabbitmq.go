package queue

import (
	"github.com/PaulBarrie/infra-worker/pkg/kernel/config"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/logger"
	"github.com/streadway/amqp"
	"log"
	"sync"
)

type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	MsgHandler func(queue string, msg amqp.Delivery, err error)
}

var Queue *RabbitMQ
var lock = &sync.Mutex{}

func init() {
	lock.Lock()
	defer lock.Unlock()

	if Queue == nil {
		logger.Info.Printf("Connecting to RabbitMQ on %s", config.Current.Queue.URI)
		conn, err := amqp.Dial(config.Current.Queue.URI)
		if err != nil {
			log.Fatalf("Failed to connect to RabbitMQ: %v", err)
		}

		ch, err := conn.Channel()
		if err != nil {
			logger.Error.Fatalf("Failed to open a channel: %v", err)
		}
		Queue = &RabbitMQ{
			Connection: conn,
			Channel:    ch,
			MsgHandler: HandleMessage,
		}
	}
}

func Close() {
	errCon := Queue.Connection.Close()
	if errCon != nil {
		logger.Error.Fatalf("Failed to close a connection: %v", errCon)
	}
	errCh := Queue.Channel.Close()
	if errCh != nil {
		logger.Error.Fatalf("Failed to close a channel: %v", errCh)
	}
}

func HandleMessage(queue string, msg amqp.Delivery, err error) {
	if err != nil {
		logger.Error.Fatalf("Error occurred in RMQ consumer", err)
	}
	logger.Info.Printf("Message received on '%s' queue: %s", queue, string(msg.Body))
}

func (queue *RabbitMQ) OnError(err error, msg string) {
	if err != nil {
		queue.MsgHandler(config.Current.Queue.Queue, amqp.Delivery{}, err)
	}
}

func (queue *RabbitMQ) StartConsumer() {
	queueName := config.Current.Queue.Queue
	q, err := queue.Channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	queue.OnError(err, "Failed to declare a queue")

	msgs, err := queue.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	queue.OnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			queue.MsgHandler(queueName, d, nil)
		}
	}()
	logger.Info.Printf("Started listening for messages on '%s' queue", queueName)
	<-forever
}
