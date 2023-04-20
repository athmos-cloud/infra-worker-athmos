package resource

import (
	"github.com/PaulBarrie/infra-worker/pkg/common"
	"helm.sh/helm/v3/pkg/release"
)

type CreateResourceRequest struct {
	ProjectID     string                 `json:"project_id"`
	ProviderType  common.ProviderType    `json:"provider_type"`
	ResourceType  common.ResourceType    `json:"resource_type"`
	ResourceSpecs map[string]interface{} `json:"resource_specs"`
}

type CreateResourceResponse struct {
	ResourceID  string
	HelmRelease *release.Release
}
