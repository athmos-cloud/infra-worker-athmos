package project

import (
	"context"
	"fmt"
	dto "github.com/athmos-cloud/infra-worker-athmos/pkg/common/dto/project"
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
		ProjectRepository = &Repository{
			MongoDAO:      mongo.Client,
			kubernetesDAO: kubernetes.Client,
		}
	}
}

func (repository *Repository) Create(ctx context.Context, optn option.Option) (interface{}, errors.Error) {
	if !optn.SetType(reflect.TypeOf(dto.CreateProjectRequest{}).String()).Validate() {
		return nil, errors.InvalidArgument.WithMessage(
			fmt.Sprintf(
				"Invalid argument type, expected %s, got %v", reflect.TypeOf(dto.CreateProjectRequest{}).Kind(), optn.Value,
			),
		)
	}
	request := optn.Value.(dto.CreateProjectRequest)
	if err := repository.validateCreateUpdate(ctx, request.ProjectName, request.OwnerID); !err.IsOk() {
		return nil, err
	}

	newProject := resource.NewProject(request.ProjectName, request.OwnerID)
	_, err := repository.kubernetesDAO.Create(ctx, option.Option{
		Value: kubernetes.CreateNamespaceRequest{
			Name: newProject.Namespace,
		},
	})
	if !err.IsOk() {
		return nil, err
	}
	// Persist

	mongoRequest := mongo.CreateRequest{
		Payload:        newProject,
		CollectionName: config.Current.Mongo.ProjectCollection,
	}
	resp, err := mongo.Client.Create(ctx, option.Option{
		Value: mongoRequest,
	})
	if !err.IsOk() {
		return dto.CreateProjectResponse{}, err
	}
	return dto.CreateProjectResponse{
		ProjectID: resp.(mongo.CreateResponse).Id,
	}, errors.OK
}

func (repository *Repository) Get(ctx context.Context, optn option.Option) (interface{}, errors.Error) {
	if !optn.SetType(reflect.TypeOf(dto.GetProjectByIDRequest{}).String()).Validate() {
		return nil, errors.InvalidArgument.WithMessage(
			fmt.Sprintf(
				"Invalid argument type, expected %s, got %v", reflect.TypeOf(dto.GetProjectByIDRequest{}).Kind(), optn.Value,
			),
		)
	}
	request := optn.Value.(dto.GetProjectByIDRequest)
	mongoGetRequest := mongo.GetRequest{
		CollectionName: config.Current.Mongo.ProjectCollection,
		Id:             request.ProjectID,
	}
	resp, err := mongo.Client.Get(ctx, option.Option{
		Value: mongoGetRequest,
	})
	if !err.IsOk() {
		return dto.GetProjectByIDResponse{}, err
	}
	project, err := fromBsonRaw(resp.(mongo.GetResponse).Payload)
	if !err.IsOk() {
		return dto.GetProjectByIDResponse{}, err
	}
	return dto.GetProjectByIDResponse{
		Payload: project,
	}, errors.OK
}

func (repository *Repository) Watch(ctx context.Context, optn option.Option) (interface{}, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (repository *Repository) List(ctx context.Context, optn option.Option) (interface{}, errors.Error) {
	if !optn.SetType(reflect.TypeOf(dto.GetProjectByOwnerIDRequest{}).String()).Validate() {
		return nil, errors.InvalidArgument.WithMessage(
			fmt.Sprintf(
				"Invalid argument type, expected %s, got %v", reflect.TypeOf(dto.GetProjectByOwnerIDRequest{}).Kind(), optn.Value,
			),
		)
	}
	request := optn.Value.(dto.GetProjectByOwnerIDRequest)
	mongoGetRequest := mongo.GetAllRequest{
		CollectionName: config.Current.Mongo.ProjectCollection,
		Filter: bson.M{
			"owner_id": request.OwnerID,
		},
	}
	resp, err := mongo.Client.GetAll(ctx, option.Option{
		Value: mongoGetRequest,
	})
	if !err.IsOk() {
		return dto.GetProjectByOwnerIDResponse{}, err
	}

	// From mongo raw to Project entity
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

	return dto.GetProjectByOwnerIDResponse{
		Payload: projects,
	}, errors.OK
}

func (repository *Repository) Update(ctx context.Context, optn option.Option) errors.Error {
	if !optn.SetType(reflect.TypeOf(dto.UpdateProjectRequest{}).String()).Validate() {
		return errors.InvalidArgument.WithMessage(
			fmt.Sprintf(
				"Invalid argument type, expected %s, got %v", reflect.TypeOf(dto.UpdateProjectRequest{}).Kind(), optn.Value,
			),
		)
	}
	request := optn.Value.(dto.UpdateProjectRequest)
	mongoGetRequest := mongo.GetRequest{
		CollectionName: config.Current.Mongo.ProjectCollection,
		Id:             request.ProjectID,
	}
	resp, err := mongo.Client.Get(ctx, option.Option{
		Value: mongoGetRequest,
	})

	if !err.IsOk() {
		return err
	}
	projectRaw := resp.(mongo.GetResponse).Payload
	projectResp, err := fromBsonRaw(projectRaw)
	if !err.IsOk() {
		return err
	}

	// Check if (name, owner) does not exist
	if errValidate := repository.validateCreateUpdate(ctx, request.ProjectName, projectResp.OwnerID); !errValidate.IsOk() {
		return errValidate
	}

	projectResp.Name = request.ProjectName
	if err = mongo.Client.Update(ctx, option.Option{
		Value: mongo.UpdateRequest{
			CollectionName: config.Current.Mongo.ProjectCollection,
			Id:             request.ProjectID,
			Payload:        projectResp,
		},
	}); !err.IsOk() {
		return err
	}
	return errors.OK
}

func (repository *Repository) Delete(ctx context.Context, optn option.Option) errors.Error {
	if !optn.SetType(reflect.TypeOf(dto.DeleteRequest{}).String()).Validate() {
		return errors.InvalidArgument.WithMessage(
			fmt.Sprintf(
				"Invalid argument type, expected %s, got %v", reflect.TypeOf(dto.DeleteRequest{}).Kind(), optn.Value,
			),
		)
	}
	request := optn.Value.(dto.DeleteRequest)
	mongoDeleteRequest := mongo.DeleteRequest{
		CollectionName: config.Current.Mongo.ProjectCollection,
		Id:             request.ProjectID,
	}
	return mongo.Client.Delete(ctx, option.Option{
		Value: mongoDeleteRequest,
	})
}
