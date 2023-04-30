package metadata

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/utils"
)

type Metadata struct {
	Name             string            `bson:"name"`
	Monitored        bool              `bson:"monitored,default=true"`
	Tags             map[string]string `bson:"tags,omitempty"`
	ReleaseReference ReleaseReference  `bson:"releaseReference"`
}

type CreateMetadataRequest struct {
	Name             string
	ProjectNamespace string
	NotMonitored     bool
	Tags             map[string]string
}

func (metadata *Metadata) Equals(other Metadata) bool {
	return metadata.Name == other.Name &&
		metadata.Monitored == other.Monitored &&
		utils.MapEquals(metadata.Tags, other.Tags) &&
		metadata.ReleaseReference.Name == other.ReleaseReference.Name &&
		metadata.ReleaseReference.Namespace == other.ReleaseReference.Namespace &&
		utils.SliceEquals(metadata.ReleaseReference.Versions, other.ReleaseReference.Versions)
}

func New(request CreateMetadataRequest) Metadata {
	return Metadata{
		Name:      request.Name,
		Monitored: !request.NotMonitored,
		Tags:      request.Tags,
		ReleaseReference: ReleaseReference{
			Name:      fmt.Sprintf("%s-%s", request.Name, utils.GenerateUUID()),
			Namespace: request.ProjectNamespace,
			Versions:  make([]string, 0),
		},
	}
}

type ReleaseReference struct {
	Name      string   `bson:"name"`
	Namespace string   `bson:"namespace"`
	Versions  []string `bson:"versions"`
}