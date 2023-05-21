package main

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/http"
	registry2 "github.com/athmos-cloud/infra-worker-athmos/pkg/registry"
)

func main() {
	registry := registry2.NewRegistry()
	ctrl := registry.NewAppController()
	server := http.New(ctrl.Project, ctrl.Secret)
	//server := http.New(&projectService, &pluginService, &resourceService, &secretService)
	//rabbitmq.Queue.SetServices(&resourceService)
	//go rabbitmq.Queue.StartConsumer(ctx)
	//defer rabbitmq.Close()
	server.Start()
}
