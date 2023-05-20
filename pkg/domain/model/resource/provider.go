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
