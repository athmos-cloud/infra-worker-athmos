package domain

import (
	dto "github.com/PaulBarrie/infra-worker/pkg/common/dto/resource"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
)

type Subnetwork struct {
	Metadata    Metadata `bson:"metadata"`
	Region      string   `bson:"region"`
	IPCIDRRange string   `bson:"ipCidrRange"`
	VMs         []VM     `bson:"vmList"`
}

func (subnet *Subnetwork) GetMetadata() Metadata {
	return subnet.Metadata
}

func (subnet *Subnetwork) WithMetadata(request CreateMetadataRequest) {
	subnet.Metadata = New(request)
}

func (subnet *Subnetwork) GetPluginReference(request dto.GetPluginReferenceRequest) (dto.GetPluginReferenceResponse, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (subnet *Subnetwork) FromMap(m map[string]interface{}) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (subnet *Subnetwork) InsertIntoProject(project Project, upsert bool) errors.Error {
	//TODO implement me
	panic("implement me")
}
