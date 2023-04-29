package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/common"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource"
)

type CreateResourceRequest struct {
	ProjectID     string                 `json:"projectId"`
	ProviderType  common.ProviderType    `json:"providerType"`
	ResourceType  common.ResourceType    `json:"resourceType"`
	ResourceSpecs map[string]interface{} `json:"resourceSpecs"`
}

type CreateResourceResponse struct {
	Resource resource.IResource
}
