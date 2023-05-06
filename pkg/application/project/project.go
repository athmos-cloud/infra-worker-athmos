package project

import (
	"context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	projectRepository "github.com/athmos-cloud/infra-worker-athmos/pkg/repository/project"
)

type Service struct {
	ProjectRepository *projectRepository.Repository
}

func (ps *Service) CreateProject(ctx context.Context, request CreateProjectRequest) CreateProjectResponse {
	resp := ps.ProjectRepository.Create(ctx, option.Option{
		Value: projectRepository.CreateProjectRequest{
			ProjectName: request.ProjectName,
			OwnerID:     request.OwnerID,
		},
	})

	projectRespo := resp.(projectRepository.CreateProjectResponse)

	return CreateProjectResponse{
		ProjectID: projectRespo.ProjectID,
	}
}

func (ps *Service) UpdateProjectName(ctx context.Context, request UpdateProjectRequest) {
	ps.ProjectRepository.Update(ctx, option.Option{
		Value: projectRepository.UpdateProjectRequest{
			ProjectID:   request.ProjectID,
			ProjectName: request.ProjectName,
		},
	})
}

func (ps *Service) GetProjectByID(ctx context.Context, request GetProjectByIDRequest) GetProjectByIDResponse {
	resp := ps.ProjectRepository.Get(ctx, option.Option{
		Value: projectRepository.GetProjectByIDRequest{
			ProjectID: request.ProjectID,
		},
	})
	projectResp := resp.(projectRepository.GetProjectByIDResponse)
	return GetProjectByIDResponse{
		Payload: projectResp.Payload,
	}
}

func (ps *Service) GetProjectByOwnerID(ctx context.Context, request GetProjectByOwnerIDRequest) GetProjectByOwnerIDResponse {
	resp := ps.ProjectRepository.List(ctx, option.Option{
		Value: projectRepository.GetProjectByOwnerIDRequest{
			OwnerID: request.OwnerID,
		},
	})
	projectList := resp.(projectRepository.GetProjectByOwnerIDResponse).Projects
	return toGetProjectByOwnerIDResponse(projectList)
}

func (ps *Service) DeleteProject(ctx context.Context, request DeleteRequest) {
	ps.ProjectRepository.Delete(ctx, option.Option{
		Value: projectRepository.DeleteRequest{
			ProjectID: request.ProjectID,
		},
	})
}
