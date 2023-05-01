package application

import (
	"fmt"
	dto "github.com/athmos-cloud/infra-worker-athmos/pkg/common/dto/project"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/config"
)

func toGetProjectByOwnerIDResponse(projectList []resource.Project) dto.GetProjectByOwnerIDResponse {
	resp := dto.GetProjectByOwnerIDResponse{}
	for _, project := range projectList {
		resp.Payload = append(resp.Payload, dto.GetProjectByOwnerIDItemResponse{
			ID:   project.ID,
			Name: project.Name,
			URI:  fmt.Sprintf("%s/%s", config.Current.RedirectionURI, project.ID),
		})
	}
	return resp
}
