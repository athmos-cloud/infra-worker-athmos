package project

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource"
)

type GetProjectByOwnerIDResponse struct {
	Projects []resource.Project `json:"projects"`
}

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

type UpdateProjectRequest struct {
	ProjectID      string           `json:"projectID"`
	ProjectName    string           `json:"projectName"`
	UpdatedProject resource.Project `json:"project"`
}

type DeleteRequest struct {
	ProjectID string
}
