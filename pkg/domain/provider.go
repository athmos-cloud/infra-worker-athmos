package domain

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/common"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource"
)

type Provider struct {
	Name         string              `json:"name"`
	Monitored    bool                `json:"monitored"`
	ProviderType common.ProviderType `json:"provider_type"`
	VPCs         VPCCollection       `json:"vpcs"`
}

func FromProviderDataMapper(provider resource.Provider) Provider {
	return Provider{
		Name:         provider.Identifier.ID,
		Monitored:    provider.Metadata.Monitored,
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
