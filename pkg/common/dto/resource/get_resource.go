package resource

import "github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"

type GetResourceRequest struct {
	ProjectID  string        `json:"project_id"`
	ResourceID identifier.ID `json:"resource_id"`
}

type GetResourceResponse struct {
	Content interface{}
}
