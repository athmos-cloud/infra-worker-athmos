package project

import "github.com/PaulBarrie/infra-worker/pkg/project"

type GetProjectByIDRequest struct {
	ProjectID string `json:"project_id"`
}

type GetProjectByIDResponse struct {
	Payload project.Project
}

type GetProjectByOwnerIDRequest struct {
	OwnerID string
}

type GetProjectByOwnerIDResponse struct {
	Payload []project.Project
}
