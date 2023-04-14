package main

import (
	"context"
	"github.com/PaulBarrie/infra-worker/pkg/exposition/http"
	"github.com/PaulBarrie/infra-worker/pkg/exposition/queue"
	"github.com/PaulBarrie/infra-worker/pkg/repository/mongo"
	"github.com/PaulBarrie/infra-worker/pkg/service"
)

var (
	DefaultWorkdir   string = "/tmp/infra-worker"
	PluginRepository        = mongo.Client
)

func main() {
	ctx := context.Background()
	projectService := service.ProjectService{
		ProjectRepository: mongo.Client,
	}
	pluginService := service.PluginService{}
	resourceService := service.ResourceService{
		ProjectRepository: mongo.Client,
		PluginRepository:  PluginRepository,
	}
	server := http.New(&projectService, &pluginService, &resourceService)
	queue.Queue.SetServices(&resourceService)
	go queue.Queue.StartConsumer(ctx)
	defer queue.Close()
	server.Start()
}
