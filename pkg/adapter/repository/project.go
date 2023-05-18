package repository

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/kubernetes"
	_ "github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/mongo"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/logger"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
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

type Project struct{}

type FindByIDRequest struct {
	ID string
}

type FindAllByOwnerRequest struct {
	Owner string
}

func (p *Project) Find(_ context.Context, opt option.Option) *model.Project {
	if !opt.SetType(reflect.TypeOf(FindByIDRequest{}).String()).Validate() {
		panic(errors.InvalidOption.WithMessage(
			fmt.Sprintf(
				"expected FindByIDRequest option, got %s",
				reflect.TypeOf(opt.Value).String(),
			),
		))
	}
	request := opt.Value.(FindByIDRequest)
	project := &model.Project{}
	err := mgm.Coll(project).FindByID(request.ID, project)
	if err != nil {
		panic(errors.NotFound.WithMessage(fmt.Sprintf("Project with id %s not found", request.ID)))
	}

	return project
}

func (p *Project) FindAll(_ context.Context, opt option.Option) []*model.Project {
	if !opt.SetType(reflect.TypeOf(FindAllByOwnerRequest{}).String()).Validate() {
		panic(errors.InvalidOption.WithMessage(
			fmt.Sprintf(
				"expected FindAllByOwnerRequest option, got %s",
				reflect.TypeOf(opt.Value).String(),
			),
		))
	}
	request := opt.Value.(FindAllByOwnerRequest)
	var projects []model.Project
	if err := mgm.Coll(&model.Project{}).SimpleFind(&projects, bson.M{OwnerIDDocumentKey: request.Owner}); err != nil {
		panic(errors.InternalError.WithMessage(err.Error()))
	}
	foundProjects := make([]*model.Project, len(projects))
	for i, project := range projects {
		foundProjects[i] = &project
	}
	return foundProjects
}

func (p *Project) Create(ctx context.Context, projectCh chan *model.Project, errCh chan errors.Error) {
	logger.Info.Println("Creating project repo")
	project := <-projectCh
	logger.Info.Println("Creating project repo1")
	projects := &[]model.Project{}
	if err := mgm.Coll(&model.Project{}).SimpleFind(projects, bson.M{NameDocumentKey: project.Name, OwnerIDDocumentKey: project.OwnerID}); err != nil {
		errCh <- errors.InternalError.WithMessage(err.Error())
		return
	}
	if len(*projects) > 0 {
		logger.Info.Println("Project already exists")
		errCh <- errors.Conflict.WithMessage(fmt.Sprintf("Project with name %s already exists", project.Name))
		return
	}
	if err := mgm.Coll(project).Create(project); err != nil {
		errCh <- errors.InternalError.WithMessage(err.Error())
		return
	}
	if err := kubernetes.Client.K8sClient.Create(ctx, &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: project.Namespace,
		},
	}); err != nil {
		errCh <- errors.KubernetesError.WithMessage(err.Error())
		return
	}
	projectCh <- project
	return
}

func (p *Project) Update(_ context.Context, project *model.Project) *model.Project {
	if err := mgm.Coll(project).Update(project); err != nil {
		panic(errors.InternalError.WithMessage(err.Error()))
	}
	return project
}

func (p *Project) Delete(ctx context.Context, project *model.Project) {
	if err := mgm.Coll(project).Delete(project); err != nil {
		panic(errors.InternalError.WithMessage(err.Error()))
	}
	if err := kubernetes.Client.K8sClient.Delete(ctx, &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: project.Namespace,
		},
	}); err != nil {
		panic(errors.KubernetesError.WithMessage(err.Error()))
	}
}
