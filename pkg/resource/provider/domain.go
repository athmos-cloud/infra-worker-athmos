package provider

import (
	"github.com/PaulBarrie/infra-worker/pkg/application/dto"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
	"github.com/PaulBarrie/infra-worker/pkg/resource/provider/auth"
)

type ProviderType string

const (
	AWS   ProviderType = "aws"
	AZURE ProviderType = "azure"
	GCP   ProviderType = "gcp"
)

type Provider struct {
	Type ProviderType `bson:"type"`
	Auth auth.Auth    `bson:"auth"`
}

func (provider *Provider) Create(request dto.CreateResourceRequest) (dto.CreateResourceResponse, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (provider *Provider) Update(request dto.UpdateResourceRequest) (dto.UpdateResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (provider *Provider) Get(request dto.GetResourceRequest) (dto.GetResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (provider *Provider) Watch(request dto.GetResourceRequest) (dto.GetResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (provider *Provider) List(request dto.GetListResourceRequest) (dto.GetListResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (provider *Provider) Delete(request dto.DeleteResourceRequest) (dto.DeleteResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}
