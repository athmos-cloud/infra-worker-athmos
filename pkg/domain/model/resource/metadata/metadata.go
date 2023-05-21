package metadata

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/utils"
)

type Metadata struct {
	Namespace string            `json:"namespace"`
	Managed   bool              `json:"managed"`
	Tags      map[string]string `json:"tags,omitempty"`
}

type CreateMetadataRequest struct {
	Name         string
	NotMonitored bool
	Tags         map[string]string
}

func (metadata *Metadata) Equals(other Metadata) bool {
	return metadata.Managed == other.Managed &&
		utils.MapEquals(metadata.Tags, other.Tags)
}

func New(request CreateMetadataRequest) Metadata {
	return Metadata{
		Managed: !request.NotMonitored,
		Tags:    request.Tags,
	}
}
