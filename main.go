package main

import (
	"github.com/PaulBarrie/infra-worker/pkg/repository/mongo"
)

var (
	DefaultWorkdir   string = "/tmp/infra-worker"
	PluginRepository        = mongo.Client
)

func main() {
	//_ := context.Background()

	return

}
