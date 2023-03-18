package dto

import "github.com/PaulBarrie/infra-worker/pkg/resource"

type CreateResourceRequest struct {
	ProjectID string
	Provider  resource.ProviderType
	Resource  interface{}
}

type CreateResourceResponse struct {
	ResourceID string
}
