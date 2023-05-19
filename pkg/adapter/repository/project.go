package repository

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/kubernetes"
	_ "github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/mongo"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
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

type FindByIDRequest struct {
	ID string
}

type FindAllByOwnerRequest struct {
	Owner string
}

func (p *projectRepository) Find(_ context.Context, opt option.Option) (*model.Project, errors.Error) {
	if !opt.SetType(reflect.TypeOf(FindByIDRequest{}).String()).Validate() {
		return nil, errors.InvalidOption.WithMessage(fmt.Sprintf(
			"expected FindByIDRequest option, got %s",
			reflect.TypeOf(opt.Value).String()))
	}
	request := opt.Value.(FindByIDRequest)
	project := &model.Project{}
	err := mgm.Coll(project).FindByID(request.ID, project)
	if err != nil {
		return nil, errors.NotFound.WithMessage(fmt.Sprintf("Project with id %s not found", request.ID))
	}

	return project, errors.OK
}

func (p *projectRepository) FindAll(_ context.Context, opt option.Option) (*[]model.Project, errors.Error) {
	if !opt.SetType(reflect.TypeOf(FindAllByOwnerRequest{}).String()).Validate() {
		return nil, errors.InvalidOption.WithMessage(
			fmt.Sprintf("expected FindAllByOwnerRequest option, got %s", reflect.TypeOf(opt.Value).String()),
		)
	}
	request := opt.Value.(FindAllByOwnerRequest)
	projects := &[]model.Project{}
	if err := mgm.Coll(&model.Project{}).SimpleFind(projects, bson.M{OwnerIDDocumentKey: request.Owner}); err != nil {
		return nil, errors.InternalError.WithMessage(err.Error())
	}

	return projects, errors.OK
}

func (p *projectRepository) Create(ctx context.Context, project *model.Project) errors.Error {
	projects := &[]model.Project{}
	if err := mgm.Coll(&model.Project{}).SimpleFind(projects, bson.M{NameDocumentKey: project.Name, OwnerIDDocumentKey: project.OwnerID}); err != nil {
		return errors.InternalError.WithMessage(err.Error())
	}
	if len(*projects) > 0 {
		return errors.Conflict.WithMessage(fmt.Sprintf("Project with name %s already exists", project.Name))
	}
	if err := mgm.Coll(project).Create(project); err != nil {
		return errors.InternalError.WithMessage(err.Error())
	}
	if err := kubernetes.Client().K8sClient.Create(ctx, &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: project.Namespace,
		},
	}); err != nil {
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
	if err := kubernetes.Client().K8sClient.Delete(ctx, &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: project.Namespace,
		},
	}); err != nil {
		return errors.KubernetesError.WithMessage(err.Error())
	}
	return errors.NoContent
}
