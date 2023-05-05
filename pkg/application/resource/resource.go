package resource

import (
	"context"
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/repository"
	projectRepository "github.com/athmos-cloud/infra-worker-athmos/pkg/repository/project"
	resourceRepository "github.com/athmos-cloud/infra-worker-athmos/pkg/repository/resource"
	"reflect"
)

type Service struct {
	ProjectRepository  repository.IRepository
	ResourceRepository repository.IRepository
}

func (service *Service) CreateResource(ctx context.Context, payload CreateResourceRequest) CreateResourceResponse {
	// Get resource
	response := service.ProjectRepository.Get(ctx, option.Option{
		Value: projectRepository.GetProjectByIDRequest{
			ProjectID: payload.ProjectID,
		},
	})
	currentProject := response.(projectRepository.GetProjectByIDResponse).Payload
	if _, ok := payload.ResourceSpecs[identifier.IdentifierKey]; !ok || reflect.TypeOf(payload.ResourceSpecs[identifier.IdentifierKey]).Kind() != reflect.TypeOf(map[string]interface{}{}).Kind() {
		panic(errors.InvalidArgument.WithMessage("Missing identifier in resource specs"))
	}
	id := identifier.BuildFromMap(payload.ResourceSpecs[identifier.IdentifierKey].(map[string]interface{}))
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

	return CreateResourceResponse{Resource: domain.FromDataMapper(createdResource)}
}

func (service *Service) GetResource(ctx context.Context, payload GetResourceRequest) CreateResourceResponse {
	resp := service.ResourceRepository.Get(ctx, option.Option{
		Value: projectRepository.GetProjectByIDRequest{
			ProjectID: payload.ProjectID,
		},
	})
	projectResponse := resp.(projectRepository.GetProjectByIDResponse).Payload
	resource := service.ResourceRepository.Get(ctx, option.Option{
		Value: resourceRepository.GetRequest{
			Project:    projectResponse,
			ResourceID: payload.ResourceID,
		},
	})
	// Return domain
	dataResource := resource.(resourceRepository.GetResourceResponse).Resource
	return CreateResourceResponse{Resource: domain.FromDataMapper(dataResource)}
}

func (service *Service) UpdateResource(ctx context.Context, payload UpdateResourceRequest) {
	projectResponse := service.ProjectRepository.Get(ctx, option.Option{
		Value: projectRepository.GetProjectByIDRequest{
			ProjectID: payload.ProjectID,
		},
	})
	currentProject := projectResponse.(projectRepository.GetProjectByIDResponse).Payload
	id := identifier.Build(payload.ResourceID)
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
		Value: projectRepository.GetProjectByIDRequest{
			ProjectID: payload.ProjectID,
		},
	})
	currentProject := project.(projectRepository.GetProjectByIDResponse).Payload
	id := identifier.Build(payload.ResourceID)
	service.ResourceRepository.Delete(ctx, option.Option{
		Value: resourceRepository.DeleteRequest{
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
