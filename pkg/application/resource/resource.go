package resource

import (
	"context"
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/application/project"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/repository"
	projectRepository "github.com/athmos-cloud/infra-worker-athmos/pkg/repository/project"
	resourceRepository "github.com/athmos-cloud/infra-worker-athmos/pkg/repository/resource"
)

type Service struct {
	ProjectRepository  repository.IRepository
	ResourceRepository repository.IRepository
}

func (service *Service) CreateResource(ctx context.Context, payload CreateResourceRequest) CreateResourceResponse {
	// Get resource
	response := service.ProjectRepository.Get(ctx, option.Option{
		Value: project.GetProjectByIDRequest{
			ProjectID: payload.ProjectID,
		},
	})
	currentProject := response.(project.GetProjectByIDResponse).Payload
	id := identifier.NewID(payload.Identifier)
	if currentProject.Exists(id) {
		panic(errors.Conflict.WithMessage(fmt.Sprintf("Resource %s already exists", id)))
	}
	resp := service.ResourceRepository.Create(ctx, option.Option{
		Value: resourceRepository.CreateRequest{
			Project:       currentProject,
			Identifier:    id,
			ProviderType:  payload.ProviderType,
			ResourceType:  payload.ResourceType,
			ResourceSpecs: payload.ResourceSpecs,
		},
	})
	createdResource := resp.(resourceRepository.CreateResponse).Resource
	service.ProjectRepository.Update(ctx, option.Option{
		Value: projectRepository.UpdateProjectRequest{
			ProjectID:      payload.ProjectID,
			UpdatedProject: currentProject,
		},
	})

	return CreateResourceResponse{Resource: createdResource}
}

func (service *Service) GetResource(ctx context.Context, payload GetResourceRequest) CreateResourceResponse {
	resp := service.ResourceRepository.Get(ctx, option.Option{
		Value: project.GetProjectByIDRequest{
			ProjectID: payload.ProjectID,
		},
	})
	projectResponse := resp.(project.GetProjectByIDResponse).Payload
	resource := service.ResourceRepository.Get(ctx, option.Option{
		Value: resourceRepository.GetRequest{
			Project:    projectResponse,
			ResourceID: payload.ResourceID,
		},
	})
	// Return domain
	return CreateResourceResponse{Resource: resource.(resourceRepository.GetResourceResponse).Resource}
}

func (service *Service) UpdateResource(ctx context.Context, payload UpdateResourceRequest) {
	projectResponse := service.ProjectRepository.Get(ctx, option.Option{
		Value: project.GetProjectByIDRequest{
			ProjectID: payload.ProjectID,
		},
	})
	currentProject := projectResponse.(projectRepository.GetProjectByIDResponse).Payload
	id := identifier.NewID(payload.ResourceID)
	if currentProject.Exists(id) {
		panic(errors.Conflict.WithMessage(fmt.Sprintf("Resource %s already exists", id)))
	}
	service.ResourceRepository.Update(ctx, option.Option{
		Value: resourceRepository.UpdateRequest{
			Project:    currentProject,
			ResourceID: id,
		},
	})
	service.ProjectRepository.Update(ctx, option.Option{
		Value: projectRepository.UpdateProjectRequest{
			ProjectID:      payload.ProjectID,
			UpdatedProject: currentProject,
		},
	})
}

func (service *Service) DeleteResource(ctx context.Context, payload DeleteResourceRequest) {
	project := service.ProjectRepository.Get(ctx, option.Option{
		Value: project.GetProjectByIDRequest{
			ProjectID: payload.ProjectID,
		},
	})
	currentProject := project.(projectRepository.GetProjectByIDResponse).Payload
	id := identifier.NewID(payload.ResourceID)
	service.ResourceRepository.Delete(ctx, option.Option{
		Value: resourceRepository.DeleteRequest{
			Project:    project.(projectRepository.GetProjectByIDResponse).Payload,
			ResourceID: id,
		},
	})
	service.ProjectRepository.Update(ctx, option.Option{
		Value: projectRepository.UpdateProjectRequest{
			ProjectID:      payload.ProjectID,
			UpdatedProject: currentProject,
		},
	})
}
