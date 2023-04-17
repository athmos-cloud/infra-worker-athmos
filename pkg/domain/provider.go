package domain

import (
	"fmt"
	"github.com/PaulBarrie/infra-worker/pkg/common"
	dto "github.com/PaulBarrie/infra-worker/pkg/common/dto/resource"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/config"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/utils"
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
	switch request.ProviderType {
	case common.GCP:
		return dto.GetPluginReferenceResponse{
			ChartName:    config.Current.Plugins.Crossplane.GCP.Provider.ChartName,
			ChartVersion: config.Current.Plugins.Crossplane.GCP.Provider.Version,
		}, errors.Error{}
	}
	return dto.GetPluginReferenceResponse{}, errors.InvalidArgument.WithMessage(fmt.Sprintf("provider type %s not supported", request.ProviderType))
}

func (provider *Provider) FromMap(m map[string]interface{}) errors.Error {
	*provider = Provider{}
	if m["id"] == nil {
		provider.Metadata.ID = utils.GenerateUUID()
	} else {
		provider.Metadata.ID = m["id"].(string)
	}
	if m["name"] == nil {
		return errors.InvalidArgument.WithMessage("name is required")
	}
	provider.Metadata.Name = m["name"].(string)
	return errors.OK
}

func (provider *Provider) InsertIntoProject(project Project, upsert bool) errors.Error {
	//TODO implement me
	panic("implement me")
}
