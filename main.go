package main

import (
	"context"
	helmrepo "github.com/PaulBarrie/infra-worker/pkg/infrastructure/repository/helm"
	"github.com/PaulBarrie/infra-worker/pkg/infrastructure/repository/mongo"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/logger"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/option"
)

var (
	DefaultWorkdir   string = "/tmp/infra-worker"
	PluginRepository        = mongo.Client
)

func main() {
	ctx := context.Background()
	helm := helmrepo.Connect()
	logger.Info.Println("Info: ", helm)
	chart, err := helm.Get(ctx, option.Option{
		Value: helmrepo.GetHelmChartRequest{ChartName: "gcp-vpc", ChartVersion: "0.1.0"}},
	)
	if !err.IsOk() {
		logger.Error.Printf("Error loading chart :  %v", err)
	}
	logger.Info.Println("Info: ", chart)
}
