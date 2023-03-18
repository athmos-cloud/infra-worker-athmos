package project

import "github.com/PaulBarrie/infra-worker/pkg/resource/provider"

type Project struct {
	ID        string              `bson:"_id,omitempty"`
	Name      string              `bson:"name"`
	OwnerID   string              `bson:"owner"`
	Resources []provider.Provider `bson:"components"`
}
