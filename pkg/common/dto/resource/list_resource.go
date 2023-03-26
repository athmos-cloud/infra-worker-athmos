package resource

import (
	"github.com/PaulBarrie/infra-worker/pkg/common"
)

type GetListResourceRequest struct {
	Provider     common.ProviderType `json:"provider"`
	ResourceType common.ResourceType `json:"type"`
	ProjectID    string              `json:"project_id"`
}
