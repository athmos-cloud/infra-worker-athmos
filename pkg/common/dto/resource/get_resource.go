package resource

import "github.com/PaulBarrie/infra-worker/pkg/common"

type GetResourceRequest struct {
	Provider     common.ProviderType `json:"provider"`
	ResourceType common.ResourceType `json:"type"`
	ProjectID    string              `json:"project_id"`
	ResourceID   string              `json:"name"`
}

type GetResourceResponse struct {
	Content interface{}
}
