package dto

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/network"
)

type GetSubnetworkResponse struct {
	ProjectID string             `json:"project_id"`
	Payload   network.Subnetwork `json:"payload"`
}

type CreateSubnetworkRequest struct {
	ParentID    identifier.Network `json:"parent_id"`
	Name        string             `json:"name"`
	Managed     bool               `json:"managed" default:"true"`
	Region      string             `json:"region"`
	IPCIDRRange string             `json:"ip_cidr_range" default:"10.0.0.1/28"`
	Tags        map[string]string  `json:"tags"`
}

type CreateSubnetworkResponse struct {
	ProjectID string             `json:"project_id"`
	Payload   network.Subnetwork `json:"payload"`
}

type UpdateSubnetworkRequest struct {
	IdentifierID identifier.Subnetwork `json:"identifier_id"`
	Managed      bool                  `json:"managed"`
	Region       *string               `json:"region"`
	IPCIDRRange  *string               `json:"ip_cidr_range"`
	Tags         *map[string]string    `json:"tags"`
}

type DeleteSubnetworkRequest struct {
	IdentifierID string `json:"identifier_id"`
	Cascade      bool   `json:"cascade" default:"true"`
}
