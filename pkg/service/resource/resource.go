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
	//projectResp, err := service.ProjectRepository.Get(ctx, option.Option{
	//	Value: mongo.GetRequest{
	//		CollectionName: config.Current.Mongo.ProjectCollection,
	//		Id:             payload.ProjectID,
	//	},
	//})
	//currentProject := projectResp.(mongo.GetResponse).Payload.(project.Project)
	//if !err.IsOk() {
	//	return dto.CreateResourceResponse{}, err
	//}
	//plugin, err := plugin.Get(payload.ProviderType, payload.ResourceType)
	//if !err.IsOk() {
	//	return dto.CreateResourceResponse{}, err
	//}
	//factory := resource.Factory(payload.ResourceType)
	//// Execute plugin
	//
	//if !err.IsOk() {
	//	return dto.CreateResourceResponse{}, err
	//}
	//helmClient.Create(ctx, option.Option{
	//	Value: factory(payload.ResourceSpecs),
	//}

	// Map plugin input to struct
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
