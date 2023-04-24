package resources

import (
	dto "github.com/athmos-cloud/infra-worker-athmos/pkg/common/dto/resource"
	domain2 "github.com/athmos-cloud/infra-worker-athmos/pkg/data/project"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/firewall"
	resources2 "github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/subnetwork"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

type Network struct {
	Metadata    domain.Metadata                  `bson:"metadata"`
	Subnetworks map[string]resources2.Subnetwork `bson:"subnetworks"`
	Firewalls   map[string]resources.Firewall    `bson:"firewalls"`
}

type NetworkHierarchyLocation struct {
	ProviderID string
	VPCID      string
}

func (network *Network) GetMetadata() domain.Metadata {
	return network.Metadata
}

func (network *Network) WithMetadata(request domain.CreateMetadataRequest) {
	network.Metadata = domain.New(request)
}

func (network *Network) GetPluginReference(request dto.GetPluginReferenceRequest) (dto.GetPluginReferenceResponse, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (network *Network) FromMap(m map[string]interface{}) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (network *Network) InsertIntoProject(project domain2.Project, upsert bool) errors.Error {
	//TODO implement me
	panic("implement me")
}
