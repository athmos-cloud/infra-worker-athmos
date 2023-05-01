package main

import (
	"context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/application"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/application/secret"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/dao/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/exposition/http"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/exposition/queue"
	projectRepository "github.com/athmos-cloud/infra-worker-athmos/pkg/repository/project"
	resourceRepository "github.com/athmos-cloud/infra-worker-athmos/pkg/repository/resource"
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
	secretService := secret.Service{
		ProjectRepository: projectRepository.ProjectRepository,
		KubernetesDAO:     kubernetes.Client,
	}
	server := http.New(&projectService, &pluginService, &resourceService, &secretService)
	queue.Queue.SetServices(&resourceService)
	go queue.Queue.StartConsumer(ctx)
	defer queue.Close()
	server.Start()
}
