package resource

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/common"
	dto "github.com/athmos-cloud/infra-worker-athmos/pkg/common/dto/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/kubernetes"
	resourcePlugin "github.com/athmos-cloud/infra-worker-athmos/pkg/data/plugin"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/config"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"reflect"
)

type Subnetwork struct {
	Metadata            metadata.Metadata       `bson:"metadata"`
	Identifier          identifier.Subnetwork   `bson:"hierarchyLocation"`
	KubernetesResources kubernetes.ResourceList `bson:"kubernetesResources"`
	VPC                 string                  `bson:"vpc" plugin:"vpc"`
	Network             string                  `bson:"network" plugin:"network"`
	Region              string                  `bson:"region" plugin:"region"`
	IPCIDRRange         string                  `bson:"ipCidrRange" plugin:"ipCidrRange"`
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

func (collection *SubnetworkCollection) Equals(other SubnetworkCollection) bool {
	if len(*collection) != len(other) {
		return false
	}
	for key, value := range *collection {
		if !value.Equals(other[key]) {
			return false
		}
	}
	return true
}

func (subnet *Subnetwork) New(id identifier.ID) (IResource, errors.Error) {
	if reflect.TypeOf(id) != reflect.TypeOf(identifier.Subnetwork{}) {
		return nil, errors.InvalidArgument.WithMessage("id type is not SubnetworkID")
	}
	res := NewSubnetwork(id.(identifier.Subnetwork))
	return &res, errors.OK
}

func (subnet *Subnetwork) GetMetadata() metadata.Metadata {
	return subnet.Metadata
}

func (subnet *Subnetwork) WithMetadata(request metadata.CreateMetadataRequest) {
	subnet.Metadata = metadata.New(request)
}

func (subnet *Subnetwork) GetPluginReference(request dto.GetPluginReferenceRequest) (dto.GetPluginReferenceResponse, errors.Error) {
	switch request.ProviderType {
	case common.GCP:
		return dto.GetPluginReferenceResponse{
			ChartName:    config.Current.Plugins.Crossplane.GCP.Subnet.Chart,
			ChartVersion: config.Current.Plugins.Crossplane.GCP.Subnet.Version,
		}, errors.Error{}
	}
	return dto.GetPluginReferenceResponse{}, errors.InvalidArgument.WithMessage(fmt.Sprintf("provider type %s not supported", request.ProviderType))
}

func (subnet *Subnetwork) FromMap(data map[string]interface{}) errors.Error {
	return resourcePlugin.InjectMapIntoStruct(data, subnet)

}

func (subnet *Subnetwork) Insert(project Project, update ...bool) errors.Error {
	shouldUpdate := false
	if len(update) > 0 {
		shouldUpdate = update[0]
	}
	id := subnet.Identifier
	_, ok := project.Resources[id.ProviderID].VPCs[id.VPCID].Networks[id.NetworkID].Subnetworks[id.ID]
	if !ok && shouldUpdate {
		return errors.NotFound.WithMessage(fmt.Sprintf("subnet %s not found in network %s", id.ID, id.NetworkID))
	}
	if ok && !shouldUpdate {
		return errors.Conflict.WithMessage(fmt.Sprintf("subnet %s already exists in network %s", id.ID, id.NetworkID))
	}
	project.Resources[id.ProviderID].VPCs[id.VPCID].Networks[id.NetworkID].Subnetworks[id.ID] = *subnet
	return errors.OK
}

func (subnet *Subnetwork) Remove(project Project) errors.Error {
	id := subnet.Identifier
	_, ok := project.Resources[id.ProviderID].VPCs[id.VPCID].Networks[id.NetworkID].Subnetworks[id.ID]
	if !ok {
		return errors.NotFound.WithMessage(fmt.Sprintf("subnet %s not found in network %s", id.ID, id.NetworkID))
	}
	delete(project.Resources[id.ProviderID].VPCs[id.VPCID].Networks[id.NetworkID].Subnetworks, id.ID)
	return errors.NoContent
}

func (subnet *Subnetwork) Equals(other Subnetwork) bool {
	return subnet.Metadata.Equals(other.Metadata) &&
		subnet.Identifier.Equals(other.Identifier) &&
		subnet.VPC == other.VPC &&
		subnet.Network == other.Network &&
		subnet.Region == other.Region &&
		subnet.IPCIDRRange == other.IPCIDRRange &&
		subnet.VMs.Equals(other.VMs)
}
