package resource

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/common"
	auth "github.com/athmos-cloud/infra-worker-athmos/pkg/data/auth"
	resourcePlugin "github.com/athmos-cloud/infra-worker-athmos/pkg/data/plugin"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/status"

	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/config"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"reflect"
)

type Provider struct {
	Metadata   metadata.Metadata     `bson:"metadata"`
	Identifier identifier.Provider   `bson:"identifier"`
	Status     status.ResourceStatus `bson:"status"`
	VPC        string                `bson:"vpc" plugin:"vpc"`
	Type       common.ProviderType   `bson:"type" plugin:"type"`
	Auth       auth.Auth             `bson:"auth" plugin:"auth"`
	VPCs       VPCCollection         `bson:"vpcs"`
}

func (provider *Provider) Equals(other Provider) bool {
	return provider.Metadata.Equals(other.Metadata) &&
		provider.Identifier.Equals(other.Identifier) &&
		provider.Status.Equals(other.Status) &&
		provider.VPC == other.VPC &&
		provider.Type == other.Type &&
		provider.Auth.Equals(other.Auth) &&
		provider.VPCs.Equals(other.VPCs)
}

type ProviderCollection map[string]Provider

func NewProvider(id identifier.Provider, providerType common.ProviderType) Provider {
	return Provider{
		Metadata: metadata.New(metadata.CreateMetadataRequest{
			Name: id.ID,
		}),
		Identifier: id,
		Status:     status.New(id.ID, common.Provider, providerType),
		VPCs:       make(VPCCollection),
	}
}

func (provider *Provider) New(id identifier.ID, providerType common.ProviderType) (IResource, errors.Error) {
	if reflect.TypeOf(id) != reflect.TypeOf(identifier.Provider{}) {
		return nil, errors.InvalidArgument.WithMessage("invalid id type")
	}
	res := NewProvider(id.(identifier.Provider), providerType)
	return &res, errors.OK
}

func (provider *Provider) GetMetadata() metadata.Metadata {
	return provider.Metadata
}

func (provider *Provider) SetStatus(resourceStatus status.ResourceStatus) {
	provider.Status = resourceStatus
}

func (provider *Provider) GetStatus() status.ResourceStatus {
	return provider.Status
}

func (provider *Provider) SetMetadata(request metadata.CreateMetadataRequest) {
	provider.Metadata = metadata.New(request)
}

func (provider *Provider) GetPluginReference() (resourcePlugin.Reference, errors.Error) {
	if !provider.Status.PluginReference.ChartReference.Empty() {
		return provider.Status.PluginReference, errors.OK
	}
	switch provider.Status.PluginReference.ResourceReference.ProviderType {
	case common.GCP:
		provider.Status.PluginReference.ChartReference = resourcePlugin.HelmChartReference{
			ChartName:    config.Current.Plugins.Crossplane.GCP.Subnet.Chart,
			ChartVersion: config.Current.Plugins.Crossplane.GCP.Subnet.Version,
		}
		return provider.Status.PluginReference, errors.OK
	}
	return resourcePlugin.Reference{}, errors.InvalidArgument.WithMessage(fmt.Sprintf("provider type %s not supported", provider.Status.PluginReference.ResourceReference.ProviderType))
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
	id := provider.Identifier
	_, ok := project.Resources[provider.Identifier.ID]
	if !ok {
		return errors.NotFound.WithMessage(fmt.Sprintf("provider %s not found", id.ID))
	}
	delete(project.Resources, provider.Identifier.ID)
	return errors.NoContent
}
