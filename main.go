package main

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/http"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/rabbitmq"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/config"
	registry2 "github.com/athmos-cloud/infra-worker-athmos/pkg/registry"
)

func main() {
	registry := registry2.NewRegistry()
	ctrl := registry.NewAppController()
	server := http.New(ctrl.Project, ctrl.Secret, ctrl.Resource)
	ctx := rabbitmq.NewContext()

	rabbitMQ := rabbitmq.New(
		config.Current.Queue.IncomingQueue,
		config.Current.Queue.OutcomingQueue,
		ctrl.Resource,
	)
	go rabbitMQ.StartConsumer(ctx)
	defer rabbitMQ.Close()
	server.Start()
}
