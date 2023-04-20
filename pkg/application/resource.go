package application

import (
	"context"
	dtoProject "github.com/PaulBarrie/infra-worker/pkg/common/dto/project"
	dtoResource "github.com/PaulBarrie/infra-worker/pkg/common/dto/resource"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/logger"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/option"
	"github.com/PaulBarrie/infra-worker/pkg/repository"
	resourceRepository "github.com/PaulBarrie/infra-worker/pkg/repository/resource"
	"helm.sh/helm/v3/pkg/release"
)

type ResourceService struct {
	ProjectRepository  repository.IRepository
	ResourceRepository repository.IRepository
}

func (service *ResourceService) CreateResource(ctx context.Context, payload dtoResource.CreateResourceRequest) (dtoResource.CreateResourceResponse, errors.Error) {
	// Get project
	response, err := service.ProjectRepository.Get(ctx, option.Option{
		Value: dtoProject.GetProjectByIDRequest{
			ProjectID: payload.ProjectID,
		},
	})
	currentProject := response.(dtoProject.GetProjectByIDResponse).Payload
	if !err.IsOk() {
		return dtoResource.CreateResourceResponse{}, err
	}
	logger.Info.Printf("Current project: %v", currentProject)
	// Install Resource
	resp, err := service.ResourceRepository.Create(ctx, option.Option{
		Value: resourceRepository.CreateRequest{
			ProjectNamespace: currentProject.Namespace,
			ProviderType:     payload.ProviderType,
			ResourceType:     payload.ResourceType,
			ResourceSpecs:    payload.ResourceSpecs,
		},
	})

	logger.Info.Printf("Response: %v", resp)
	if !err.IsOk() {
		logger.Info.Printf("Error creating resource : %v", err)
		return dtoResource.CreateResourceResponse{}, err
	}

	return dtoResource.CreateResourceResponse{ResourceID: "", HelmRelease: resp.(*release.Release)}, errors.OK
}

func (service *ResourceService) GetResource(ctx context.Context, payload dtoResource.GetResourceRequest) (dtoResource.CreateResourceResponse, errors.Error) {
	panic("")
}

func (service *ResourceService) UpdateResource(ctx context.Context, payload dtoResource.UpdateResourceRequest) errors.Error {
	panic("")
}

func (service *ResourceService) DeleteResource(ctx context.Context, payload dtoResource.DeleteResourceRequest) errors.Error {
	panic("")
}
