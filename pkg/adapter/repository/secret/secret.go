package secret

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/secret"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/share"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository"
	"github.com/kamva/mgm/v3"
	"reflect"
)

type secretRepository struct{}

func NewSecretRepository() repository.Secret {
	return &secretRepository{}
}

type GetSecretByProjectIdAndName struct {
	ProjectId string
	Name      string
}

type GetSecretInProject struct {
	ProjectId string
}

func (s *secretRepository) Find(_ context.Context, opt option.Option) (*secret.Secret, errors.Error) {
	if !opt.SetType(reflect.TypeOf(GetSecretByProjectIdAndName{}).String()).Validate() {
		return nil, errors.InvalidOption.WithMessage(
			fmt.Sprintf("expected GetSecretByProjectIdAndName option, got %s", reflect.TypeOf(opt.Value).String()),
		)
	}
	request := opt.Value.(GetSecretByProjectIdAndName)
	project, err := s.getProjectByID(request.ProjectId)
	if !err.IsOk() {
		return nil, err
	}
	secretAuth, ok := project.Secrets[request.Name]
	if !ok {
		return nil, errors.NotFound.WithMessage(
			fmt.Sprintf("Secret with name %s not found in project %s",
				request.Name, request.ProjectId))
	}
	return &secretAuth, errors.OK
}

func (s *secretRepository) FindAll(_ context.Context, opt option.Option) (*[]secret.Secret, errors.Error) {
	if !opt.SetType(reflect.TypeOf(GetSecretInProject{}).String()).Validate() {
		return nil, errors.InvalidOption.WithMessage(
			fmt.Sprintf("expected GetSecretInProject option, got %s", reflect.TypeOf(opt.Value).String()),
		)
	}
	request := opt.Value.(GetSecretInProject)
	project, err := s.getProjectByID(request.ProjectId)
	if !err.IsOk() {
		return nil, err
	}
	var secrets []secret.Secret
	for _, secretAuth := range project.Secrets {
		secrets = append(secrets, secretAuth)
	}
	return &secrets, errors.OK
}

func (s *secretRepository) Create(ctx context.Context, secretAuth *secret.Secret) errors.Error {
	project, err := s.getProjectByID(ctx.Value(share.ProjectIDKey).(string))
	if !err.IsOk() {
		return err
	}
	if _, ok := project.Secrets[secretAuth.Name]; ok {
		return errors.Conflict.WithMessage(
			fmt.Sprintf("Secret with name %s already exist in project %s",
				secretAuth.Name, project.ID.Hex()))
	}
	project.Secrets[secretAuth.Name] = *secretAuth
	if errUp := mgm.Coll(project).Update(project); errUp != nil {
		return errors.InternalError.WithMessage(errUp.Error())
	}
	return errors.Created
}

func (s *secretRepository) Update(ctx context.Context, secretAuth *secret.Secret) errors.Error {
	project, err := s.getProjectByID(ctx.Value(share.ProjectIDKey).(string))
	if !err.IsOk() {
		return err
	}
	if _, ok := project.Secrets[secretAuth.Name]; !ok {
		return errors.NotFound.WithMessage(
			fmt.Sprintf("Secret with name %s not found in project %s",
				secretAuth.Name, project.ID.Hex()))
	}
	project.Secrets[secretAuth.Name] = *secretAuth
	if errUp := mgm.Coll(project).Update(project); errUp != nil {
		return errors.InternalError.WithMessage(errUp.Error())
	}
	return errors.NoContent
}

func (s *secretRepository) Delete(ctx context.Context, secretAuth *secret.Secret) errors.Error {
	project, err := s.getProjectByID(ctx.Value(share.ProjectIDKey).(string))
	if !err.IsOk() {
		return err
	}
	if _, ok := project.Secrets[secretAuth.Name]; !ok {
		return errors.NotFound.WithMessage(
			fmt.Sprintf("Secret with name %s not found in project %s",
				secretAuth.Name, project.ID.Hex()))
	}
	delete(project.Secrets, secretAuth.Name)
	if errUp := mgm.Coll(project).Update(project); errUp != nil {
		return errors.InternalError.WithMessage(errUp.Error())
	}
	return errors.NoContent
}

func (s *secretRepository) getProjectByID(id string) (*model.Project, errors.Error) {
	project := &model.Project{}
	err := mgm.Coll(project).FindByID(id, project)
	if err != nil {
		return nil, errors.NotFound.WithMessage(fmt.Sprintf("Project with id %s not found", id))
	}
	return project, errors.OK
}
