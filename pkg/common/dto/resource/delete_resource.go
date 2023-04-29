package resource

import "github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"

type DeleteResourceRequest struct {
	ProjectID    string `json:"projectID"`
	ResourceType string `json:"type"`
	ResourceID   string `json:"resourceID"`
}

type DeleteResourceResponse struct {
	Message errors.Error `json:"message"`
}
