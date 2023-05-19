package usecase

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	arepo "github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/share"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository"
)

type Project interface {
	Get(context.Context, *model.Project) errors.Error
	List(context.Context, *[]model.Project) errors.Error
	Create(context.Context, *model.Project) errors.Error
	Update(context.Context, *model.Project) errors.Error
	Delete(context.Context, *model.Project) errors.Error
}

func NewProjectUseCase(repo repository.Project) Project {
	return &projectUseCase{repo: repo}
}

type projectUseCase struct {
	repo repository.Project
}

func (pu *projectUseCase) Get(ctx context.Context, project *model.Project) errors.Error {
	foundProject, err := pu.repo.Find(ctx, option.Option{
		Value: arepo.FindByIDRequest{
			ID: ctx.Value(share.ProjectIDKey).(string),
		},
	})
	if !err.IsOk() {
		return err
	}
	*project = *foundProject

	return errors.OK
}

func (pu *projectUseCase) List(ctx context.Context, projects *[]model.Project) errors.Error {
	foundProjects, err := pu.repo.FindAll(ctx, option.Option{
		Value: arepo.FindAllByOwnerRequest{
			Owner: ctx.Value(share.OwnerIDKey).(string),
		},
	})
	if !err.IsOk() {
		return err
	}
	*projects = *foundProjects

	return errors.OK
}

func (pu *projectUseCase) Create(ctx context.Context, project *model.Project) errors.Error {
	return pu.repo.Create(ctx, project)
}

func (pu *projectUseCase) Update(ctx context.Context, project *model.Project) errors.Error {
	project, err := pu.repo.Find(ctx, option.Option{Value: arepo.FindByIDRequest{ID: ctx.Value(share.ProjectIDKey).(string)}})
	if !err.IsOk() {
		return err
	}
	project.Name = ctx.Value(share.RequestContextKey).(dto.UpdateProjectRequest).Name

	return pu.repo.Update(ctx, project)
}

func (pu *projectUseCase) Delete(ctx context.Context, project *model.Project) errors.Error {
	project, err := pu.repo.Find(ctx, option.Option{Value: arepo.FindByIDRequest{ID: ctx.Value(share.ProjectIDKey).(string)}})
	if !err.IsOk() {
		return err
	}

	return pu.repo.Delete(ctx, project)
}
