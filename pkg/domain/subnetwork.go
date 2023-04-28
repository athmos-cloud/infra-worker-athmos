package domain

import "github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource"

type Subnetwork struct {
	Name        string       `json:"name"`
	Monitored   bool         `json:"monitored"`
	IPCIDRRange string       `json:"ip_cidr_range"`
	Region      string       `json:"region"`
	VMs         VMCollection `json:"vms"`
}

func FromSubnetworkDataMapper(subnet resource.Subnetwork) Subnetwork {
	return Subnetwork{
		Name:        subnet.Identifier.ID,
		Monitored:   subnet.Metadata.Monitored,
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
