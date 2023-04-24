package project

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/project"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"go.mongodb.org/mongo-driver/bson"
)

func fromBsonRaw(raw bson.Raw) (domain.Project, errors.Error) {
	var project domain.Project
	if err := bson.Unmarshal(raw, &project); err != nil {
		return domain.Project{}, errors.InternalError.WithMessage(err.Error())
	}
	return project, errors.OK
}
