package resource

import (
	"github.com/PaulBarrie/infra-worker/pkg/common"
)

type CreateResourceRequest struct {
	ProjectID     string
	Provider      common.ProviderType
	ResourceType  common.ResourceType
	ResourceSpecs interface{}
}

type CreateResourceResponse struct {
	ResourceID string
}
