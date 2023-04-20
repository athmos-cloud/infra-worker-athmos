package resource

import "github.com/athmos-cloud/infra-worker-athmos/pkg/common"

type CreateRequest struct {
	ProjectNamespace string
	ProviderType     common.ProviderType
	ResourceType     common.ResourceType
	ResourceSpecs    map[string]interface{}
}

type UpdateRequest struct {
	ReleaseName      string
	ReleaseVersion   string
	NewResourceSpecs map[string]interface{}
}
