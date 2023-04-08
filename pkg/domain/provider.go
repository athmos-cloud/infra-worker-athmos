package domain

import (
	"github.com/PaulBarrie/infra-worker/pkg/common"
	dto "github.com/PaulBarrie/infra-worker/pkg/common/dto/resource"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
)

type Provider struct {
	Metadata Metadata            `bson:"metadata"`
	Type     common.ProviderType `bson:"type"`
	Auth     Auth                `bson:"auth"`
	VPCs     []VPC               `bson:"vpcs"`
}

type ProviderList []Provider

func (provider *Provider) GetMetadata() Metadata {
	return provider.Metadata
}

func (provider *Provider) WithMetadata(request CreateMetadataRequest) {
	provider.Metadata = New(request)
}
func (provider *Provider) GetPluginReference(request dto.GetPluginReferenceRequest) (dto.GetPluginReferenceResponse, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (provider *Provider) FromMap(m map[string]interface{}) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (provider *Provider) InsertIntoProject(project Project, upsert bool) errors.Error {
	//TODO implement me
	panic("implement me")
}
