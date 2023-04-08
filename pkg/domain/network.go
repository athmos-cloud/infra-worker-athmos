package domain

import (
	dto "github.com/PaulBarrie/infra-worker/pkg/common/dto/resource"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
)

type Network struct {
	Metadata    Metadata     `bson:"metadata"`
	Subnetworks []Subnetwork `bson:"subnetworks"`
	Firewalls   []Firewall   `bson:"firewalls"`
}

func (network *Network) GetMetadata() Metadata {
	return network.Metadata
}

func (network *Network) WithMetadata(request CreateMetadataRequest) {
	network.Metadata = New(request)
}

func (network *Network) GetPluginReference(request dto.GetPluginReferenceRequest) (dto.GetPluginReferenceResponse, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (network *Network) FromMap(m map[string]interface{}) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (network *Network) InsertIntoProject(project Project, upsert bool) errors.Error {
	//TODO implement me
	panic("implement me")
}
