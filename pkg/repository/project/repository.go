package project

import (
	"context"
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/dao/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/dao/mongo"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/logger"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
	"sync"
)

const (
	NameDocumentKey    = "name"
	OwnerIDDocumentKey = "owner_id"
)

var ProjectRepository *Repository
var lock = &sync.Mutex{}

type Repository struct {
	MongoDAO      *mongo.DAO
	kubernetesDAO *kubernetes.DAO
}

func init() {
	lock.Lock()
	defer lock.Unlock()
	if ProjectRepository == nil {
		logger.Info.Printf("Init project repository...")
		ProjectRepository = &Repository{
			MongoDAO:      mongo.Client,
			kubernetesDAO: kubernetes.Client,
		}
	}
}

func (repository *Repository) Create(ctx context.Context, opt option.Option) interface{} {
	if !opt.SetType(reflect.TypeOf(CreateProjectRequest{}).String()).Validate() {
		panic(errors.InvalidArgument.WithMessage(
			fmt.Sprintf(
				"Invalid argument type, expected %s, got %v", reflect.TypeOf(CreateProjectRequest{}).Kind(), opt.Value,
			),
		))
	}
	request := opt.Value.(CreateProjectRequest)

	var projects []resource.Project
	if err := mgm.Coll(&resource.Project{}).SimpleFind(&projects, bson.M{NameDocumentKey: request.ProjectName, OwnerIDDocumentKey: request.OwnerID}); err != nil {
		panic(errors.InternalError.WithMessage(err.Error()))
	}
	if len(projects) > 0 {
		panic(errors.Conflict.WithMessage(fmt.Sprintf("Project with name %s owned by %s already exists", request.ProjectName, request.OwnerID)))
	}
	newProject := resource.NewProject(request.ProjectName, request.OwnerID)
	_ = repository.kubernetesDAO.Create(ctx, option.Option{
		Value: kubernetes.CreateNamespaceRequest{
			Name: newProject.Namespace,
		},
	})

	if err := mgm.Coll(newProject).Create(newProject); err != nil {
		panic(errors.ExternalServiceError.WithMessage(err.Error()))
	}
	logger.Info.Printf("Project %s created", newProject.ID)
	return CreateProjectResponse{
		ProjectID: newProject.ID.Hex(),
	}
}

func (repository *Repository) Get(ctx context.Context, opt option.Option) interface{} {
	if !opt.SetType(reflect.TypeOf(GetProjectByIDRequest{}).String()).Validate() {
		panic(errors.InvalidArgument.WithMessage(
			fmt.Sprintf(
				"Invalid argument type, expected %s, got %v", reflect.TypeOf(GetProjectByIDRequest{}).Kind(), opt.Value,
			),
		))
	}
	request := opt.Value.(GetProjectByIDRequest)

	project := &resource.Project{}

	err := mgm.Coll(project).FindByID(request.ProjectID, project)
	if err != nil {
		panic(errors.NotFound.WithMessage(fmt.Sprintf("Project with id %s not found", request.ProjectID)))
	}
	return GetProjectByIDResponse{
		Payload: *project,
	}
}

func (repository *Repository) Watch(_ context.Context, _ option.Option) interface{} {
	//TODO implement me
	panic("implement me")
}

func (repository *Repository) List(ctx context.Context, opt option.Option) interface{} {
	if !opt.SetType(reflect.TypeOf(GetProjectByOwnerIDRequest{}).String()).Validate() {
		panic(errors.InvalidArgument.WithMessage(
			fmt.Sprintf(
				"Invalid argument type, expected %s, got %v", reflect.TypeOf(GetProjectByOwnerIDRequest{}).Kind(), opt.Value,
			),
		))
	}
	request := opt.Value.(GetProjectByOwnerIDRequest)
	var projects []resource.Project
	if err := mgm.Coll(&resource.Project{}).SimpleFind(&projects, bson.M{OwnerIDDocumentKey: request.OwnerID}); err != nil {
		panic(errors.InternalError.WithMessage(err.Error()))
	}

	return GetProjectByOwnerIDResponse{
		Projects: projects,
	}
}

func (repository *Repository) Update(ctx context.Context, opt option.Option) interface{} {
	if !opt.SetType(reflect.TypeOf(UpdateProjectRequest{}).String()).Validate() {
		panic(errors.InvalidArgument.WithMessage(
			fmt.Sprintf(
				"Invalid argument type, expected %s, got %v", reflect.TypeOf(UpdateProjectRequest{}).Kind(), opt.Value,
			),
		))
	}
	request := opt.Value.(UpdateProjectRequest)
	var projectToUpdate resource.Project
	if request.ProjectName != "" {
		projectToUpdate = repository.Get(ctx, option.Option{
			Value: GetProjectByIDRequest{
				ProjectID: request.ProjectID,
			},
		}).(GetProjectByIDResponse).Payload
		projectToUpdate.Name = request.ProjectName
	} else if !reflect.DeepEqual(request.UpdatedProject, resource.Project{}) {
		projectToUpdate = request.UpdatedProject
	} else {
		panic(errors.InvalidArgument.WithMessage("A project or a project name must be provided"))
	}

	err := mgm.Coll(&projectToUpdate).Update(&projectToUpdate)
	if err != nil {
		panic(errors.ExternalServiceError.WithMessage(err.Error()))
	}
	return nil
}

func (repository *Repository) Delete(ctx context.Context, opt option.Option) {
	if !opt.SetType(reflect.TypeOf(DeleteRequest{}).String()).Validate() {
		panic(errors.InvalidArgument.WithMessage(
			fmt.Sprintf(
				"Invalid argument type, expected %s, got %v", reflect.TypeOf(DeleteRequest{}).Kind(), opt.Value,
			),
		))
	}
	request := opt.Value.(DeleteRequest)
	var project resource.Project
	if err := mgm.Coll(&resource.Project{}).FindByID(request.ProjectID, &project); err != nil {
		return
	}
	if err := mgm.Coll(&resource.Project{}).Delete(&project); err != nil {
		panic(errors.InternalError.WithMessage(err.Error()))
	}
}
