package project

import (
	"github.com/PaulBarrie/infra-worker/pkg/domain"
)

type GetProjectByIDRequest struct {
	ProjectID string `json:"project_id"`
}

type GetProjectByIDResponse struct {
	Payload domain.Project
}

type GetProjectByOwnerIDRequest struct {
	OwnerID string
}

type GetProjectByOwnerIDResponse struct {
	Payload []domain.Project
}
