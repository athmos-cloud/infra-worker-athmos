package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/common"
)

type GetListResourceRequest struct {
	Provider     common.ProviderType `json:"provider"`
	ResourceType common.ResourceType `json:"type"`
	ProjectID    string              `json:"projectID"`
}
