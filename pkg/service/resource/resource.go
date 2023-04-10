package resource

import (
	"context"
	"fmt"
	dto "github.com/PaulBarrie/infra-worker/pkg/common/dto/resource"
	"github.com/PaulBarrie/infra-worker/pkg/domain"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/config"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/logger"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/option"
	"github.com/PaulBarrie/infra-worker/pkg/plugin"
	"github.com/PaulBarrie/infra-worker/pkg/repository"
	"github.com/PaulBarrie/infra-worker/pkg/repository/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type Service struct {
	ProjectRepository repository.IRepository
	PluginRepository  repository.IRepository
}

func (service *Service) CreateResource(ctx context.Context, payload dto.CreateResourceRequest) (dto.CreateResourceResponse, errors.Error) {
	// Get project
	projectResp, err := service.ProjectRepository.Get(ctx, option.Option{
		Value: mongo.GetRequest{
			CollectionName: config.Current.Mongo.ProjectCollection,
			Id:             payload.ProjectID,
		},
	})
	if !err.IsOk() && !err.IsNotFound() {
		return dto.CreateResourceResponse{}, err
	} else if err.IsNotFound() {
		return dto.CreateResourceResponse{}, err.WithMessage(fmt.Sprintf("Project with ID %s not found", payload.ProjectID))
	}
	response := projectResp.(mongo.GetResponse)
	var currentProject domain.Project
	if convErr := bson.Unmarshal(response.Payload, &currentProject); convErr != nil {
		return dto.CreateResourceResponse{}, errors.InternalError.WithMessage(convErr.Error())
	}
	curPlugin, err := plugin.Get(payload.ProviderType, payload.ResourceType)
	if !err.IsOk() {
		return dto.CreateResourceResponse{}, err
	}
	logger.Info.Printf("Plugin: %v", curPlugin)
	completedPayload, err := curPlugin.ValidateAndComplete(payload.ResourceSpecs)
	if !err.IsOk() {
		return dto.CreateResourceResponse{}, err
	}
	resource := domain.Factory(payload.ResourceType)
	logger.Info.Printf("Completed payload: %v", completedPayload)
	// Execute curPlugin
	pluginReference, err := resource.GetPluginReference(dto.GetPluginReferenceRequest{
		ProviderType: payload.ProviderType,
	})
	if !err.IsOk() {
		return dto.CreateResourceResponse{}, err
	}
	logger.Info.Printf("Plugin reference: %v", pluginReference)
	// Validator
	resource.FromMap(payload.ResourceSpecs)
	resource.WithMetadata(domain.CreateMetadataRequest{
		Name:             payload.ResourceSpecs["name"].(string),
		ProjectNamespace: currentProject.Namespace,
		NotMonitored:     !(payload.ResourceSpecs["notMonitored"].(bool)),
		Tags:             payload.ResourceSpecs["tags"].(map[string]string),
	})

	//_, err = service.PluginRepository.Create(ctx, option.Option{
	//	Value: helm.CreateHelmReleaseRequest{
	//		ChartName:    pluginReference.ChartName,
	//		ChartVersion: pluginReference.ChartVersion,
	//		ReleaseName:  resource.GetMetadata().ReleaseReference.Name,
	//		Values:       completedPayload,
	//		Namespace:    resource.GetMetadata().ReleaseReference.Namespace,
	//	},
	//})
	//if !err.IsOk() {
	//	return dto.CreateResourceResponse{}, err
	//}

	//Map curPlugin input to struct
	//Save domain
	return dto.CreateResourceResponse{ResourceID: resource.GetMetadata().ID}, errors.OK

}

func (service *Service) GetResource(ctx context.Context, payload dto.GetResourceRequest) (dto.CreateResourceResponse, errors.Error) {
	panic("")
}

func (service *Service) UpdateResource(ctx context.Context, payload dto.UpdateResourceRequest) errors.Error {
	panic("")
}

func (service *Service) DeleteResource(ctx context.Context, payload dto.DeleteResourceRequest) errors.Error {
	panic("")
}
