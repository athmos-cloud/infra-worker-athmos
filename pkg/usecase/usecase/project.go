package usecase

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/context"
	arepo "github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository"
)

const (
	ProjectIDKeyContext    = "project_id"
	ProjectOwnerKeyContext = "project_owner"
)

type Project interface {
	Get(context.Context, *model.Project) *model.Project
	List(context.Context, []*model.Project) []*model.Project
	Create(context.Context, chan *model.Project, chan errors.Error)
	Update(context.Context, *model.Project) *model.Project
	Delete(context.Context, *model.Project)
}

func NewProjectUseCase(repo repository.Project) Project {
	return &projectUseCase{repo: repo}
}

type projectUseCase struct {
	repo repository.Project
}

func (pu *projectUseCase) Get(ctx context.Context, project *model.Project) *model.Project {
	return pu.repo.Find(ctx, option.Option{
		Value: arepo.FindByIDRequest{
			ID: ctx.Value(ProjectIDKeyContext).(string),
		},
	})
}

func (pu *projectUseCase) List(ctx context.Context, project []*model.Project) []*model.Project {
	project = pu.repo.FindAll(ctx, option.Option{
		Value: arepo.FindAllByOwnerRequest{
			Owner: ctx.Value(ProjectOwnerKeyContext).(string),
		},
	})
	return project
}

func (pu *projectUseCase) Create(ctx context.Context, projectCh chan *model.Project, errCh chan errors.Error) {
	pu.repo.Create(ctx, projectCh, errCh)
}

func (pu *projectUseCase) Update(ctx context.Context, project *model.Project) *model.Project {
	return pu.repo.Update(ctx, project)
}

func (pu *projectUseCase) Delete(ctx context.Context, project *model.Project) {
	pu.repo.Delete(ctx, project)
}
