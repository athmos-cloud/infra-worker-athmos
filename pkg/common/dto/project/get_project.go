package project

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource"
)

type GetProjectByIDRequest struct {
	ProjectID string `json:"projectID"`
}

type GetProjectByIDResponse struct {
	Payload resource.Project `json:"payload"`
}

type GetProjectByOwnerIDRequest struct {
	OwnerID string `json:"ownerId"`
}

type GetProjectByOwnerIDResponse struct {
	Payload []GetProjectByOwnerIDItemResponse `json:"payload"`
}

type GetProjectByOwnerIDItemResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URI  string `json:"uri"`
}
