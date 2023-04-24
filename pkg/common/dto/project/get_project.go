package project

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/project"
)

type GetProjectByIDRequest struct {
	ProjectID string `json:"project_id"`
}

type GetProjectByIDResponse struct {
	Payload domain.Project `json:"payload"`
}

type GetProjectByOwnerIDRequest struct {
	OwnerID string `json:"owner_id"`
}

type GetProjectByOwnerIDResponse struct {
	Payload []domain.Project `json:"payload"`
}
