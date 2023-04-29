package resource

import "github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"

type GetResourceRequest struct {
	ProjectID  string        `json:"projectID"`
	ResourceID identifier.ID `json:"resourceID"`
}

type GetResourceResponse struct {
	Content interface{}
}
