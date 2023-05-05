package main

import (
	"context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/application/plugin"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/application/project"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/application/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/application/secret"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/dao/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/exposition/http"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/exposition/queue"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/logger"
	projectRepository "github.com/athmos-cloud/infra-worker-athmos/pkg/repository/project"
	resourceRepository "github.com/athmos-cloud/infra-worker-athmos/pkg/repository/resource"
)

func main() {
	logger.Info.Printf("Starting the go app...")
	ctx := context.Background()
	projectService := project.Service{
		ProjectRepository: projectRepository.ProjectRepository,
	}
	pluginService := plugin.Service{}
	resourceService := resource.Service{
		ProjectRepository:  projectRepository.ProjectRepository,
		ResourceRepository: resourceRepository.ResourceRepository,
	}
	secretService := secret.Service{
		ProjectRepository: projectRepository.ProjectRepository,
		KubernetesDAO:     kubernetes.Client,
	}
	defer func() {
		logger.Info.Printf("Stop go app")
		if r := recover(); r != nil {
			logger.Error.Printf("Error: %v", r)
		}
	}()
	server := http.New(&projectService, &pluginService, &resourceService, &secretService)
	queue.Queue.SetServices(&resourceService)
	go queue.Queue.StartConsumer(ctx)
	defer queue.Close()
	server.Start()
}
