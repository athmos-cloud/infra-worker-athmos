package resource

import (
	"github.com/PaulBarrie/infra-worker/pkg/common/dto/resource"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
	"github.com/PaulBarrie/infra-worker/pkg/repository"
)

type ResourceService struct {
	ProjectRepository *repository.IRepository
	PluginRepository  *repository.IRepository
}

func (rs *ResourceService) CreateResource(payload resource.CreateResourceRequest) (resource.CreateResourceResponse, errors.Error) {
	panic("")
}

func (rs *ResourceService) GetResource(payload resource.GetResourceRequest) (resource.CreateResourceResponse, errors.Error) {
	panic("")
}

func (rs *ResourceService) UpdateResource(payload resource.UpdateResourceRequest) errors.Error {
	panic("")
}

func (rs *ResourceService) DeleteResource(payload resource.DeleteResourceRequest) errors.Error {
	panic("")
}
