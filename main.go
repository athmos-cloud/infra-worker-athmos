package main

import (
	"context"
	"github.com/PaulBarrie/infra-worker/pkg/http"
	"github.com/PaulBarrie/infra-worker/pkg/queue"
	"github.com/PaulBarrie/infra-worker/pkg/repository/mongo"
	"github.com/PaulBarrie/infra-worker/pkg/service/plugin"
	"github.com/PaulBarrie/infra-worker/pkg/service/project"
	"github.com/PaulBarrie/infra-worker/pkg/service/resource"
)

var (
	DefaultWorkdir   string = "/tmp/infra-worker"
	PluginRepository        = mongo.Client
)

func main() {
	ctx := context.Background()
	projectService := project.Service{
		ProjectRepository: mongo.Client,
	}
	pluginService := plugin.Service{}
	resourceService := resource.Service{
		ProjectRepository: mongo.Client,
		PluginRepository:  PluginRepository,
	}
	server := http.New(&projectService, &pluginService, &resourceService)
	queue.Queue.SetServices(&resourceService)
	go queue.Queue.StartConsumer(ctx)
	defer queue.Close()
	server.Start()
}
