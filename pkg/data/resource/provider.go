package resource

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/common"
	dto "github.com/athmos-cloud/infra-worker-athmos/pkg/common/dto/resource"
	auth "github.com/athmos-cloud/infra-worker-athmos/pkg/data/auth"
	resourcePlugin "github.com/athmos-cloud/infra-worker-athmos/pkg/data/plugin"

	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/config"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"reflect"
)

type Provider struct {
	Metadata            metadata.Metadata       `bson:"metadata"`
	Identifier          identifier.Provider     `bson:"identifier"`
	VPC                 string                  `bson:"vpc" plugin:"vpc"`
	KubernetesResources kubernetes.ResourceList `bson:"kubernetesResources"`
	Type                common.ProviderType     `bson:"type" plugin:"type"`
	Auth                auth.Auth               `bson:"auth" plugin:"auth"`
	VPCs                VPCCollection           `bson:"vpcs"`
}

func (provider *Provider) Equals(other Provider) bool {
	return provider.Metadata.Equals(other.Metadata) &&
		provider.Identifier.Equals(other.Identifier) &&
		provider.VPC == other.VPC &&
		provider.KubernetesResources.Equals(other.KubernetesResources) &&
		provider.Type == other.Type &&
		provider.Auth.Equals(other.Auth) &&
		provider.VPCs.Equals(other.VPCs)
}

type ProviderCollection map[string]Provider

func NewProvider(id identifier.Provider) Provider {
	return Provider{
		Metadata: metadata.New(metadata.CreateMetadataRequest{
			Name: id.ID,
		}),
		Identifier: id,
		VPCs:       make(VPCCollection),
	}
}

func (provider *Provider) New(id identifier.ID) (IResource, errors.Error) {
	if reflect.TypeOf(id) != reflect.TypeOf(identifier.Provider{}) {
		return nil, errors.InvalidArgument.WithMessage("invalid id type")
	}
	res := NewProvider(id.(identifier.Provider))
	return &res, errors.OK
}

func (provider *Provider) GetMetadata() metadata.Metadata {
	return provider.Metadata
}

func (provider *Provider) WithMetadata(request metadata.CreateMetadataRequest) {
	provider.Metadata = metadata.New(request)
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
	return resourcePlugin.InjectMapIntoStruct(m, provider)
}

func (provider *Provider) Insert(project Project, update ...bool) errors.Error {
	shouldUpdate := false
	if len(update) > 0 {
		shouldUpdate = update[0]
	}
	id := provider.Identifier
	_, ok := project.Resources[provider.Identifier.ID]
	if !ok && shouldUpdate {
		return errors.NotFound.WithMessage(fmt.Sprintf("provider %s not found", id.ID))
	} else if ok && !shouldUpdate {
		return errors.Conflict.WithMessage(fmt.Sprintf("provider %s already exists", id.ID))
	}
	project.Resources[provider.Identifier.ID] = *provider
	return errors.OK
}

func (provider *Provider) Remove(project Project) errors.Error {
	//TODO implement me
	panic("implement me")
}

func validateEntryMap(entry map[string]interface{}) errors.Error {
	if entry["name"] == "" || reflect.TypeOf(entry["name"]).Kind() != reflect.String {
		return errors.InvalidArgument.WithMessage("a string name is required")
	}
	if entry["value"] == "" {
		return errors.InvalidArgument.WithMessage("value is required")
	}
	return errors.OK
}
