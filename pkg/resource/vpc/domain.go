package vpc

import (
	resource2 "github.com/PaulBarrie/infra-worker/pkg/common/dto/resource"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
	"github.com/PaulBarrie/infra-worker/pkg/resource"
	"github.com/PaulBarrie/infra-worker/pkg/resource/network"
)

type VPC struct {
	ID                string `bson:"_id,omitempty"`
	ResourceReference resource.Reference
	Networks          []network.Network `bson:"networks"`
}

func (vpc *VPC) Create(request resource2.CreateResourceRequest) (resource2.CreateResourceResponse, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (vpc *VPC) Update(request resource2.UpdateResourceRequest) (resource2.UpdateResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (vpc *VPC) Get(request resource2.GetResourceRequest) (resource2.GetResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (vpc *VPC) Watch(request resource2.GetResourceRequest) (resource2.GetResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (vpc *VPC) List(request resource2.GetListResourceRequest) (resource2.GetListResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (vpc *VPC) Delete(request resource2.DeleteResourceRequest) (resource2.DeleteResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}
