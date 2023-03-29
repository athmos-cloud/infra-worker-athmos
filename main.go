package main

import (
	"github.com/PaulBarrie/infra-worker/pkg/plugin"
	"github.com/PaulBarrie/infra-worker/pkg/repository/mongo"
)

var (
	DefaultWorkdir   string = "/tmp/infra-worker"
	PluginRepository        = mongo.Client
)

func main() {
	//_ := context.Background()
	pluginEntry := map[string]interface{}{
		"vpc":         "vpc-1234567890",
		"zone":        "us-east-1a",
		"network":     "network-1234567890",
		"subnetwork":  "subnet-1234567890",
		"machineType": "n1-standard-1",
		"disk": map[string]interface{}{
			"size":       10,
			"type":       "pd-standard",
			"mode":       "READ_WRITE",
			"autoDelete": true,
		},
		"os": map[string]interface{}{
			"type":    "ubuntu",
			"version": "20.04",
		},
	}
	plugin, err := plugin.Get("gcp", "vm")
	if !err.IsOk() {
		panic(err)
	}
	//logger.Info.Println(plugin.Types[0])
	if err1 := plugin.Validate(pluginEntry); err1.IsOk() {
		panic(err1)
	}

	return

}
