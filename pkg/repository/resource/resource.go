package resource

import (
	"context"
	"fmt"
	dtoResource "github.com/athmos-cloud/infra-worker-athmos/pkg/common/dto/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/dao/helm"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/dao/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/plugin"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/logger"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	"reflect"
	"sync"
)

var ResourceRepository *Repository
var lock = &sync.Mutex{}

func init() {
	lock.Lock()
	defer lock.Unlock()
	if ResourceRepository == nil {
		ResourceRepository = &Repository{
			KubernetesDAO: kubernetes.Client,
			HelmDAO:       helm.ReleaseClient,
		}
	}
}

type Repository struct {
	KubernetesDAO *kubernetes.DAO
	HelmDAO       *helm.ReleaseDAO
}

func (repository *Repository) Create(ctx context.Context, optn option.Option) (interface{}, errors.Error) {
	if !optn.SetType(reflect.TypeOf(CreateRequest{}).String()).Validate() {
		return nil, errors.InvalidArgument.WithMessage(
			fmt.Sprintf(
				"Invalid argument type, expected CreateRequest, got %v",
				reflect.TypeOf(CreateRequest{}).Kind(), optn.Value,
			),
		)
	}
	logger.Info.Printf("Creating resource %v", optn.Value)
	request := optn.Value.(CreateRequest)
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
	logger.Info.Printf("Creating resource %s", pluginReference.ChartName)

	resp, err := repository.HelmDAO.Create(ctx, option.Option{
		Value: helm.CreateHelmReleaseRequest{
			ReleaseName:  resource.GetMetadata().Name,
			ChartName:    pluginReference.ChartName,
			ChartVersion: pluginReference.ChartVersion,
			Values:       completedPayload,
			Namespace:    request.ProjectNamespace,
		},
	})
	if !err.IsOk() {
		logger.Info.Printf("Error creating resource %s", resource.GetMetadata().Name)
		return dtoResource.CreateResourceResponse{}, err
	}
	// Parse manifest

	// Insert into project

	// Persist project

	return dtoResource.CreateResourceResponse{
		ResourceID:  resource.GetMetadata().ID,
		HelmRelease: resp.(helm.CreateHelmReleaseResponse).Release,
	}, errors.OK
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
	panic("implement me")
}

func (repository *Repository) Delete(ctx context.Context, option option.Option) errors.Error {
	//TODO implement me
	panic("implement me")
}
