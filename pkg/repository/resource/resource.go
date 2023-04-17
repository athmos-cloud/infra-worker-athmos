package resource

import (
	"context"
	"fmt"
	dtoResource "github.com/PaulBarrie/infra-worker/pkg/common/dto/resource"
	"github.com/PaulBarrie/infra-worker/pkg/dao/helm"
	"github.com/PaulBarrie/infra-worker/pkg/dao/kubernetes"
	"github.com/PaulBarrie/infra-worker/pkg/domain"
	"github.com/PaulBarrie/infra-worker/pkg/domain/plugin"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/option"
	"reflect"
)

type Repository struct {
	KubernetesDAO kubernetes.DAO
	HelmDAO       helm.ReleaseDAO
}

func (repository *Repository) Create(ctx context.Context, optn option.Option) (interface{}, errors.Error) {
	if optn.SetType(reflect.TypeOf(dtoResource.CreateResourceRequest{}).String()).Validate() {
		return nil, errors.InvalidArgument.WithMessage(
			fmt.Sprintf(
				"Invalid argument type, expected %s, got %v",
				reflect.TypeOf(dtoResource.CreateResourceRequest{}).Kind(), optn.Value,
			),
		)
	}
	request := optn.Value.(dtoResource.CreateResourceRequest)
	curPlugin, err := plugin.Get(request.ProviderType, request.ResourceType)
	if !err.IsOk() {
		return dtoResource.CreateResourceResponse{}, err
	}
	completedPayload, err := curPlugin.ValidateAndCompletePluginEntry(request.ResourceSpecs)
	if !err.IsOk() {
		return dtoResource.CreateResourceResponse{}, err
	}
	resource := domain.Factory(request.ResourceType)
	// Execute curPlugin
	pluginReference, err := resource.GetPluginReference(dtoResource.GetPluginReferenceRequest{
		ProviderType: request.ProviderType,
	})
	if !err.IsOk() {
		return dtoResource.CreateResourceResponse{}, err
	}
	// Validator
	resource.FromMap(request.ResourceSpecs)
	resource.WithMetadata(domain.CreateMetadataRequest{
		Name:             completedPayload["name"].(string),
		ProjectNamespace: request.ProjectNamespace,
		NotMonitored:     !(completedPayload["monitored"].(bool)),
		Tags:             completedPayload["tags"].(map[string]string),
	})

	return repository.HelmDAO.Create(ctx, option.Option{
		Value: option.Option{
			Value: helm.CreateHelmReleaseRequest{
				ReleaseName:  resource.GetMetadata().Name,
				ChartName:    pluginReference.ChartName,
				ChartVersion: pluginReference.ChartVersion,
				Values:       completedPayload,
				Namespace:    request.ProjectNamespace,
			},
		},
	})
}

func (repository *Repository) Get(ctx context.Context, option option.Option) (interface{}, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (repository *Repository) Watch(ctx context.Context, option option.Option) (interface{}, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (repository *Repository) List(ctx context.Context, option option.Option) (interface{}, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (repository *Repository) Update(ctx context.Context, option option.Option) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (repository *Repository) Delete(ctx context.Context, option option.Option) errors.Error {
	//TODO implement me
	panic("implement me")
}
