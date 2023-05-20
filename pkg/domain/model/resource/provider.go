package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/secret"
)

type Provider struct {
	Metadata       metadata.Metadata   `json:"metadata"`
	IdentifierID   identifier.Provider `json:"identifierID"`
	IdentifierName identifier.Provider `json:"identifierName"`
	Auth           secret.Secret       `json:"auth"`
	VPCs           VPCCollection       `json:"vpcs,omitempty"`
	Networks       NetworkCollection   `json:"networks,omitempty"`
}

type ProviderCollection map[string]Provider

func (provider *Provider) Equals(other Provider) bool {
	return provider.Metadata.Equals(other.Metadata) &&
		provider.IdentifierID.Equals(&other.IdentifierID) &&
		provider.IdentifierName.Equals(&other.IdentifierName) &&
		provider.Auth.Equals(other.Auth) &&
		provider.VPCs.Equals(other.VPCs) &&
		provider.Networks.Equals(other.Networks)
}

func (collection *ProviderCollection) Equals(other ProviderCollection) bool {
	if len(*collection) != len(other) {
		return false
	}
	for key, value := range *collection {
		if !value.Equals(other[key]) {
			return false
		}
	}
	return true
}
