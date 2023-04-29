package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/common"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource"
)

type CreateResourceRequest struct {
	ProjectID     string                 `json:"project_id"`
	ProviderType  common.ProviderType    `json:"provider_type"`
	ResourceType  common.ResourceType    `json:"resource_type"`
	ResourceSpecs map[string]interface{} `json:"resource_specs"`
}

type CreateResourceResponse struct {
	Resource resource.IResource
}
