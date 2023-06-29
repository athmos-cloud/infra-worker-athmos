package dto

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/network"
)

type GetNetworkResponse struct {
	ProjectID string          `json:"project_id"`
	Payload   network.Network `json:"payload"`
}

type CreateNetworkRequest struct {
	ParentIDProvider identifier.Provider `json:"parent_id"`
	Name             string              `json:"name"`
	Region           string              `json:"region,omitempty"`
	Managed          bool                `json:"managed" default:"true"`
	Tags             map[string]string   `json:"tags"`
}

type CreateNetworkResponse struct {
	ProjectID string          `json:"project_id"`
	Payload   network.Network `json:"payload"`
}

type UpdateNetworkRequest struct {
	IdentifierID identifier.Network `json:"identifier_id"`
	Tags         *map[string]string `json:"tags"`
	Managed      *bool              `json:"managed"`
}

type DeleteNetworkRequest struct {
	IdentifierID string `json:"identifier_id"`
	Cascade      bool   `json:"cascade" default:"true"`
}
