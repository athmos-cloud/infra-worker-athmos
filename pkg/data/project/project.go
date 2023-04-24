package domain

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/provider"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"reflect"
)

type Project struct {
	ID        string                 `bson:"_id,omitempty"`
	Name      string                 `bson:"name"`
	Namespace string                 `bson:"namespace"`
	OwnerID   string                 `bson:"owner_id"`
	Resources resources.ProviderList `bson:"providers"`
}

func (project *Project) Insert(resource domain.IResource) errors.Error {
	if reflect.TypeOf(resource) != reflect.TypeOf(&resources.Provider{}) {

	}
	provider := resource.(*resources.Provider)
	for _, p := range project.Resources {
		if p.Metadata.Name == provider.Metadata.Name {
			return errors.AlreadyExists.WithMessage(
				fmt.Sprintf("provider with name %s already exists", provider.Metadata.Name),
			)
		}
		if p.Metadata.ID == provider.Metadata.ID {
			return errors.AlreadyExists.WithMessage(
				fmt.Sprintf("provider with id %s already exists", provider.Metadata.ID),
			)
		}
	}
	project.Resources = append(project.Resources, *provider)
	return errors.OK
}
