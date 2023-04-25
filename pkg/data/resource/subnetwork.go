package resource

import (
	"fmt"
	dto "github.com/athmos-cloud/infra-worker-athmos/pkg/common/dto/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

type Subnetwork struct {
	Metadata            metadata.Metadata       `bson:"metadata"`
	Identifier          identifier.Subnetwork   `bson:"hierarchyLocation"`
	KubernetesResources kubernetes.ResourceList `bson:"kubernetesResources"`
	VPC                 string                  `bson:"vpc"`
	Network             string                  `bson:"network"`
	Region              string                  `bson:"region"`
	IPCIDRRange         string                  `bson:"ipCidrRange"`
	VMs                 VMCollection            `bson:"vmList"`
}

func NewSubnetwork(id identifier.Subnetwork) Subnetwork {
	return Subnetwork{
		Metadata:            metadata.New(metadata.CreateMetadataRequest{Name: id.ID}),
		Identifier:          id,
		KubernetesResources: kubernetes.ResourceList{},
		VMs:                 make(VMCollection),
	}
}

type SubnetworkCollection map[string]Subnetwork

func (subnet *Subnetwork) GetMetadata() metadata.Metadata {
	return subnet.Metadata
}

func (subnet *Subnetwork) WithMetadata(request metadata.CreateMetadataRequest) {
	subnet.Metadata = metadata.New(request)
}

func (subnet *Subnetwork) GetPluginReference(request dto.GetPluginReferenceRequest) (dto.GetPluginReferenceResponse, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (subnet *Subnetwork) FromMap(m map[string]interface{}) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (subnet *Subnetwork) Insert(project Project, update ...bool) errors.Error {
	shouldUpdate := false
	if len(update) > 0 {
		shouldUpdate = update[0]
	}
	id := subnet.Identifier
	if _, ok := project.Resources[id.ProviderID].VPCs[id.VPCID].Networks[id.NetworkID].Subnetworks[id.ID]; !ok && !shouldUpdate {
		return errors.NotFound.WithMessage(fmt.Sprintf("subnet %s not found in network %s", id.ID, id.NetworkID))
	}
	project.Resources[id.ProviderID].VPCs[id.VPCID].Networks[id.NetworkID].Subnetworks[id.ID] = *subnet
	return errors.OK
}

func (subnet *Subnetwork) ToDomain() (interface{}, errors.Error) {
	//TODO implement me
	panic("implement me")
}
