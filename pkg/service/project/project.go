package project

import (
	"context"
	"fmt"
	dto "github.com/PaulBarrie/infra-worker/pkg/common/dto/project"
	"github.com/PaulBarrie/infra-worker/pkg/domain"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/config"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/logger"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/option"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/utils"
	"github.com/PaulBarrie/infra-worker/pkg/repository"
	"github.com/PaulBarrie/infra-worker/pkg/repository/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type Service struct {
	ProjectRepository repository.IRepository
}

func (ps *Service) Create(ctx context.Context, request dto.CreateProjectRequest) (dto.CreateProjectResponse, errors.Error) {
	newProject := domain.Project{
		Name:      request.ProjectName,
		OwnerID:   request.OwnerID,
		Namespace: fmt.Sprintf("%s-%s", request.ProjectName, utils.RandomString(5)),
	}
	mongoResquest := mongo.CreateRequest{
		Payload:        newProject,
		CollectionName: config.Current.Mongo.ProjectCollection,
	}
	resp, err := mongo.Client.Create(ctx, option.Option{
		Value: mongoResquest,
	})
	if !err.IsOk() {
		return dto.CreateProjectResponse{}, err
	}
	return dto.CreateProjectResponse{
		ProjectID: resp.(mongo.CreateResponse).Id,
	}, errors.OK
}

func (ps *Service) Update(ctx context.Context, request dto.UpdateProjectRequest) errors.Error {
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

func (ps *Service) GetByID(ctx context.Context, request dto.GetProjectByIDRequest) (dto.GetProjectByIDResponse, errors.Error) {
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

func (ps *Service) GetByOwnerID(ctx context.Context, request dto.GetProjectByOwnerIDRequest) (dto.GetProjectByOwnerIDResponse, errors.Error) {
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
	var projects []domain.Project
	for _, p := range resp.(mongo.GetAllResponse).Payload {
		primitive := p.(bson.D)
		doc, errMarshal := bson.Marshal(primitive)
		if errMarshal != nil {
			logger.Error.Printf("Error marshalling bson: %v", errMarshal)
		}
		var projectItem domain.Project
		if errUnmarshall := bson.Unmarshal(doc, &projectItem); errUnmarshall != nil {
			logger.Error.Printf("Error unmarshalling bson: %v", errUnmarshall)
		}
		projects = append(projects, projectItem)
	}
	return dto.GetProjectByOwnerIDResponse{
		Payload: projects,
	}, errors.OK
}

func (ps *Service) Delete(ctx context.Context, request dto.DeleteRequest) errors.Error {
	mongoDeleteRequest := mongo.DeleteRequest{
		CollectionName: config.Current.Mongo.ProjectCollection,
		Id:             request.ProjectID,
	}
	if err := mongo.Client.Delete(ctx, option.Option{
		Value: mongoDeleteRequest,
	}); !err.IsOk() {
		return err
	}

	return errors.OK
}

func fromBsonRaw(raw bson.Raw) (domain.Project, errors.Error) {
	var project domain.Project
	if err := bson.Unmarshal(raw, &project); err != nil {
		return domain.Project{}, errors.InternalError.WithMessage(err.Error())
	}
	return project, errors.OK
}
