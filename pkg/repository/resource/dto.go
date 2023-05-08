package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
)

type GetRequest struct {
	ResourceID identifier.ID
	Project    resource.Project
}

type GetResourceResponse struct {
	Resource resource.IResource
}

type CreateRequest struct {
	Project          resource.Project
	Name             string
	ProviderType     types.ProviderType
	ResourceType     types.ResourceType
	ParentIdentifier identifier.ID
	Monitored        bool
	Tags             map[string]string
	ResourceSpecs    map[string]interface{}
}

type CreateResponse struct {
	Resource resource.IResource
}

type UpdateRequest struct {
	Project          resource.Project
	ResourceID       identifier.ID
	NewResourceSpecs map[string]interface{}
}

type DeleteRequest struct {
	Project    resource.Project
	ResourceID identifier.ID
}
