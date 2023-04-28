package domain

import "github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource"

type Subnetwork struct {
	Name        string
	IPCIDRRange string
	Region      string `bson:"region"`
	VMs         VMCollection
}

func FromSubnetworkDataMapper(subnet resource.Subnetwork) Subnetwork {
	return Subnetwork{
		Name:        subnet.Identifier.ID,
		IPCIDRRange: subnet.IPCIDRRange,
		Region:      subnet.Region,
		VMs:         FromVMCollectionDataMapper(subnet.VMs),
	}
}

type SubnetworkCollection map[string]Subnetwork

func FromSubnetworkCollectionDataMapper(subnets resource.SubnetworkCollection) SubnetworkCollection {
	result := make(SubnetworkCollection)
	for _, subnet := range subnets {
		result[subnet.Identifier.ID] = FromSubnetworkDataMapper(subnet)
	}
	return result
}
