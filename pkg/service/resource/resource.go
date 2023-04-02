package resource

import (
	"context"
	dto "github.com/PaulBarrie/infra-worker/pkg/common/dto/resource"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
	"github.com/PaulBarrie/infra-worker/pkg/repository"
)

type Service struct {
	ProjectRepository repository.IRepository
	PluginRepository  repository.IRepository
}

func (service *Service) CreateResource(ctx context.Context, payload dto.CreateResourceRequest) (dto.CreateResourceResponse, errors.Error) {
	// Get project
	//currentProject, err := service.ProjectRepository.Get(ctx, option.Option{
	//	Value: mongo.GetRequest{
	//		CollectionName: config.Current.Mongo.ProjectCollection,
	//		Id:             payload.ProjectID,
	//	},
	//})
	//if !err.IsOk() {
	//	return dto.CreateResourceResponse{}, err
	//}
	//_ := resource.Factory(payload.ResourceType)

	// Execute plugin

	// Save resource

	panic("")
}

func (service *Service) GetResource(payload dto.GetResourceRequest) (dto.CreateResourceResponse, errors.Error) {
	panic("")
}

func (service *Service) UpdateResource(payload dto.UpdateResourceRequest) errors.Error {
	panic("")
}

func (service *Service) DeleteResource(payload dto.DeleteResourceRequest) errors.Error {
	panic("")
}
