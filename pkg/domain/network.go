package domain

import "github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource"

type Network struct {
	Name        string
	Subnetworks map[string]Subnetwork
}

func FromNetworkDataMapper(network resource.Network) Network {
	return Network{
		Name:        network.Identifier.ID,
		Subnetworks: FromSubnetworkCollectionDataMapper(network.Subnetworks),
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
