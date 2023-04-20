package resource

import "github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"

type DeleteResourceRequest struct {
	ProjectID    string `json:"project_id"`
	ResourceType string `json:"type"`
	ResourceID   string `json:"resource_id"`
}

type DeleteResourceResponse struct {
	Message errors.Error `json:"message"`
}
