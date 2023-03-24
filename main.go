package main

import (
	"context"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/logger"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/option"
	"github.com/PaulBarrie/infra-worker/pkg/repository/helm"
	"github.com/PaulBarrie/infra-worker/pkg/repository/mongo"
)

var (
	DefaultWorkdir   string = "/tmp/infra-worker"
	PluginRepository        = mongo.Client
)

func main() {
	ctx := context.Background()
	client, err := helm.Client("crossplane-system")
	if !err.IsOk() {
		logger.Error.Println("Error: ", err)
	}
	release, err := client.Create(ctx, option.Option{
		Value: helm.CreateHelmReleaseRequest{
			ChartName:    "plugins/gcp-vpc",
			ChartVersion: "0.1.0",
			ReleaseName:  "vpc-test",
			Namespace:    "default",
			Values: map[string]interface{}{
				"name":    "test",
				"managed": true,
			},
		},
	},
	)
	resp, err := client.Get(ctx, option.Option{
		Value: helm.GetHelmReleaseRequest{
			ReleaseName: "vpc-test",
		},
	})
	getResp := resp.(helm.GetHelmReleaseResponse)
	logger.Info.Println("Info: ", release)

	if !err.IsOk() {
		logger.Error.Printf("Error loading chart :  %v", err)
	}
	logger.Info.Println("Info: ", getResp.Release.Manifest)
}
