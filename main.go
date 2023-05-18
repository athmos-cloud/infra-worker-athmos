package main

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/presenter"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/http"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase"
)

func main() {
	server := http.New(
		controller.NewProjectController(
			usecase.NewProjectUseCase(
				&repository.Project{},
			),
			&presenter.Project{},
		),
	)
	//server := http.New(&projectService, &pluginService, &resourceService, &secretService)
	//rabbitmq.Queue.SetServices(&resourceService)
	//go rabbitmq.Queue.StartConsumer(ctx)
	//defer rabbitmq.Close()
	server.Start()
}
