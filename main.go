package main

import (
	"github.com/PaulBarrie/infra-worker/pkg/infrastructure/repository/mongo"
)

var kafkaServer, kafkaTopic string

const (
	groupID = "test-group"
)

var (
	DefaultWorkdir   string = "/tmp/infra-worker"
	PluginRepository        = mongo.Client
)

func main() {
	return
}
