package resource

import (
	resource2 "github.com/PaulBarrie/infra-worker/pkg/common/dto/resource"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
)

type IResource interface {
	Create(resource2.CreateResourceRequest) (resource2.CreateResourceResponse, errors.Error)
	Update(resource2.UpdateResourceRequest) (resource2.UpdateResourceRequest, errors.Error)
	Get(resource2.GetResourceRequest) (resource2.GetResourceRequest, errors.Error)
	Watch(resource2.GetResourceRequest) (resource2.GetResourceRequest, errors.Error)
	List(resource2.GetListResourceRequest) (resource2.GetListResourceRequest, errors.Error)
	Delete(resource2.DeleteResourceRequest) (resource2.DeleteResourceRequest, errors.Error)
}
