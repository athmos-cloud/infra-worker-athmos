package main

import (
	"context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/application"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/dao/mongo"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/exposition/http"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/exposition/queue"
	projectRepository "github.com/athmos-cloud/infra-worker-athmos/pkg/repository/project"
	resourceRepository "github.com/athmos-cloud/infra-worker-athmos/pkg/repository/resource"
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
