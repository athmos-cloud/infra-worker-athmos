package metadata

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/utils"
	"github.com/kamva/mgm/v3"
)

type Metadata struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string            `bson:"name"`
	Managed          bool              `bson:"managed,default=true" plugin:"managed"`
	Tags             map[string]string `bson:"tags,omitempty"`
}

type CreateMetadataRequest struct {
	Name         string
	NotMonitored bool
	Tags         map[string]string
}

func (metadata *Metadata) Equals(other Metadata) bool {
	return metadata.Name == other.Name &&
		metadata.Managed == other.Managed &&
		utils.MapEquals(metadata.Tags, other.Tags)
}

func New(request CreateMetadataRequest) Metadata {
	return Metadata{
		Name:    request.Name,
		Managed: !request.NotMonitored,
		Tags:    request.Tags,
	}
}
