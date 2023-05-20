package dto

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
)

type GetNetworkRequest struct {
	IdentifierID identifier.Network `json:"identifierID"`
}

type GetNetworkResponse struct {
	ProjectID string           `json:"projectID"`
	Payload   resource.Network `json:"payload"`
}

type ListNetworksRequest struct {
	ParentID  identifier.ID `json:"parentID"`
	Recursive bool          `json:"recursive"`
}

type ListNetworkResponse struct {
	Payload resource.NetworkCollection `json:"payload"`
}

type CreateNetworkRequest struct {
	ParentID identifier.ID     `json:"parentID"`
	Name     string            `json:"name"`
	Managed  bool              `json:"managed"`
	Tags     map[string]string `json:"tags"`
}

type CreateNetworkResponse struct {
	Payload resource.Network `json:"payload"`
}

type UpdateNetworkRequest struct {
	IdentifierID   identifier.Network `json:"identifierID"`
	Name           string             `json:"name"`
	SecretAuthName string             `json:"secretAuthName"`
}

type DeleteNetworkRequest struct {
	IdentifierID identifier.Network `json:"identifierID"`
	Cascade      bool               `json:"cascade"`
}
