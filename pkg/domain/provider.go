package domain

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
)

type Provider struct {
	Name         string             `json:"name"`
	Monitored    bool               `json:"monitored"`
	ProviderType types.ProviderType `json:"providerType"`
	VPCs         VPCCollection      `json:"vpcs"`
}

func (provider Provider) ToDataMapper(resourceInput resource.IResource) resource.IResource {
	providerInput := resourceInput.(*resource.Provider)
	providerInput.Identifier.ID = provider.Name
	providerInput.Metadata.Managed = provider.Monitored
	providerInput.Status.PluginReference.ResourceReference.ProviderType = provider.ProviderType
	return providerInput
}

func FromProviderDataMapper(provider *resource.Provider) Provider {
	return Provider{
		Name:         provider.Identifier.ID,
		Monitored:    provider.Metadata.Managed,
		ProviderType: provider.GetPluginReference().ResourceReference.ProviderType,
		VPCs:         FromVPCCollectionDataMapper(provider.VPCs),
	}
}

type ProviderCollection map[string]Provider

func FromProviderCollectionDataMapper(providers resource.ProviderCollection) ProviderCollection {
	result := make(ProviderCollection)
	for _, provider := range providers {
		result[provider.Identifier.ID] = FromProviderDataMapper(&provider)
	}
	return result
}
