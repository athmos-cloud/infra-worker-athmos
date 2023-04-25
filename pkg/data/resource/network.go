package resource

import (
	"fmt"
	dto "github.com/athmos-cloud/infra-worker-athmos/pkg/common/dto/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

type Network struct {
	Metadata            metadata.Metadata       `bson:"metadata"`
	Identifier          identifier.Network      `bson:"identifier"`
	KubernetesResources kubernetes.ResourceList `bson:"kubernetesResources"`
	Subnetworks         SubnetworkCollection    `bson:"subnetworks"`
	Firewalls           FirewallCollection      `bson:"firewalls"`
}

func NewNetwork(id identifier.Network) Network {
	return Network{
		Metadata:            metadata.New(metadata.CreateMetadataRequest{Name: id.ID}),
		Identifier:          id,
		KubernetesResources: kubernetes.ResourceList{},
		Subnetworks:         make(SubnetworkCollection),
		Firewalls:           make(FirewallCollection),
	}
}

type NetworkCollection map[string]Network

func (network *Network) GetMetadata() metadata.Metadata {
	return network.Metadata
}

func (network *Network) WithMetadata(request metadata.CreateMetadataRequest) {
	network.Metadata = metadata.New(request)
}

func (network *Network) GetPluginReference(request dto.GetPluginReferenceRequest) (dto.GetPluginReferenceResponse, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (network *Network) FromMap(m map[string]interface{}) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (network *Network) Insert(project Project, update ...bool) errors.Error {
	shouldUpdate := false
	if len(update) > 0 {
		shouldUpdate = update[0]
	}
	id := network.Identifier
	_, ok := project.Resources[id.ProviderID].VPCs[id.VPCID].Networks[id.ID]
	if !ok && shouldUpdate {
		return errors.NotFound.WithMessage(fmt.Sprintf("network %s not found in vpc %s", id.ID, id.VPCID))
	}
	if ok && !shouldUpdate {
		return errors.Conflict.WithMessage(fmt.Sprintf("network %s already exists in vpc %s", id.ID, id.VPCID))
	}
	project.Resources[id.ProviderID].VPCs[id.VPCID].Networks[id.ID] = *network
	return errors.OK
}

func (network *Network) ToDomain() (interface{}, errors.Error) {
	//TODO implement me
	panic("implement me")
}
