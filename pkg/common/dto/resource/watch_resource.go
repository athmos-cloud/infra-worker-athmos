package resource

import (
	"github.com/PaulBarrie/infra-worker/pkg/common"
)

type WatchResourceRequest struct {
	ResourceType common.ResourceType `json:"type"`
	ProjectID    string              `json:"project_id"`
	ResourceID   string              `json:"name"`
}

type WatchResourceResponse struct {
	Content interface{}
}
