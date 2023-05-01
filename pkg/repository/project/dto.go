package project

import "github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource"

type GetProjectByOwnerIDResponse struct {
	Projects []resource.Project `json:"projects"`
}
