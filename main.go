package main

import (
	"context"
	"github.com/PaulBarrie/infra-worker/pkg/common/dto/project"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/logger"
	"github.com/PaulBarrie/infra-worker/pkg/repository/mongo"
	projectService "github.com/PaulBarrie/infra-worker/pkg/service/project"
)

var (
	DefaultWorkdir   string = "/tmp/infra-worker"
	PluginRepository        = mongo.Client
)

func main() {
	ctx := context.Background()
	service := projectService.ProjectService{
		ProjectRepository: mongo.Client,
	}

	id1, _ := service.Create(ctx, project.CreateProjectRequest{
		ProjectName: "test1",
		OwnerID:     "toto",
	})
	_, _ = service.Create(ctx, project.CreateProjectRequest{
		ProjectName: "test2",
		OwnerID:     "toto",
	})
	projectByID, err := service.GetByID(ctx, project.GetProjectByIDRequest{
		ProjectID: id1.ProjectID,
	})
	if !err.IsOk() {
		logger.Info.Println("Err: ", err)
	}
	logger.Info.Println(ctx, "Project found with id: ", projectByID)
	err = service.Update(ctx, project.UpdateProjectRequest{
		ProjectID:   id1.ProjectID,
		ProjectName: "test1-updated",
	})
	if !err.IsOk() {
		logger.Error.Println(ctx, "Error: ", err)
	}
	projectAll, err := service.GetByOwnerID(ctx, project.GetProjectByOwnerIDRequest{
		OwnerID: "toto",
	})

	err = service.Delete(ctx, project.DeleteRequest{
		ProjectID: id1.ProjectID,
	})
	if !err.IsOk() {
		logger.Error.Println(ctx, "Error: ", err)
	}
	err = service.Delete(ctx, project.DeleteRequest{
		ProjectID: id1.ProjectID,
	})
	if !err.IsOk() {
		logger.Error.Println(ctx, "Error: ", err)
	}
	logger.Info.Println(ctx, "Project found with ownerID: ", projectAll)

}
