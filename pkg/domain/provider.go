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

func FromProviderDataMapper(provider resource.Provider) Provider {
	return Provider{
		Name:         provider.Identifier.ID,
		Monitored:    provider.Metadata.Managed,
		ProviderType: provider.Type,
		VPCs:         FromVPCCollectionDataMapper(provider.VPCs),
	}
}

type ProviderCollection map[string]Provider

func FromProviderCollectionDataMapper(providers resource.ProviderCollection) ProviderCollection {
	result := make(ProviderCollection)
	for _, provider := range providers {
		result[provider.Identifier.ID] = FromProviderDataMapper(provider)
	}
	return result
}
