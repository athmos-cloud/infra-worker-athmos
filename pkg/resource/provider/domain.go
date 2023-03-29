package provider

import (
	"github.com/PaulBarrie/infra-worker/pkg/common"
	resource2 "github.com/PaulBarrie/infra-worker/pkg/common/dto/resource"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
	"github.com/PaulBarrie/infra-worker/pkg/resource/provider/auth"
	"github.com/PaulBarrie/infra-worker/pkg/resource/vpc"
)

type Provider struct {
	Id   string              `bson:"id"`
	Name string              `bson:"name"`
	Type common.ProviderType `bson:"type"`
	Auth auth.Auth           `bson:"auth"`
	VPCs []vpc.VPC           `bson:"vpcs"`
}

func (provider *Provider) Create(request resource2.CreateResourceRequest) (resource2.CreateResourceResponse, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (provider *Provider) Update(request resource2.UpdateResourceRequest) (resource2.UpdateResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (provider *Provider) Get(request resource2.GetResourceRequest) (resource2.GetResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (provider *Provider) Watch(request resource2.GetResourceRequest) (resource2.GetResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (provider *Provider) List(request resource2.GetListResourceRequest) (resource2.GetListResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (provider *Provider) Delete(request resource2.DeleteResourceRequest) (resource2.DeleteResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}
