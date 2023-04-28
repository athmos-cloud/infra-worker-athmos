package domain

import "github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource"

type Network struct {
	Monitored   bool                 `json:"monitored"`
	Name        string               `json:"name"`
	Subnetworks SubnetworkCollection `json:"subnetworks"`
	Firewalls   FirewallCollection   `json:"firewalls"`
}

func FromNetworkDataMapper(network resource.Network) Network {
	return Network{
		Monitored:   network.Metadata.Monitored,
		Name:        network.Identifier.ID,
		Subnetworks: FromSubnetworkCollectionDataMapper(network.Subnetworks),
		Firewalls:   FromFirewallCollectionDataMapper(network.Firewalls),
	}
}

type NetworkCollection map[string]Network

func FromNetworkCollectionDataMapper(networks resource.NetworkCollection) NetworkCollection {
	result := make(NetworkCollection)
	for _, network := range networks {
		result[network.Identifier.ID] = FromNetworkDataMapper(network)
	}
	return result
}
