package usecase

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository"
	resourceRepos "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository/resource"
)

type Project interface {
	Get(context.Context, *model.Project) errors.Error
	List(context.Context, *[]model.Project) errors.Error
	Create(context.Context, *model.Project) errors.Error
	Update(context.Context, *model.Project) errors.Error
	Delete(context.Context, *model.Project) errors.Error
}

func NewProjectUseCase(repo repository.Project, gcpRepo resourceRepos.Resource, awsRepo resourceRepos.Resource) Project {
	return &projectUseCase{
		projectRepo: repo,
		gcpRepo:     gcpRepo,
		awsRepo:     awsRepo,
	}
}

type projectUseCase struct {
	projectRepo repository.Project
	gcpRepo     resourceRepos.Resource
	awsRepo     resourceRepos.Resource
}

func (pu *projectUseCase) Get(ctx context.Context, project *model.Project) errors.Error {
	foundProject, err := pu.projectRepo.Find(ctx, option.Option{
		Value: repository.FindProjectByIDRequest{
			ID: ctx.Value(context.ProjectIDKey).(string),
		},
	})
	if !err.IsOk() {
		return err
	}
	*project = *foundProject

	return errors.OK
}

func (pu *projectUseCase) List(ctx context.Context, projects *[]model.Project) errors.Error {
	foundProjects, err := pu.projectRepo.FindAll(ctx, option.Option{
		Value: repository.FindAllProjectByOwnerRequest{
			Owner: ctx.Value(context.OwnerIDKey).(string),
		},
	})
	if !err.IsOk() {
		return err
	}
	*projects = *foundProjects

	return errors.OK
}

func (pu *projectUseCase) Create(ctx context.Context, project *model.Project) errors.Error {
	return pu.projectRepo.Create(ctx, project)
}

func (pu *projectUseCase) Update(ctx context.Context, project *model.Project) errors.Error {
	project, err := pu.projectRepo.Find(ctx, option.Option{Value: repository.FindProjectByIDRequest{ID: ctx.Value(context.ProjectIDKey).(string)}})
	if !err.IsOk() {
		return err
	}
	project.Name = ctx.Value(context.RequestKey).(dto.UpdateProjectRequest).Name

	return pu.projectRepo.Update(ctx, project)
}

func (pu *projectUseCase) Delete(ctx context.Context, project *model.Project) errors.Error {
	project, err := pu.projectRepo.Find(ctx, option.Option{Value: repository.FindProjectByIDRequest{ID: ctx.Value(context.ProjectIDKey).(string)}})
	if !err.IsOk() {
		return err
	}
	searchLabels := map[string]string{model.ProjectIDLabelKey: project.ID.Hex()}
	if providers, err := pu.gcpRepo.FindAllProviders(ctx, option.Option{Value: resourceRepos.FindAllResourceOption{Labels: searchLabels}}); !err.IsOk() {
		return err
	} else {
		if len(*providers) > 0 {
			return errors.Conflict.WithMessage("Cannot delete project with providers")
		}
	}
	if providers, err := pu.awsRepo.FindAllProviders(ctx, option.Option{Value: resourceRepos.FindAllResourceOption{Labels: searchLabels}}); !err.IsOk() {
		return err
	} else {
		if len(*providers) > 0 {
			return errors.Conflict.WithMessage("Cannot delete project with providers")
		}
	}

	return pu.projectRepo.Delete(ctx, project)
}
