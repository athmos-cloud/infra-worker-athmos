package resource

import (
	"context"
	"fmt"
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

func (repository *Repository) Create(ctx context.Context, opt option.Option) interface{} {
	if !opt.SetType(reflect.TypeOf(CreateRequest{}).String()).Validate() {
		panic(errors.InvalidArgument.WithMessage(
			fmt.Sprintf(
				"Invalid argument type, expected CreateRequest, got %v",
				opt.Value,
			),
		))
	}
	request := opt.Value.(CreateRequest)
	logger.Info.Printf("Creating resource %s to cloud provider %s", request.ResourceType, request.ProviderType)

	curResource := resource.Factory(request.ResourceType)
	// Execute curPlugin
	curResource = curResource.New(request.Identifier, request.ProviderType)
	pluginReference := curResource.GetPluginReference()

	completedPayload, err := pluginReference.Plugin.ValidateAndCompletePluginEntry(request.ResourceSpecs)
	if !err.IsOk() {
		panic(err)
	}
	curResource.FromMap(completedPayload)
	logger.Info.Printf("Completed payload: %v", completedPayload)
	curResource.SetMetadata(metadata.CreateMetadataRequest{
		Name:             completedPayload["name"].(string),
		ProjectNamespace: request.Project.Namespace,
		NotMonitored:     !(completedPayload["monitored"].(bool)),
		Tags:             completedPayload["tags"].(map[string]string),
	})
	updatedStatus := curResource.GetStatus()

	resp := repository.HelmDAO.Create(ctx, option.Option{
		Value: helm.CreateHelmReleaseRequest{
			ReleaseName:  updatedStatus.HelmRelease.Name,
			ChartName:    pluginReference.ChartReference.ChartName,
			ChartVersion: pluginReference.ChartReference.ChartVersion,
			Values:       completedPayload,
			Namespace:    request.Project.Namespace,
		},
	})

	releaseResp := resp.(helm.CreateHelmReleaseResponse).Release
	// Parse manifest
	resID := kubernetesData.GetResourcesIdentifiersFromManifests(releaseResp.Manifest)

	updatedStatus.KubernetesResources = kubernetesData.NewResourceList(resID)
	curResource.SetStatus(updatedStatus)
	curResource.FromMap(releaseResp.Config)
	// Insert curResource into project
	request.Project.Insert(curResource)

	return CreateResponse{
		Resource: curResource,
	}
}

func (repository *Repository) Get(_ context.Context, opt option.Option) interface{} {
	if !opt.SetType(reflect.TypeOf(GetRequest{}).String()).Validate() {
		panic(errors.InvalidArgument.WithMessage(
			fmt.Sprintf(
				"Invalid argument type, expected GetRequest, got %v", reflect.TypeOf(opt.Value).Kind(),
			),
		))
	}
	request := opt.Value.(GetRequest)
	res := request.Project.Get(request.ResourceID)

	return GetResourceResponse{
		Resource: res,
	}
}

func (repository *Repository) Watch(ctx context.Context, opt option.Option) interface{} {
	//TODO implement me
	panic("implement me")
}

func (repository *Repository) List(ctx context.Context, opt option.Option) interface{} {
	//TODO implement me
	panic("implement me")
}

func (repository *Repository) Update(ctx context.Context, opt option.Option) interface{} {
	if !opt.SetType(reflect.TypeOf(UpdateRequest{}).String()).Validate() {
		panic(errors.InvalidArgument.WithMessage(
			fmt.Sprintf(
				"Invalid argument type, expected UpdateRequest, got %v", reflect.TypeOf(opt.Value).Kind(),
			),
		))
	}
	request := opt.Value.(UpdateRequest)
	res := request.Project.Get(request.ResourceID)
	pluginReference := res.GetPluginReference()

	completedPayload, err := pluginReference.Plugin.ValidateAndCompletePluginEntry(request.NewResourceSpecs)
	if !err.IsOk() {
		panic(err)
	}
	updatedStatus := res.GetStatus()

	repository.HelmDAO.Update(ctx, option.Option{
		Value: helm.UpdateHelmReleaseRequest{
			ReleaseName:  updatedStatus.HelmRelease.Name,
			ChartName:    pluginReference.ChartReference.ChartName,
			ChartVersion: pluginReference.ChartReference.ChartVersion,
			Values:       completedPayload,
			Namespace:    request.Project.Namespace,
		},
	})
	return nil
}

func (repository *Repository) Delete(ctx context.Context, opt option.Option) {
	if !opt.SetType(reflect.TypeOf(DeleteRequest{}).String()).Validate() {
		panic(errors.InvalidArgument.WithMessage(
			fmt.Sprintf(
				"Invalid argument type, expected DeleteRequest, got %v", reflect.TypeOf(opt.Value).Kind(),
			),
		))
	}
	request := opt.Value.(DeleteRequest)
	currentResource := request.Project.Get(request.ResourceID)
	//Uninstall helm release
	repository.HelmDAO.Delete(ctx, option.Option{
		Value: helm.DeleteHelmReleaseRequest{
			ReleaseName: currentResource.GetStatus().HelmRelease.Name,
			Namespace:   request.Project.Namespace,
		},
	})
	request.Project.Delete(currentResource)
}
