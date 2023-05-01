package project

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/config"
)

func toGetProjectByOwnerIDResponse(projectList []resource.Project) GetProjectByOwnerIDResponse {
	resp := GetProjectByOwnerIDResponse{}
	for _, project := range projectList {
		resp.Payload = append(resp.Payload, GetProjectByOwnerIDItemResponse{
			ID:   project.ID,
			Name: project.Name,
			URI:  fmt.Sprintf("%s/%s", config.Current.RedirectionURI, project.ID),
		})
	}
	return resp
}