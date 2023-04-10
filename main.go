package main

import (
	"github.com/PaulBarrie/infra-worker/pkg/http"
	"github.com/PaulBarrie/infra-worker/pkg/queue"
	"github.com/PaulBarrie/infra-worker/pkg/repository/mongo"
)

var (
	DefaultWorkdir   string = "/tmp/infra-worker"
	PluginRepository        = mongo.Client
)

func main() {
	go queue.Queue.StartConsumer()
	defer queue.Close()
	http.Start()
}

//ctx := context.Background()
//logger.Info.Println("Starting infra-worker")
//projectService := project.Service{
//	ProjectRepository: mongo.Client,
//}
//resourceService := resource.Service{
//	ProjectRepository: mongo.Client,
//	PluginRepository:  helm.ReleaseClient,
//}
//
//resp, err := projectService.Create(ctx, projectDto.CreateProjectRequest{
//	ProjectName: "test",
//	OwnerID:     "test",
//},
//)
//if !err.IsOk() {
//	logger.Error.Printf("Error creating project: %s", err)
//} else {
//	logger.Info.Printf("Created project with ID: %s", resp.ProjectID)
//}
//gcpCreds, errFile := os.ReadFile(config.Current.Test.Credentials.GCP)
//if errFile != nil {
//	logger.Error.Printf("Error reading GCP credentials: %s", errFile)
//	return
//}
//res1, err := resourceService.CreateResource(ctx, resourceDto.CreateResourceRequest{
//	ProjectID:    "64312d1d9cc6ed73c5a26fc2",
//	ProviderType: common.GCP,
//	ResourceType: common.Provider,
//	ResourceSpecs: map[string]interface{}{
//		"name": "gcp-test",
//		"vpc":  "plugin-playground",
//		"auth": map[string]interface{}{
//			"tyoe": "secret",
//			"secret": map[string]interface{}{
//				"value": gcpCreds,
//			},
//		},
//	},
//})
//if !err.IsOk() {
//	logger.Error.Printf("Error creating resource: %s", err)
//} else {
//	logger.Info.Printf("Created resource with ID: %s", res1.ResourceID)
//}
//res1, err = resourceService.CreateResource(ctx, resourceDto.CreateResourceRequest{
//	ProjectID:    "64312d1d9cc6ed73c5a26fc2",
//	ProviderType: common.GCP,
//	ResourceType: common.VPC,
//	ResourceSpecs: map[string]interface{}{
//		"name": "test",
//	},
//})
//if !err.IsOk() {
//	logger.Error.Printf("Error creating resource: %s", err)
//} else {
//	logger.Info.Printf("Created resource with ID: %s", res1.ResourceID)
//}
