package project

import (
	"context"
	"fmt"
	"github.com/PaulBarrie/infra-worker/pkg/dao/mongo"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/config"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/option"
	"go.mongodb.org/mongo-driver/bson"
)

func (repository *Repository) validateCreateUpdate(ctx context.Context, projectName string, ownerID string) errors.Error {
	exists, err := repository.MongoDAO.Exists(ctx, option.Option{
		Value: mongo.ExistsRequest{
			CollectionName: config.Current.Mongo.ProjectCollection,
			Filter:         bson.M{NameDocumentKey: projectName, OwnerIDDocumentKey: ownerID},
		},
	})
	if !err.IsOk() {
		return err
	}
	if exists {
		return errors.AlreadyExists.WithMessage(
			fmt.Sprintf("Project with name %s owned by %s already exists", projectName, ownerID),
		)
	}
	return errors.OK
}
