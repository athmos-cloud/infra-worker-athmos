package domain

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
)

type Network struct {
	Monitored    bool                 `json:"monitored"`
	Name         string               `json:"name"`
	ProviderType types.ProviderType   `json:"providerType"`
	Subnetworks  SubnetworkCollection `json:"subnetworks"`
	Firewalls    FirewallCollection   `json:"firewalls"`
}

func (network Network) ToDataMapper(resourceInput resource.IResource) resource.IResource {
	networkInput := resourceInput.(*resource.Network)
	networkInput.Identifier.ID = network.Name
	networkInput.Metadata.Managed = network.Monitored
	networkInput.Status.PluginReference.ResourceReference.ProviderType = network.ProviderType
	return networkInput
}

func FromNetworkDataMapper(network *resource.Network) Network {
	return Network{
		Monitored:   network.Metadata.Managed,
		Name:        network.Identifier.ID,
		Subnetworks: FromSubnetworkCollectionDataMapper(network.Subnetworks),
		Firewalls:   FromFirewallCollectionDataMapper(network.Firewalls),
	}
}

type NetworkCollection map[string]Network

func FromNetworkCollectionDataMapper(networks resource.NetworkCollection) NetworkCollection {
	result := make(NetworkCollection)
	for _, network := range networks {
		result[network.Identifier.ID] = FromNetworkDataMapper(&network)
	}
	return result
}
