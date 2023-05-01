package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

type CreateResourceRequest struct {
	ProjectID     string                 `json:"projectID"`
	Identifier    identifier.IdPayload   `json:"identifier"`
	ProviderType  types.ProviderType     `json:"providerType"`
	ResourceType  types.ResourceType     `json:"resourceType"`
	ResourceSpecs map[string]interface{} `json:"resourceSpecs"`
}

type CreateResourceResponse struct {
	Resource resource.IResource
}

type GetResourceRequest struct {
	ProjectID  string        `json:"projectID"`
	ResourceID identifier.ID `json:"resourceID"`
}

type GetResourceResponse struct {
	Content interface{}
}

type GetListResourceRequest struct {
	Provider     types.ProviderType `json:"provider"`
	ResourceType types.ResourceType `json:"type"`
	ProjectID    string             `json:"projectID"`
}

type UpdateResourceRequest struct {
	ProjectID        string
	ResourceID       identifier.IdPayload
	NewResourceSpecs map[string]interface{}
}

type DeleteResourceRequest struct {
	ProjectID  string               `json:"projectID"`
	ResourceID identifier.IdPayload `json:"resourceID"`
}

type DeleteResourceResponse struct {
	Message errors.Error `json:"message"`
}
