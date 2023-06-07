package repository

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/kubernetes"
	_ "github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/mongo"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/logger"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"
)

const (
	NameDocumentKey    = "name"
	OwnerIDDocumentKey = "owner_id"
)

type projectRepository struct{}

func NewProjectRepository() repository.Project {
	return &projectRepository{}
}

func (p *projectRepository) Find(_ context.Context, opt option.Option) (*model.Project, errors.Error) {
	if !opt.SetType(reflect.TypeOf(repository.FindProjectByIDRequest{}).String()).Validate() {
		return nil, errors.InvalidOption.WithMessage(fmt.Sprintf(
			"expected FindProjectByIDRequest option, got %s",
			reflect.TypeOf(opt.Value).String()))
	}
	request := opt.Value.(repository.FindProjectByIDRequest)
	project := &model.Project{}
	err := mgm.Coll(project).FindByID(request.ID, project)
	if err != nil {
		return nil, errors.NotFound.WithMessage(fmt.Sprintf("Project with id %s not found", request.ID))
	}

	return project, errors.OK
}

func (p *projectRepository) FindAll(_ context.Context, opt option.Option) (*[]model.Project, errors.Error) {
	if !opt.SetType(reflect.TypeOf(repository.FindAllProjectByOwnerRequest{}).String()).Validate() {
		return nil, errors.InvalidOption.WithMessage(
			fmt.Sprintf("expected FindAllProjectByOwnerRequest option, got %s", reflect.TypeOf(opt.Value).String()),
		)
	}
	request := opt.Value.(repository.FindAllProjectByOwnerRequest)
	projects := &[]model.Project{}
	if err := mgm.Coll(&model.Project{}).SimpleFind(projects, bson.M{OwnerIDDocumentKey: request.Owner}); err != nil {
		return nil, errors.InternalError.WithMessage(err.Error())
	}

	return projects, errors.OK
}

func (p *projectRepository) Create(ctx context.Context, project *model.Project) errors.Error {
	projects := &[]model.Project{}
	logger.Info.Printf("Project: %v", project)
	if err := mgm.Coll(&model.Project{}).SimpleFind(projects, bson.M{NameDocumentKey: project.Name, OwnerIDDocumentKey: project.OwnerID}); err != nil {
		return errors.InternalError.WithMessage(err.Error())
	}

	if len(*projects) > 0 {
		return errors.Conflict.WithMessage(fmt.Sprintf("Project with name %s already exists", project.Name))
	}
	if err := mgm.Coll(project).Create(project); err != nil {
		logger.Info.Printf("err: %v", err)
		return errors.InternalError.WithMessage(err.Error())
	}
	if err := kubernetes.Client().Client.Create(ctx, &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: project.Namespace,
		},
	}); err != nil {
		logger.Info.Printf("err: %v", err)
		return errors.KubernetesError.WithMessage(err.Error())
	}
	return errors.Created
}

func (p *projectRepository) Update(_ context.Context, project *model.Project) errors.Error {
	if err := mgm.Coll(project).Update(project); err != nil {
		return errors.InternalError.WithMessage(err.Error())
	}
	return errors.NoContent
}

func (p *projectRepository) Delete(ctx context.Context, project *model.Project) errors.Error {
	if err := mgm.Coll(project).Delete(project); err != nil {
		return errors.InternalError.WithMessage(err.Error())
	}
	if err := kubernetes.Client().Client.Delete(ctx, &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: project.Namespace,
		},
	}); err != nil {
		return errors.KubernetesError.WithMessage(err.Error())
	}
	return errors.NoContent
}
