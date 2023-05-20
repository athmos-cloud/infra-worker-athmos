package dto

import "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"

type GetSubnetworkRequest struct {
	IdentifierID identifier.Subnetwork `json:"identifierID"`
}

type GetSubnetworkResponse struct {
	ProjectID string                `json:"projectID"`
	Payload   identifier.Subnetwork `json:"payload"`
}

type ListSubnetworksRequest struct {
	ParentID  identifier.Network `json:"parentID"`
	Recursive bool               `json:"recursive"`
}

type ListSubnetworksResponse struct {
	ProjectID string                `json:"projectID"`
	Payload   identifier.Subnetwork `json:"payload"`
}

type CreateSubnetworkRequest struct {
	ParentID identifier.Network `json:"parentID"`
	Name     string             `json:"name"`
}

type CreateSubnetworkResponse struct {
	ProjectID string                `json:"projectID"`
	Payload   identifier.Subnetwork `json:"payload"`
}

type UpdateSubnetworkRequest struct {
	IdentifierID identifier.Subnetwork `json:"identifierID"`
	Name         string                `json:"name"`
}

type DeleteSubnetworkRequest struct {
	IdentifierID identifier.Subnetwork `json:"identifierID"`
}
