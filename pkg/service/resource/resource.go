package resource

import (
	"context"
	dto "github.com/PaulBarrie/infra-worker/pkg/common/dto/resource"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/config"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/option"
	"github.com/PaulBarrie/infra-worker/pkg/repository"
	"github.com/PaulBarrie/infra-worker/pkg/repository/mongo"
	"github.com/PaulBarrie/infra-worker/pkg/resource"
)

type Service struct {
	ProjectRepository repository.IRepository
	PluginRepository  repository.IRepository
}

func (service *Service) CreateResource(ctx context.Context, payload dto.CreateResourceRequest) (dto.CreateResourceResponse, errors.Error) {
	// Get project
	currentProject, err := service.ProjectRepository.Get(ctx, option.Option{
		Value: mongo.GetRequest{
			CollectionName: config.Current.Mongo.ProjectCollection,
			Id:             payload.ProjectID,
		},
	})
	if !err.IsOk() {
		return dto.CreateResourceResponse{}, err
	}
	resource := resource.ResourceFactory(payload.ResourceType)
	resource.Create(ctx, option.Option{
		Value: resource.CreateResourceRequest{
		}
	})
	// Execute plugin

	// Save resource

	panic("")
}

func (service *Service) GetResource(payload resource.GetResourceRequest) (resource.CreateResourceResponse, errors.Error) {
	panic("")
}

func (service *Service) UpdateResource(payload resource.UpdateResourceRequest) errors.Error {
	panic("")
}

func (service *Service) DeleteResource(payload resource.DeleteResourceRequest) errors.Error {
	panic("")
}
