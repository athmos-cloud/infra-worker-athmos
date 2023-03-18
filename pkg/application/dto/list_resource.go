package dto

import "github.com/PaulBarrie/infra-worker/pkg/resource"

type GetListResourceRequest struct {
	Provider     resource.ProviderType `json:"provider"`
	ResourceType resource.ResourceType `json:"type"`
	ProjectID    string                `json:"project_id"`
}
