package dto

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
)

type GetNetworkRequest struct {
	ProjectID  string             `json:"projectID"`
	Identifier identifier.Network `json:"identifier"`
}

type GetNetworkResponse struct {
	ProjectID string           `json:"projectID"`
	Payload   resource.Network `json:"payload"`
}

type GetAllNetworksRequest struct {
	ProjectID string        `json:"projectID"`
	ParentID  identifier.ID `json:"parentID"`
	Recursive bool          `json:"recursive"`
}

type GetAllNetworkResponse struct {
	ProjectID string                     `json:"projectID"`
	Payload   resource.NetworkCollection `json:"payload"`
}

type CreateNetworkRequest struct {
	ProjectID      string `json:"projectID"`
	Name           string `json:"name"`
	SecretAuthName string `json:"secretAuthName"`
}

type CreateNetworkResponse struct {
	ProjectID string           `json:"projectID"`
	Payload   resource.Network `json:"payload"`
}

type UpdateNetworkRequest struct {
	ProjectID      string             `json:"projectID"`
	Identifier     identifier.Network `json:"identifier"`
	Name           string             `json:"name"`
	SecretAuthName string             `json:"secretAuthName"`
}

type DeleteNetworkRequest struct {
	ProjectID  string             `json:"projectID"`
	Identifier identifier.Network `json:"identifier"`
	Cascade    bool               `json:"cascade"`
}
