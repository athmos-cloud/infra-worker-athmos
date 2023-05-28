package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/secret"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
)

type Provider struct {
	IdentifierID   identifier.Provider `json:"identifierID"`
	IdentifierName identifier.Provider `json:"identifierName"`
	Type           types.Provider      `json:"provider"`
	Auth           ProviderAuth        `json:"auth"`
	VPCs           VPCCollection       `json:"vpcs,omitempty"`
	Networks       NetworkCollection   `json:"networks,omitempty"`
}

type ProviderAuth struct {
	Name             string            `json:"name"`
	KubernetesSecret secret.Kubernetes `json:"kubernetesSecret"`
}
type ProviderCollection map[string]Provider

func (provider *Provider) Equals(other Provider) bool {
	return provider.IdentifierID.Equals(&other.IdentifierID) &&
		provider.IdentifierName.Equals(&other.IdentifierName) &&
		provider.Type == other.Type &&
		provider.Auth.Equals(other.Auth) &&
		provider.VPCs.Equals(other.VPCs) &&
		provider.Networks.Equals(other.Networks)
}

func (s *ProviderAuth) Equals(auth ProviderAuth) bool {
	return s.Name == auth.Name &&
		s.KubernetesSecret.Equals(auth.KubernetesSecret)
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
