package resource

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/utils"
)

type Project struct {
	ID        string             `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Namespace string             `bson:"namespace"`
	OwnerID   string             `bson:"owner_id"`
	Resources ProviderCollection `bson:"providers"`
}

func NewProject(name string, ownerID string) Project {
	return Project{
		Name:      name,
		Namespace: fmt.Sprintf("%s-%s", name, utils.RandomString(5)),
		OwnerID:   ownerID,
		Resources: make(ProviderCollection, 10000),
	}
}

func (project *Project) Insert(resource IResource) *Project {
	resource.Insert(*project)
	return project
}

func (project *Project) Update(resource IResource) *Project {
	resource.Insert(*project, true)
	return project
}

func (project *Project) Delete(resource IResource) *Project {
	resource.Remove(*project)
	return project
}
