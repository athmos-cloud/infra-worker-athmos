package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/network"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/secret"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
)

type Provider struct {
	IdentifierID   identifier.Provider       `json:"identifier_id"`
	IdentifierName identifier.Provider       `json:"identifier_name"`
	Type           types.Provider            `json:"provider"`
	Auth           ProviderAuth              `json:"auth"`
	Networks       network.NetworkCollection `json:"networks,omitempty"`
}

type ProviderAuth struct {
	Name             string            `json:"name"`
	KubernetesSecret secret.Kubernetes `json:"kubernetes_secret"`
}
type ProviderCollection map[string]Provider
