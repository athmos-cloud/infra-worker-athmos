package application

import (
	"context"
	dtoProject "github.com/athmos-cloud/infra-worker-athmos/pkg/common/dto/project"
	dtoResource "github.com/athmos-cloud/infra-worker-athmos/pkg/common/dto/resource"
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

	resp := service.ResourceRepository.Create(ctx, option.Option{
		Value: resourceRepository.CreateRequest{
			Project:       currentProject,
			ProviderType:  payload.ProviderType,
			ResourceType:  payload.ResourceType,
			ResourceSpecs: payload.ResourceSpecs,
		},
	})
	createdResource := resp.(resourceRepository.CreateResponse).Resource

	return dtoResource.CreateResourceResponse{Resource: createdResource}
}

func (service *ResourceService) GetResource(ctx context.Context, payload dtoResource.GetResourceRequest) dtoResource.CreateResourceResponse {
	panic("")
}

func (service *ResourceService) UpdateResource(ctx context.Context, payload dtoResource.UpdateResourceRequest) {
	panic("")
}

func (service *ResourceService) DeleteResource(ctx context.Context, payload dtoResource.DeleteResourceRequest) {
	panic("")
}
