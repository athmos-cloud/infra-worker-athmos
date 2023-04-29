package resource

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/common"
	resourcePlugin "github.com/athmos-cloud/infra-worker-athmos/pkg/data/plugin"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/status"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/config"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"reflect"
)

type Subnetwork struct {
	Metadata    metadata.Metadata     `bson:"metadata"`
	Identifier  identifier.Subnetwork `bson:"hierarchyLocation"`
	Status      status.ResourceStatus `bson:"status"`
	VPC         string                `bson:"subnet" plugin:"subnet"`
	Network     string                `bson:"network" plugin:"network"`
	Region      string                `bson:"region" plugin:"region"`
	IPCIDRRange string                `bson:"ipCidrRange" plugin:"ipCidrRange"`
	VMs         VMCollection          `bson:"vmList"`
}

func NewSubnetwork(id identifier.Subnetwork, providerType common.ProviderType) Subnetwork {
	return Subnetwork{
		Metadata:   metadata.New(metadata.CreateMetadataRequest{Name: id.ID}),
		Identifier: id,
		Status:     status.New(id.ID, common.Subnetwork, providerType),
		VMs:        make(VMCollection),
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

func (subnet *Subnetwork) New(id identifier.ID, provider common.ProviderType) (IResource, errors.Error) {
	if reflect.TypeOf(id) != reflect.TypeOf(identifier.Subnetwork{}) {
		return nil, errors.InvalidArgument.WithMessage("id type is not SubnetworkID")
	}
	res := NewSubnetwork(id.(identifier.Subnetwork), provider)
	return &res, errors.OK
}

func (subnet *Subnetwork) GetMetadata() metadata.Metadata {
	return subnet.Metadata
}

func (subnet *Subnetwork) SetMetadata(request metadata.CreateMetadataRequest) {
	subnet.Metadata = metadata.New(request)
}

func (subnet *Subnetwork) SetStatus(resourceStatus status.ResourceStatus) {
	subnet.Status = resourceStatus
}

func (subnet *Subnetwork) GetStatus() status.ResourceStatus {
	return subnet.Status
}

func (subnet *Subnetwork) GetPluginReference() (resourcePlugin.Reference, errors.Error) {
	if !subnet.Status.PluginReference.ChartReference.Empty() {
		return subnet.Status.PluginReference, errors.OK
	}
	switch subnet.Status.PluginReference.ResourceReference.ProviderType {
	case common.GCP:
		subnet.Status.PluginReference.ChartReference = resourcePlugin.HelmChartReference{
			ChartName:    config.Current.Plugins.Crossplane.GCP.Subnet.Chart,
			ChartVersion: config.Current.Plugins.Crossplane.GCP.Subnet.Version,
		}
		return subnet.Status.PluginReference, errors.OK
	}
	return resourcePlugin.Reference{}, errors.InvalidArgument.WithMessage(fmt.Sprintf("provider type %s not supported", subnet.Status.PluginReference.ResourceReference.ProviderType))
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
		subnet.Status.Equals(other.Status) &&
		subnet.VPC == other.VPC &&
		subnet.Network == other.Network &&
		subnet.Region == other.Region &&
		subnet.IPCIDRRange == other.IPCIDRRange &&
		subnet.VMs.Equals(other.VMs)
}
