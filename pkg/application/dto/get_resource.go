package dto

import "github.com/PaulBarrie/infra-worker/pkg/resource"

type GetResourceRequest struct {
	Provider     resource.ProviderType `json:"provider"`
	ResourceType resource.ResourceType `json:"type"`
	ProjectID    string                `json:"project_id"`
	ResourceID   string                `json:"name"`
}

type GetResourceResponse struct {
	Content interface{}
}
