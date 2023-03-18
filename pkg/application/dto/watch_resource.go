package dto

import "github.com/PaulBarrie/infra-worker/pkg/resource"

type WatchResourceRequest struct {
	ResourceType resource.ResourceType `json:"type"`
	ProjectID    string                `json:"project_id"`
	ResourceID   string                `json:"name"`
}

type WatchResourceResponse struct {
	Content interface{}
}
