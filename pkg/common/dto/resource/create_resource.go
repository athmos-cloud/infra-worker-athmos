package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/common"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
)

type CreateResourceRequest struct {
	ProjectID     string                 `json:"projectID"`
	Identifier    identifier.IdPayload   `json:"identifier"`
	ProviderType  common.ProviderType    `json:"providerType"`
	ResourceType  common.ResourceType    `json:"resourceType"`
	ResourceSpecs map[string]interface{} `json:"resourceSpecs"`
}

type CreateResourceResponse struct {
	Resource resource.IResource
}
