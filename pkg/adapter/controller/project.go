package controller

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/error"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/validator"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model"
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
	if err := validator.ListProjectByOwner(ctx); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	}
	projects := &[]model.Project{}
	if err := pc.projectUseCase.List(ctx, projects); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	} else {
		pc.projectPort.RenderAll(ctx, projects)
	}
}

func (pc *projectController) GetProject(ctx context.Context) {
	if err := validator.GetProject(ctx); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	}
	project := &model.Project{}
	if err := pc.projectUseCase.Get(ctx, project); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	} else {
		pc.projectPort.Render(ctx, project)
	}
}

func (pc *projectController) CreateProject(ctx context.Context) {
	if err := validator.CreateProject(ctx); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	}
	req := ctx.Value(context.RequestKey).(dto.CreateProjectRequest)
	project := model.NewProject(req.Name, req.Owner)
	if err := pc.projectUseCase.Create(ctx, project); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	} else {
		pc.projectPort.RenderCreate(ctx, project)
	}
}

func (pc *projectController) UpdateProject(ctx context.Context) {
	if err := validator.UpdateProject(ctx); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	}
	project := &model.Project{}
	if err := pc.projectUseCase.Update(ctx, project); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	} else {
		pc.projectPort.RenderUpdate(ctx, project)
	}
}

func (pc *projectController) DeleteProject(ctx context.Context) {
	if err := validator.DeleteProject(ctx); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	}
	project := &model.Project{}
	pc.projectUseCase.Delete(ctx, project)
}
