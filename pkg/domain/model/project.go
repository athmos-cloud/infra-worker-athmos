package model

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/auth"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/utils"
	"github.com/kamva/mgm/v3"
	"regexp"
)

type Project struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string                      `bson:"name"`
	Namespace        string                      `bson:"namespace"`
	OwnerID          string                      `bson:"owner_id"`
	Resources        resource.ProviderCollection `bson:"providers"`
	Authentications  auth.List                   `bson:"authentications"`
}

func NewProject(name string, ownerID string) *Project {
	return &Project{
		Name:            name,
		Namespace:       namespaceFormat(fmt.Sprintf("%s-%s", name, utils.RandomString(5))),
		OwnerID:         ownerID,
		Resources:       make(resource.ProviderCollection),
		Authentications: make(auth.List),
	}
}

func namespaceFormat(namespace string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		panic(errors.InternalError.WithMessage(err.Error()))
	}
	return reg.ReplaceAllString(namespace, "")
}
