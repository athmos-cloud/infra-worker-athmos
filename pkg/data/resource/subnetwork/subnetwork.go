package resources

import (
	dto "github.com/athmos-cloud/infra-worker-athmos/pkg/common/dto/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/kubernetes"
	domain2 "github.com/athmos-cloud/infra-worker-athmos/pkg/data/project"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/vm"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

type Subnetwork struct {
	Metadata            domain.Metadata         `bson:"metadata"`
	Position            Identifier              `bson:"hierarchyLocation"`
	KubernetesResources kubernetes.ResourceList `bson:"kubernetesResources"`
	VPC                 string                  `bson:"vpc"`
	Network             string                  `bson:"network"`
	Region              string                  `bson:"region"`
	IPCIDRRange         string                  `bson:"ipCidrRange"`
	VMs                 map[string]resources.VM `bson:"vmList"`
}

type Identifier struct {
	ID         string
	ProviderID string
	VPCID      string
}

func (subnet *Subnetwork) GetMetadata() domain.Metadata {
	return subnet.Metadata
}

func (subnet *Subnetwork) WithMetadata(request domain.CreateMetadataRequest) {
	subnet.Metadata = domain.New(request)
}

func (subnet *Subnetwork) GetPluginReference(request dto.GetPluginReferenceRequest) (dto.GetPluginReferenceResponse, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (subnet *Subnetwork) FromMap(m map[string]interface{}) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (subnet *Subnetwork) InsertIntoProject(project domain2.Project, upsert bool) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (subnet *Subnetwork) ToDomain() (interface{}, errors.Error) {
	//TODO implement me
	panic("implement me")
}
