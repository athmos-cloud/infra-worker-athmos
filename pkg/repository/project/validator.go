package project

import (
	"context"
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/dao/mongo"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/config"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
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
		return errors.Conflict.WithMessage(
			fmt.Sprintf("Project with name %s owned by %s already exists", projectName, ownerID),
		)
	}
	return errors.OK
}
