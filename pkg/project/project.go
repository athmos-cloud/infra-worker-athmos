package project

import (
	"github.com/PaulBarrie/infra-worker/pkg/plugin"
	"github.com/PaulBarrie/infra-worker/pkg/resource/provider"
)

type Project struct {
	ID        string              `bson:"_id,omitempty"`
	Name      string              `bson:"name"`
	Namespace string              `bson:"namespace"`
	OwnerID   string              `bson:"owner_id"`
	Resources []provider.Provider `bson:"components"`
}

type Resource struct {
	ID               string           `bson:"_id,omitempty"`
	Plugin           plugin.Plugin    `bson:"plugin"`
	ReleaseReference ReleaseReference `bson:"releaseReference"`
}

type ReleaseReference struct {
	Name      string `bson:"name"`
	Namespace string `bson:"namespace"`
}
