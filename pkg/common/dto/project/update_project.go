package project

import "github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource"

type UpdateProjectRequest struct {
	ProjectID      string           `json:"projectID"`
	ProjectName    string           `json:"projectName"`
	UpdatedProject resource.Project `json:"project"`
}
