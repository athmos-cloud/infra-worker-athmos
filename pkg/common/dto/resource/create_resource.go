package resource

import (
	"github.com/PaulBarrie/infra-worker/pkg/common"
)

type CreateResourceRequest struct {
	ProjectID string
	Provider  common.Plugin
	Resource  interface{}
}

type CreateResourceResponse struct {
	ResourceID string
}
