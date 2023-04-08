package domain

import (
	"fmt"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/utils"
)

type Metadata struct {
	ID               string            `bson:"_id"`
	Name             string            `bson:"name"`
	Monitored        bool              `bson:"monitored,default=true"`
	Tags             map[string]string `bson:"tags"`
	ReleaseReference ReleaseReference  `bson:"releaseReference"`
}

type CreateMetadataRequest struct {
	Name             string
	ProjectNamespace string
	NotMonitored     bool
	Tags             map[string]string
}

func New(request CreateMetadataRequest) Metadata {
	return Metadata{
		ID:        utils.GenerateUUID(),
		Name:      request.Name,
		Monitored: !request.NotMonitored,
		Tags:      request.Tags,
		ReleaseReference: ReleaseReference{
			Name:      fmt.Sprintf("%s-%s", request.Name, utils.GenerateUUID()),
			Namespace: request.ProjectNamespace,
			Versions:  []string{},
		},
	}
}

type ReleaseReference struct {
	Name      string   `bson:"name"`
	Namespace string   `bson:"namespace"`
	Versions  []string `bson:"versions"`
}
