package project

import (
	"context"
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/dao/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/dao/mongo"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/config"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/logger"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
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
	exists := repository.MongoDAO.Exists(ctx, option.Option{
		Value: mongo.ExistsRequest{
			CollectionName: config.Current.Mongo.ProjectCollection,
			Filter:         bson.M{NameDocumentKey: request.ProjectName, OwnerIDDocumentKey: request.OwnerID},
		},
	})
	logger.Info.Printf("UpdatedProject %s owned by %s already exists: %v", request.ProjectName, request.OwnerID, exists)
	if exists {
		panic(errors.Conflict.WithMessage(fmt.Sprintf("UpdatedProject %s owned by %s already exists", request.ProjectName, request.OwnerID)))
	}

	newProject := resource.NewProject(request.ProjectName, request.OwnerID)
	_ = repository.kubernetesDAO.Create(ctx, option.Option{
		Value: kubernetes.CreateNamespaceRequest{
			Name: newProject.Namespace,
		},
	})

	mongoRequest := mongo.CreateRequest{
		Payload:        newProject,
		CollectionName: config.Current.Mongo.ProjectCollection,
	}
	resp := mongo.Client.Create(ctx, option.Option{
		Value: mongoRequest,
	})

	return CreateProjectResponse{
		ProjectID: resp.(mongo.CreateResponse).Id,
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
	mongoGetRequest := mongo.GetRequest{
		CollectionName: config.Current.Mongo.ProjectCollection,
		Id:             request.ProjectID,
	}
	resp := mongo.Client.Get(ctx, option.Option{
		Value: mongoGetRequest,
	})

	project := fromBsonRaw(resp.(mongo.GetResponse).Payload)

	return GetProjectByIDResponse{
		Payload: project,
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
	mongoGetRequest := mongo.GetAllRequest{
		CollectionName: config.Current.Mongo.ProjectCollection,
		Filter: bson.M{
			"owner_id": request.OwnerID,
		},
	}
	resp := mongo.Client.GetAll(ctx, option.Option{
		Value: mongoGetRequest,
	})

	// From mongo raw to UpdatedProject entity
	var projects []resource.Project
	for _, p := range resp.(mongo.GetAllResponse).Payload {
		primitive := p.(bson.D)
		doc, errMarshal := bson.Marshal(primitive)
		if errMarshal != nil {
			logger.Error.Printf("Error marshalling bson: %v", errMarshal)
		}
		var projectItem resource.Project
		if errUnmarshall := bson.Unmarshal(doc, &projectItem); errUnmarshall != nil {
			logger.Error.Printf("Error unmarshalling bson: %v", errUnmarshall)
		}
		projects = append(projects, projectItem)
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
		resp := mongo.Client.Get(ctx, option.Option{
			Value: mongo.GetRequest{
				CollectionName: config.Current.Mongo.ProjectCollection,
				Id:             request.ProjectID,
			},
		})
		projectRaw := resp.(mongo.GetResponse).Payload
		projectToUpdate = fromBsonRaw(projectRaw)
		projectToUpdate.Name = request.ProjectName
	} else if !reflect.DeepEqual(request.UpdatedProject, resource.Project{}) {
		projectToUpdate = request.UpdatedProject
	} else {
		panic(errors.InvalidArgument.WithMessage("A project or a project name must be provided"))
	}
	//bsonDoc, err := bson.Marshal(projectToUpdate)
	//if err != nil {
	//	panic(errors.InternalError.WithMessage(err.Error()))
	//}
	//var bsonDocMap map[string]interface{}
	//if errUnmarshal := bson.Unmarshal(bsonDoc, &bsonDocMap); errUnmarshal != nil {
	//	panic(errors.InternalError.WithMessage(errUnmarshal.Error()))
	//}
	mongo.Client.Update(ctx, option.Option{
		Value: mongo.UpdateRequest{
			CollectionName: config.Current.Mongo.ProjectCollection,
			Id:             request.ProjectID,
			Payload:        projectToUpdate,
		},
	})
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
	mongoDeleteRequest := mongo.DeleteRequest{
		CollectionName: config.Current.Mongo.ProjectCollection,
		Id:             request.ProjectID,
	}
	mongo.Client.Delete(ctx, option.Option{
		Value: mongoDeleteRequest,
	})
}
