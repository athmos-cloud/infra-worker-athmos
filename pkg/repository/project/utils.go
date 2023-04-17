package project

import (
	"github.com/PaulBarrie/infra-worker/pkg/domain"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
	"go.mongodb.org/mongo-driver/bson"
)

func fromBsonRaw(raw bson.Raw) (domain.Project, errors.Error) {
	var project domain.Project
	if err := bson.Unmarshal(raw, &project); err != nil {
		return domain.Project{}, errors.InternalError.WithMessage(err.Error())
	}
	return project, errors.OK
}
