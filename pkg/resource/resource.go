package resource

import (
	"github.com/PaulBarrie/infra-worker/pkg/application/dto"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
)

type IResource interface {
	Create(dto.CreateResourceRequest) (dto.CreateResourceResponse, errors.Error)
	Update(dto.UpdateResourceRequest) (dto.UpdateResourceRequest, errors.Error)
	Get(dto.GetResourceRequest) (dto.GetResourceRequest, errors.Error)
	Watch(dto.GetResourceRequest) (dto.GetResourceRequest, errors.Error)
	List(dto.GetListResourceRequest) (dto.GetListResourceRequest, errors.Error)
	Delete(dto.DeleteResourceRequest) (dto.DeleteResourceRequest, errors.Error)
}
