package resource

import (
	"context"
	"fmt"
	dtoResource "github.com/athmos-cloud/infra-worker-athmos/pkg/common/dto/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/dao/helm"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/dao/kubernetes"
	kubernetesData "github.com/athmos-cloud/infra-worker-athmos/pkg/data/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/metadata"
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

func (repository *Repository) Create(ctx context.Context, opt option.Option) (interface{}, errors.Error) {
	if !opt.SetType(reflect.TypeOf(CreateRequest{}).String()).Validate() {
		return nil, errors.InvalidArgument.WithMessage(
			fmt.Sprintf(
				"Invalid argument type, expected CreateRequest, got %v",
				reflect.TypeOf(CreateRequest{}).Kind(), opt.Value,
			),
		)
	}
	request := opt.Value.(CreateRequest)
	logger.Info.Printf("Creating resource %s to cloud provider %s", request.ResourceType, request.ProviderType)

	curResource := resource.Factory(request.ResourceType)
	// Execute curPlugin
	curResource, err := curResource.New(request.Identifier, request.ProviderType)
	if !err.IsOk() {
		return nil, err
	}
	pluginReference, err := curResource.GetPluginReference()
	if !err.IsOk() {
		return nil, err
	}

	completedPayload, err := pluginReference.Plugin.ValidateAndCompletePluginEntry(request.ResourceSpecs)
	if !err.IsOk() {
		return nil, err
	}
	curResource.FromMap(request.ResourceSpecs)
	curResource.SetMetadata(metadata.CreateMetadataRequest{
		Name:             completedPayload["name"].(string),
		ProjectNamespace: request.Project.Namespace,
		NotMonitored:     !(completedPayload["monitored"].(bool)),
		Tags:             completedPayload["tags"].(map[string]string),
	})
	updatedStatus := curResource.GetStatus()

	resp, err := repository.HelmDAO.Create(ctx, option.Option{
		Value: helm.CreateHelmReleaseRequest{
			ReleaseName:  updatedStatus.HelmRelease.Name,
			ChartName:    pluginReference.ChartReference.ChartName,
			ChartVersion: pluginReference.ChartReference.ChartVersion,
			Values:       completedPayload,
			Namespace:    request.Project.Namespace,
		},
	})
	if !err.IsOk() {
		logger.Info.Printf("Error creating curResource %s", curResource.GetMetadata().Name)
		return nil, err
	}
	releaseResp := resp.(helm.CreateHelmReleaseResponse).Release
	// Parse manifest
	resID, err := kubernetesData.GetResourcesIdentifiersFromManifests(releaseResp.Manifest)
	if !err.IsOk() {
		logger.Info.Printf("Error parsing manifest for curResource %s", curResource.GetMetadata().Name)
		return nil, err
	}
	updatedStatus.KubernetesResources = kubernetesData.NewResourceList(resID)
	curResource.SetStatus(updatedStatus)
	curResource.FromMap(releaseResp.Config)
	// Insert curResource into project
	if errInsert := request.Project.Insert(curResource); errInsert.IsOk() {
		return nil, errInsert
	}

	return dtoResource.CreateResourceResponse{
		Resource: curResource,
		Project:  request.Project,
	}, errors.Created
}

func (repository *Repository) Get(_ context.Context, opt option.Option) (interface{}, errors.Error) {
	if !opt.SetType(reflect.TypeOf(GetRequest{}).String()).Validate() {
		return nil, errors.InvalidArgument.WithMessage(
			fmt.Sprintf(
				"Invalid argument type, expected GetRequest, got %v", reflect.TypeOf(opt.Value).Kind(),
			),
		)
	}
	request := opt.Value.(GetRequest)
	res, err := request.Project.Get(request.ResourceID)
	if !err.IsOk() {
		return nil, err
	}
	return GetResourceResponse{
		Resource: res,
	}, errors.OK
}

func (repository *Repository) Watch(ctx context.Context, opt option.Option) (interface{}, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (repository *Repository) List(ctx context.Context, opt option.Option) (interface{}, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (repository *Repository) Update(ctx context.Context, opt option.Option) errors.Error {
	if !opt.SetType(reflect.TypeOf(UpdateRequest{}).String()).Validate() {
		return errors.InvalidArgument.WithMessage(
			fmt.Sprintf(
				"Invalid argument type, expected UpdateRequest, got %v", reflect.TypeOf(opt.Value).Kind(),
			),
		)
	}
	request := opt.Value.(UpdateRequest)
	// Get curResource
	res, err := request.Project.Get(request.ResourceID)
	if !err.IsOk() {
		return err
	}
	// Get plugin
	pluginReference, err := res.GetPluginReference()
	if !err.IsOk() {
		return err
	}
	completedPayload, err := pluginReference.Plugin.ValidateAndCompletePluginEntry(request.NewResourceSpecs)
	if !err.IsOk() {
		return err
	}
	updatedStatus := res.GetStatus()

	err = repository.HelmDAO.Update(ctx, option.Option{
		Value: helm.UpdateHelmReleaseRequest{
			ReleaseName:  updatedStatus.HelmRelease.Name,
			ChartName:    pluginReference.ChartReference.ChartName,
			ChartVersion: pluginReference.ChartReference.ChartVersion,
			Values:       completedPayload,
			Namespace:    request.Project.Namespace,
		},
	})
	if !err.IsOk() {
		logger.Info.Printf("Error creating curResource %s", res.GetMetadata().Name)
		return err
	}
	return errors.OK
}

func (repository *Repository) Delete(ctx context.Context, opt option.Option) errors.Error {
	//TODO implement me
	panic("implement me")
}
