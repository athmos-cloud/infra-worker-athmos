package domain

import "github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource"

type Project struct {
	ID        string
	Name      string
	Owner     string
	Providers ProviderCollection
}

func FromProjectDataMapper(project resource.Project) Project {
	return Project{
		ID:        project.ID,
		Name:      project.Name,
		Owner:     project.OwnerID,
		Providers: FromProviderCollectionDataMapper(project.Resources),
	}
}
