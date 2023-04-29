package application

import (
	"context"

	dto "github.com/athmos-cloud/infra-worker-athmos/pkg/common/dto/project"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	projectRepository "github.com/athmos-cloud/infra-worker-athmos/pkg/repository/project"
)

type ProjectService struct {
	ProjectRepository *projectRepository.Repository
}

func (ps *ProjectService) CreateProject(ctx context.Context, request dto.CreateProjectRequest) dto.CreateProjectResponse {
	resp := ps.ProjectRepository.Create(ctx, option.Option{
		Value: request,
	})

	return resp.(dto.CreateProjectResponse)
}

func (ps *ProjectService) UpdateProjectName(ctx context.Context, request dto.UpdateProjectRequest) {
	ps.ProjectRepository.Update(ctx, option.Option{
		Value: request,
	})
}

func (ps *ProjectService) GetProjectByID(ctx context.Context, request dto.GetProjectByIDRequest) dto.GetProjectByIDResponse {
	resp := ps.ProjectRepository.Get(ctx, option.Option{
		Value: request,
	})
	return resp.(dto.GetProjectByIDResponse)
}

func (ps *ProjectService) GetProjectByOwnerID(ctx context.Context, request dto.GetProjectByOwnerIDRequest) dto.GetProjectByOwnerIDResponse {
	resp := ps.ProjectRepository.List(ctx, option.Option{
		Value: request,
	})
	return resp.(dto.GetProjectByOwnerIDResponse)
}

func (ps *ProjectService) DeleteProject(ctx context.Context, request dto.DeleteRequest) {
	ps.ProjectRepository.Delete(ctx, option.Option{
		Value: request,
	})
}
