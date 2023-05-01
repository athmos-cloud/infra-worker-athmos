package project

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain"
)

type CreateProjectRequest struct {
	ProjectName string `json:"projectName"`
	OwnerID     string `json:"ownerID"`
}

type CreateProjectResponse struct {
	ProjectID string `json:"projectID"`
}

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

type UpdateProjectRequest struct {
	ProjectID      string         `json:"projectID"`
	ProjectName    string         `json:"projectName"`
	UpdatedProject domain.Project `json:"project"`
}

type DeleteRequest struct {
	ProjectID string
}
