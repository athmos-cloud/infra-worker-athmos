package project

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"go.mongodb.org/mongo-driver/bson"
)

func fromBsonRaw(raw bson.Raw) (resource.Project, errors.Error) {
	var resProject resource.Project
	if err := bson.Unmarshal(raw, &resProject); err != nil {
		return resource.Project{}, errors.InternalError.WithMessage(err.Error())
	}
	return resProject, errors.OK
}
