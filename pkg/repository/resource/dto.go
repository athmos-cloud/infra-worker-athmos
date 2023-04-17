package resource

import "github.com/PaulBarrie/infra-worker/pkg/common"

type CreateRequest struct {
	ProjectNamespace string
	ProviderType     common.ProviderType
	ResourceType     common.ResourceType
	ResourceSpecs    map[string]interface{}
}
