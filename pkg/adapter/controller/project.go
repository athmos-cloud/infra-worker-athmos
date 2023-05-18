package controller

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/validator"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/output"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase"
)

type Project interface {
	ListProjectByOwner(context.Context)
	GetProject(context.Context)
	CreateProject(context.Context)
	UpdateProject(context.Context)
	DeleteProject(context.Context)
}

type projectController struct {
	projectUseCase usecase.Project
	projectPort    output.ProjectPort
}

func NewProjectController(projectUseCase usecase.Project, projectPort output.ProjectPort) Project {
	return &projectController{projectUseCase: projectUseCase, projectPort: projectPort}
}

func (pc *projectController) ListProjectByOwner(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			raiseError(ctx, r)
		}
	}()
	validator.ListProjectByOwner(ctx)
	var projects []*model.Project
	projects = pc.projectUseCase.List(ctx, projects)
	pc.projectPort.RenderAll(ctx, projects)
}

func (pc *projectController) GetProject(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			raiseError(ctx, r)
		}
	}()
	validator.GetProject(ctx)
	var project *model.Project
	project = pc.projectUseCase.Get(ctx, project)
	pc.projectPort.Render(ctx, project)
}

func (pc *projectController) CreateProject(ctx context.Context) {
	var req dto.CreateProjectRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(400, errors.BadRequest.WithMessage(fmt.Sprintf("Expected : %+v", dto.CreateProjectRequest{})))
	}
	init := false
	projectCh := make(chan *model.Project, 1)
	errCh := make(chan errors.Error)

	go func() {
		projectCh <- model.NewProject(req.Name, req.Owner)
		pc.projectUseCase.Create(ctx, projectCh, errCh)
	}()
	for {
		select {
		case project := <-projectCh:
			if !init {
				projectCh <- project
				init = true
				continue
			}
			pc.projectPort.RenderCreate(ctx, project)
			close(errCh)
			return
		case e := <-errCh:
			raiseError(ctx, e)
			close(errCh)
			close(projectCh)
			return
		}
	}
}

func (pc *projectController) UpdateProject(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			raiseError(ctx, r)
		}
	}()
	validator.UpdateProject(ctx)
	var project *model.Project
	pc.projectUseCase.Update(ctx, project)
}

func (pc *projectController) DeleteProject(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			raiseError(ctx, r)
		}
	}()
	validator.DeleteProject(ctx)
	var project *model.Project
	pc.projectUseCase.Delete(ctx, project)
}
