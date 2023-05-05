package resource

import (
	"fmt"
	resourcePlugin "github.com/athmos-cloud/infra-worker-athmos/pkg/data/plugin"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/status"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/config"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/kamva/mgm/v3"
	"reflect"
)

type Network struct {
	mgm.DefaultModel `bson:",inline"`
	Metadata         metadata.Metadata     `bson:"metadata"`
	Identifier       identifier.Network    `bson:"identifier"`
	Status           status.ResourceStatus `bson:"status"`
	Subnetworks      SubnetworkCollection  `bson:"subnetworks,omitempty"`
	Firewalls        FirewallCollection    `bson:"firewalls,omitempty"`
}

func NewNetwork(id identifier.Network, providerType types.ProviderType) Network {
	return Network{
		Metadata:    metadata.New(metadata.CreateMetadataRequest{Name: id.ID}),
		Identifier:  id,
		Status:      status.New(id.ID, types.Network, providerType),
		Subnetworks: make(SubnetworkCollection),
		Firewalls:   make(FirewallCollection),
	}
}

type NetworkCollection map[string]Network

func (netCollection *NetworkCollection) Equals(other NetworkCollection) bool {
	if len(*netCollection) != len(other) {
		return false
	}
	for key, value := range *netCollection {
		if !value.Equals(other[key]) {
			return false
		}
	}
	return true
}

func (network *Network) New(id identifier.ID, providerType types.ProviderType) IResource {
	if reflect.TypeOf(id) != reflect.TypeOf(identifier.Network{}) {
		panic(errors.InvalidArgument.WithMessage("id type is not NetworkID"))
	}
	res := NewNetwork(id.(identifier.Network), providerType)
	return &res
}

func (network *Network) GetMetadata() metadata.Metadata {
	return network.Metadata
}

func (network *Network) SetMetadata(request metadata.CreateMetadataRequest) {
	network.Metadata = metadata.New(request)
}

func (network *Network) SetStatus(resourceStatus status.ResourceStatus) {
	network.Status = resourceStatus
}

func (network *Network) GetStatus() status.ResourceStatus {
	return network.Status
}

func (network *Network) GetPluginReference() resourcePlugin.Reference {
	if !network.Status.PluginReference.ChartReference.Empty() {
		return network.Status.PluginReference
	}
	switch network.Status.PluginReference.ResourceReference.ProviderType {
	case types.GCP:
		network.Status.PluginReference.ChartReference = resourcePlugin.HelmChartReference{
			ChartName:    config.Current.Plugins.Crossplane.GCP.Network.Chart,
			ChartVersion: config.Current.Plugins.Crossplane.GCP.Network.Version,
		}
		return network.Status.PluginReference
	}
	panic(errors.InvalidArgument.WithMessage(fmt.Sprintf("network type %s not supported", network.Status.PluginReference.ResourceReference.ProviderType)))
}

func (network *Network) FromMap(data map[string]interface{}) {
	if err := resourcePlugin.InjectMapIntoStruct(data, network); !err.IsOk() {
		panic(err)
	}
}

func (network *Network) Insert(project Project, update ...bool) {
	shouldUpdate := false
	if len(update) > 0 {
		shouldUpdate = update[0]
	}
	id := network.Identifier
	_, ok := project.Resources[id.ProviderID].VPCs[id.VPCID].Networks[id.ID]
	if !ok && shouldUpdate {
		panic(errors.NotFound.WithMessage(fmt.Sprintf("network %s not found in vpc %s", id.ID, id.VPCID)))
	}
	if ok && !shouldUpdate {
		panic(errors.Conflict.WithMessage(fmt.Sprintf("network %s already exists in vpc %s", id.ID, id.VPCID)))
	}
	project.Resources[id.ProviderID].VPCs[id.VPCID].Networks[id.ID] = *network
}

func (network *Network) Remove(project Project) {
	id := network.Identifier
	_, ok := project.Resources[id.ProviderID].VPCs[id.VPCID].Networks[id.ID]
	if !ok {
		panic(errors.NotFound.WithMessage(fmt.Sprintf("network %s not found in vpc %s", id.ID, id.VPCID)))
	}
	delete(project.Resources[id.ProviderID].VPCs[id.VPCID].Networks, id.ID)
}

func (network *Network) Equals(other Network) bool {
	return network.Metadata.Equals(other.Metadata) &&
		network.Identifier.Equals(other.Identifier) &&
		network.Status.Equals(other.Status) &&
		network.Subnetworks.Equals(other.Subnetworks) &&
		network.Firewalls.Equals(other.Firewalls)
}
