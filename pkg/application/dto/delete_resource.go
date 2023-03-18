package dto

import "github.com/PaulBarrie/infra-worker/pkg/kernel/errors"

type DeleteResourceRequest struct {
	ProjectID    string `json:"project_id"`
	ResourceType string `json:"type"`
	ResourceID   string `json:"resource_id"`
}

type DeleteResourceResponse struct {
	Message errors.Error `json:"message"`
}
