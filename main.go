package main

import (
	"fmt"
	"github.com/PaulBarrie/infra-worker/pkg/application/service"
	"github.com/PaulBarrie/infra-worker/pkg/auth"
	"github.com/PaulBarrie/infra-worker/pkg/infrastructure/repository/mongo"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/logger"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/types"
	plugin2 "github.com/PaulBarrie/infra-worker/pkg/plugin/common"
	"github.com/PaulBarrie/infra-worker/pkg/plugin/common/config"
	"os"
)

var kafkaServer, kafkaTopic string

const (
	groupID = "test-group"
)

var (
	DefaultWorkdir   string = "/tmp/infra-worker"
	PluginRepository        = mongo.Client
)

func init() {
	kafkaServer = os.Getenv("KAFKA_BROKER_ADDRESS")
	kafkaTopic = readFromENV("KAFKA_TOPIC", "test")
	fmt.Println("Kafka Broker - ", kafkaServer)
	fmt.Println("Kafka topic - ", kafkaTopic)
}

func main() {

		_, err = service.NewContext(ctx, "1", user, types.Create, nil, DefaultWorkdir)
		if err != nil {
			logger.Error.Printf("Error creating plugin context: %s", err)
			os.Exit(1)
		}

		_, err = svc.LoadPlugin(ctx, "1")
		if err != nil {
			return
		}
	}

	//config := kafka.ConfigMap{"bootstrap.servers": kafkaServer, "group.id": groupID, "go.events.channel.enable": true}
	//consumer, consumerCreateErr := kafka.NewConsumer(&config)
	//if consumerCreateErr != nil {
	//	fmt.Println("consumer not created ", consumerCreateErr.Error())
	//	os.Exit(1)
	//}
	//subscriptionErr := consumer.Subscribe(kafkaTopic, nil)
	//if subscriptionErr != nil {
	//	fmt.Println("Unable to subscribe to topic " + kafkaTopic + " due to error - " + subscriptionErr.Error())
	//	os.Exit(1)
	//} else {
	//	fmt.Println("subscribed to topic ", kafkaTopic)
	//}
	//
	//for {
	//	fmt.Println("waiting for event...")
	//	kafkaEvent := <-consumer.Events()
	//	if kafkaEvent != nil {
	//		switch event := kafkaEvent.(type) {
	//		case *kafka.WithMessage:
	//			fmt.Println("WithMessage " + string(event.Value))
	//		case kafka.Error:
	//			fmt.Println("Consumer error ", event.String())
	//		case kafka.PartitionEOF:
	//			fmt.Println(kafkaEvent)
	//		default:
	//			fmt.Println(kafkaEvent)
	//		}
	//	} else {
	//		fmt.Println("Event was null")
	//	}
	//}

}

func readFromENV(key, defaultVal string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultVal
	}
	return value
}
