package service

import (
	dto2 "github.com/PaulBarrie/infra-worker/pkg/application/dto"
	"github.com/PaulBarrie/infra-worker/pkg/infrastructure/repository"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
)

type ResourceService struct {
	ProjectRepository *repository.IRepository
	PluginRepository  *repository.IRepository
}

func (rs *ResourceService) CreateResource(payload dto2.CreateResourceRequest) (dto2.CreateResourceResponse, errors.Error) {
	panic("")
}

func (rs *ResourceService) GetResource(payload dto2.GetResourceRequest) (dto2.CreateResourceResponse, errors.Error) {
	panic("")
}

func (rs *ResourceService) UpdateResource(payload dto2.UpdateResourceRequest) errors.Error {
	panic("")
}

func (rs *ResourceService) DeleteResource(payload dto2.DeleteResourceRequest) errors.Error {
	panic("")
}
