package resource

import (
	"fmt"
	auth "github.com/athmos-cloud/infra-worker-athmos/pkg/data/auth"
	resourcePlugin "github.com/athmos-cloud/infra-worker-athmos/pkg/data/plugin"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/status"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/config"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/utils"
	"github.com/kamva/mgm/v3"
	"reflect"
)

type ProviderCollection map[string]Provider

func NewProvider(payload NewResourcePayload) Provider {
	payload.Validate()
	id := identifier.Provider{
		ProviderID: fmt.Sprintf("%s-%s", payload.Name, utils.RandomString(resourceIDSuffixLength)),
	}
	return Provider{
		Metadata: metadata.New(metadata.CreateMetadataRequest{
			Name:         id.ProviderID,
			NotMonitored: !payload.Monitored,
			Tags:         payload.Tags,
		}),
		Identifier: id,
		Status:     status.New(id.ProviderID, types.Provider, payload.Provider),
		VPCs:       make(VPCCollection),
	}
}

type Provider struct {
	mgm.DefaultModel `bson:",inline"`
	Metadata         metadata.Metadata     `bson:"metadata"`
	Identifier       identifier.Provider   `bson:"identifier"`
	Status           status.ResourceStatus `bson:"status"`
	Auth             auth.Auth             `bson:"auth" plugin:"auth"`
	VPCs             VPCCollection         `bson:"vpcs"`
	Networks         NetworkCollection     `bson:"networks"`
}

func (provider *Provider) GetIdentifier() identifier.ID {
	return provider.Identifier
}

func (provider *Provider) Equals(other Provider) bool {
	return provider.Metadata.Equals(other.Metadata) &&
		provider.Identifier.Equals(other.Identifier) &&
		provider.Status.Equals(other.Status) &&
		provider.Auth.Equals(other.Auth) &&
		provider.VPCs.Equals(other.VPCs) &&
		provider.Networks.Equals(other.Networks)
}

func (provider *Provider) New(payload NewResourcePayload) IResource {
	res := NewProvider(payload)
	return &res
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

func (provider *Provider) GetPluginReference() resourcePlugin.Reference {
	if !provider.Status.PluginReference.ChartReference.Empty() {
		return provider.Status.PluginReference
	}
	switch provider.Status.PluginReference.ResourceReference.ProviderType {
	case types.GCP:
		provider.Status.PluginReference.ChartReference = resourcePlugin.HelmChartReference{
			ChartName:    config.Current.Plugins.Crossplane.GCP.Provider.Chart,
			ChartVersion: config.Current.Plugins.Crossplane.GCP.Provider.Version,
		}
		return provider.Status.PluginReference
	}
	panic(errors.InvalidArgument.WithMessage(fmt.Sprintf("provider type %s not supported", provider.Status.PluginReference.ResourceReference.ProviderType)))
}

func (provider *Provider) FromMap(m map[string]interface{}) {
	if err := resourcePlugin.InjectMapIntoStruct(m, provider); !err.IsOk() {
		panic(err)
	}
}

func (provider *Provider) Insert(resource IResource, update ...bool) {
	shouldUpdate := false
	if len(update) > 0 {
		shouldUpdate = update[0]
	}
	idPayload := identifier.IDToPayload(resource.GetIdentifier())

	if idPayload.VPCID != "" {
		provider.insertInVPC(resource, shouldUpdate)
		return
	} else if idPayload.NetworkID != "" {
		provider.insertInNetwork(resource, shouldUpdate)
		return
	}

	panic(errors.InvalidArgument.WithMessage("ID type must contain either VPC or Network"))
}

func (provider *Provider) insertInVPC(resource IResource, update ...bool) {
	shouldUpdate := false
	if len(update) > 0 {
		shouldUpdate = update[0]
	}
	idPayload := identifier.IDToPayload(resource.GetIdentifier())
	_, ok := provider.VPCs[idPayload.VPCID]
	if reflect.TypeOf(resource) == reflect.TypeOf(&VPC{}) && !ok && shouldUpdate {
		panic(errors.NotFound.WithMessage(fmt.Sprintf("vpc %s not found in provider %s", idPayload.VPCID, provider.Identifier.ProviderID)))
	}
	if reflect.TypeOf(resource) == reflect.TypeOf(&VPC{}) && ok && !shouldUpdate {
		panic(errors.Conflict.WithMessage(fmt.Sprintf("vpc %s already exists in provider %s", idPayload.ProviderID, provider.Identifier.ProviderID)))
	}
	if reflect.TypeOf(resource) == reflect.TypeOf(&VPC{}) && (ok && shouldUpdate || !ok && !shouldUpdate) {
		vpc := resource.(*VPC)
		provider.VPCs[idPayload.VPCID] = *vpc
		return
	} else if reflect.TypeOf(resource) != reflect.TypeOf(&VPC{}) && idPayload.VPCID != "" {
		vpc := provider.VPCs[idPayload.VPCID]
		vpc.Insert(resource, update...)
		return
	}
	panic(errors.InternalError.WithMessage(fmt.Sprintf("Invalid vpc insertion %v", resource)))
}

func (provider *Provider) insertInNetwork(resource IResource, update ...bool) {
	shouldUpdate := false
	if len(update) > 0 {
		shouldUpdate = update[0]
	}
	idPayload := identifier.IDToPayload(resource.GetIdentifier())
	_, ok := provider.Networks[idPayload.NetworkID]
	if reflect.TypeOf(resource) == reflect.TypeOf(&Network{}) && !ok && shouldUpdate {
		panic(errors.NotFound.WithMessage(fmt.Sprintf("network %s not found in provider %s", idPayload.NetworkID, provider.Identifier.ProviderID)))
	}
	if reflect.TypeOf(resource) == reflect.TypeOf(&Network{}) && ok && !shouldUpdate {
		panic(errors.Conflict.WithMessage(fmt.Sprintf("network %s already exists in provider %s", idPayload.ProviderID, provider.Identifier.ProviderID)))
	}
	if reflect.TypeOf(resource) == reflect.TypeOf(&Network{}) && (ok && shouldUpdate || !ok && !shouldUpdate) {
		network := resource.(*Network)
		provider.Networks[idPayload.NetworkID] = *network
	} else if reflect.TypeOf(resource) != reflect.TypeOf(&Network{}) && idPayload.NetworkID != "" {
		network := provider.Networks[idPayload.NetworkID]
		network.Insert(resource, update...)
	}
	panic(errors.InternalError.WithMessage(fmt.Sprintf("Invalid vpc insertion %v", resource)))
}

func (provider *Provider) Remove(resource IResource) {
	idPayload := identifier.IDToPayload(resource.GetIdentifier())
	if idPayload.VPCID != "" && reflect.TypeOf(resource) == reflect.TypeOf(&VPC{}) {
		delete(provider.VPCs, idPayload.VPCID)
	} else if idPayload.VPCID != "" && reflect.TypeOf(resource) != reflect.TypeOf(&Network{}) {
		vpc := provider.VPCs[idPayload.VPCID]
		vpc.Remove(resource)
	} else if idPayload.NetworkID != "" && reflect.TypeOf(resource) == reflect.TypeOf(&Network{}) {
		delete(provider.Networks, idPayload.NetworkID)
	} else if idPayload.NetworkID != "" && reflect.TypeOf(resource) != reflect.TypeOf(&Network{}) {
		network := provider.Networks[idPayload.NetworkID]
		network.Remove(resource)
	}
}
