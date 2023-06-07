package presenter

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/config"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/logger"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/output"
	"github.com/gin-gonic/gin"
)

type Project struct{}

func NewProjectPresenter() output.ProjectPort {
	return &Project{}
}

func (p *Project) Render(ctx context.Context, project *model.Project) {
	resp := dto.GetProjectResponse{
		ID:    project.ID.Hex(),
		Name:  project.Name,
		Owner: project.OwnerID,
	}
	ctx.JSON(200, gin.H{"payload": resp})
}

func (p *Project) RenderCreate(ctx context.Context, project *model.Project) {
	resp := dto.CreateProjectResponse{
		ID:             project.ID.Hex(),
		RedirectionURL: fmt.Sprintf("%s/projects/%s", config.Current.RedirectionURL, project.ID.Hex()),
	}
	ctx.JSON(201, gin.H{"payload": resp})
}

func (p *Project) RenderUpdate(ctx context.Context, project *model.Project) {
	ctx.JSON(204, gin.H{"message": fmt.Sprintf("Project %s updated", project.ID.Hex())})
}

func (p *Project) RenderAll(ctx context.Context, projects *[]model.Project) {
	resp := make([]dto.ListProjectResponseItem, len(*projects))
	for i, project := range *projects {
		resp[i] = dto.ListProjectResponseItem{
			ID:             project.ID.Hex(),
			Name:           project.Name,
			Owner:          project.OwnerID,
			RedirectionURL: fmt.Sprintf("%s/projects/%s", config.Current.RedirectionURL, project.ID.Hex()),
		}
	}
	logger.Info.Printf("resp: %v", resp)
	ctx.JSON(200, gin.H{"payload": resp})
}
