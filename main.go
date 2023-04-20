package main

import (
	"context"
	"github.com/PaulBarrie/infra-worker/pkg/application"
	"github.com/PaulBarrie/infra-worker/pkg/dao/mongo"
	"github.com/PaulBarrie/infra-worker/pkg/exposition/http"
	"github.com/PaulBarrie/infra-worker/pkg/exposition/queue"
	projectRepository "github.com/PaulBarrie/infra-worker/pkg/repository/project"
	resourceRepository "github.com/PaulBarrie/infra-worker/pkg/repository/resource"
)

var (
	DefaultWorkdir   string = "/tmp/infra-worker"
	PluginRepository        = mongo.Client
)

func main() {
	ctx := context.Background()
	projectService := application.ProjectService{
		ProjectRepository: projectRepository.ProjectRepository,
	}
	pluginService := application.PluginService{}
	resourceService := application.ResourceService{
		ProjectRepository:  projectRepository.ProjectRepository,
		ResourceRepository: resourceRepository.ResourceRepository,
	}
	server := http.New(&projectService, &pluginService, &resourceService)
	queue.Queue.SetServices(&resourceService)
	go queue.Queue.StartConsumer(ctx)
	defer queue.Close()
	server.Start()
}
