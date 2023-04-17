package main

import (
	"context"
	"github.com/PaulBarrie/infra-worker/pkg/application"
	"github.com/PaulBarrie/infra-worker/pkg/dao/mongo"
	"github.com/PaulBarrie/infra-worker/pkg/exposition/http"
	"github.com/PaulBarrie/infra-worker/pkg/exposition/queue"
)

var (
	DefaultWorkdir   string = "/tmp/infra-worker"
	PluginRepository        = mongo.Client
)

func main() {
	ctx := context.Background()
	projectService := application.ProjectService{
		ProjectRepository: mongo.Client,
	}
	pluginService := application.PluginService{}
	resourceService := application.ResourceService{
		MongoDAO:         mongo.Client,
		PluginRepository: PluginRepository,
	}
	server := http.New(&projectService, &pluginService, &resourceService)
	queue.Queue.SetServices(&resourceService)
	go queue.Queue.StartConsumer(ctx)
	defer queue.Close()
	server.Start()
}
