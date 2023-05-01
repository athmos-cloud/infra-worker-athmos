package application

import (
	"context"
	"fmt"
	dtoProject "github.com/athmos-cloud/infra-worker-athmos/pkg/common/dto/project"
	dtoResource "github.com/athmos-cloud/infra-worker-athmos/pkg/common/dto/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/repository"
	resourceRepository "github.com/athmos-cloud/infra-worker-athmos/pkg/repository/resource"
)

type ResourceService struct {
	ProjectRepository  repository.IRepository
	ResourceRepository repository.IRepository
}

func (service *ResourceService) CreateResource(ctx context.Context, payload dtoResource.CreateResourceRequest) dtoResource.CreateResourceResponse {
	// Get resource
	response := service.ProjectRepository.Get(ctx, option.Option{
		Value: dtoProject.GetProjectByIDRequest{
			ProjectID: payload.ProjectID,
		},
	})
	currentProject := response.(dtoProject.GetProjectByIDResponse).Payload
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
		Value: dtoProject.UpdateProjectRequest{
			ProjectID:      payload.ProjectID,
			UpdatedProject: currentProject,
		},
	})

	return dtoResource.CreateResourceResponse{Resource: createdResource}
}

func (service *ResourceService) GetResource(ctx context.Context, payload dtoResource.GetResourceRequest) dtoResource.CreateResourceResponse {
	resp := service.ResourceRepository.Get(ctx, option.Option{
		Value: dtoProject.GetProjectByIDRequest{
			ProjectID: payload.ProjectID,
		},
	})
	project := resp.(dtoProject.GetProjectByIDResponse).Payload
	resource := service.ResourceRepository.Get(ctx, option.Option{
		Value: resourceRepository.GetRequest{
			Project:    project,
			ResourceID: payload.ResourceID,
		},
	})
	// Return domain
	return dtoResource.CreateResourceResponse{Resource: resource.(resourceRepository.GetResourceResponse).Resource}
}

func (service *ResourceService) UpdateResource(ctx context.Context, payload dtoResource.UpdateResourceRequest) {
	project := service.ProjectRepository.Get(ctx, option.Option{
		Value: dtoProject.GetProjectByIDRequest{
			ProjectID: payload.ProjectID,
		},
	})
	currentProject := project.(dtoProject.GetProjectByIDResponse).Payload
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
		Value: dtoProject.UpdateProjectRequest{
			ProjectID:      payload.ProjectID,
			UpdatedProject: currentProject,
		},
	})
}

func (service *ResourceService) DeleteResource(ctx context.Context, payload dtoResource.DeleteResourceRequest) {
	project := service.ProjectRepository.Get(ctx, option.Option{
		Value: dtoProject.GetProjectByIDRequest{
			ProjectID: payload.ProjectID,
		},
	})
	currentProject := project.(dtoProject.GetProjectByIDResponse).Payload
	id := identifier.NewID(payload.ResourceID)
	service.ResourceRepository.Delete(ctx, option.Option{
		Value: resourceRepository.DeleteRequest{
			Project:    project.(dtoProject.GetProjectByIDResponse).Payload,
			ResourceID: id,
		},
	})
	service.ProjectRepository.Update(ctx, option.Option{
		Value: dtoProject.UpdateProjectRequest{
			ProjectID:      payload.ProjectID,
			UpdatedProject: currentProject,
		},
	})
}
