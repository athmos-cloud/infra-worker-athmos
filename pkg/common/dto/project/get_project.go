package project

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource"
)

type GetProjectByIDRequest struct {
	ProjectID string `json:"project_id"`
}

type GetProjectByIDResponse struct {
	Payload resource.Project `json:"payload"`
}

type GetProjectByOwnerIDRequest struct {
	OwnerID string `json:"owner_id"`
}

type GetProjectByOwnerIDResponse struct {
	Payload []resource.Project `json:"payload"`
}
