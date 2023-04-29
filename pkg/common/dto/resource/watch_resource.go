package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/common"
)

type WatchResourceRequest struct {
	ResourceType common.ResourceType `json:"type"`
	ProjectID    string              `json:"projectID"`
	ResourceID   string              `json:"name"`
}

type WatchResourceResponse struct {
	Content interface{}
}
