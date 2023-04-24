package resources

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/common"
	dto "github.com/athmos-cloud/infra-worker-athmos/pkg/common/dto/resource"
	domain2 "github.com/athmos-cloud/infra-worker-athmos/pkg/data/auth"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/kubernetes"
	domain3 "github.com/athmos-cloud/infra-worker-athmos/pkg/data/project"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/vpc"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/config"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/utils"
)

type Provider struct {
	Metadata            domain.Metadata          `bson:"metadata"`
	KubernetesResources kubernetes.ResourceList  `bson:"kubernetesResources"`
	Identifier          ProviderIdentifier       `bson:"identifier"`
	Type                common.ProviderType      `bson:"type"`
	Auth                domain2.Auth             `bson:"auth"`
	VPCs                map[string]resources.VPC `bson:"vpcs"`
}

type ProviderIdentifier struct {
	ID string `bson:"id"`
}

type ProviderList []Provider

func (provider *Provider) GetMetadata() domain.Metadata {
	return provider.Metadata
}

func (provider *Provider) WithMetadata(request domain.CreateMetadataRequest) {
	provider.Metadata = domain.New(request)
}
func (provider *Provider) GetPluginReference(request dto.GetPluginReferenceRequest) (dto.GetPluginReferenceResponse, errors.Error) {
	switch request.ProviderType {
	case common.GCP:
		return dto.GetPluginReferenceResponse{
			ChartName:    config.Current.Plugins.Crossplane.GCP.Provider.Chart,
			ChartVersion: config.Current.Plugins.Crossplane.GCP.Provider.Version,
		}, errors.Error{}
	}
	return dto.GetPluginReferenceResponse{}, errors.InvalidArgument.WithMessage(fmt.Sprintf("provider type %s not supported", request.ProviderType))
}

func (provider *Provider) FromMap(m map[string]interface{}) errors.Error {
	*provider = Provider{}
	if m["id"] == nil {
		provider.Identifier.ID = utils.GenerateUUID()
	} else {
		provider.Identifier.ID = m["id"].(string)
	}
	if m["name"] == nil {
		return errors.InvalidArgument.WithMessage("name is required")
	}
	provider.Metadata.Name = m["name"].(string)
	return errors.OK
}

func (provider *Provider) InsertIntoProject(project domain3.Project, upsert bool) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (provider *Provider) ToDomain() (interface{}, errors.Error) {
	//TODO implement me
	panic("implement me")
}
