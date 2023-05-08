package domain

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
)

type Subnetwork struct {
	Name         string             `json:"name"`
	Monitored    bool               `json:"monitored"`
	ProviderType types.ProviderType `json:"providerType"`
	IPCIDRRange  string             `json:"IPCIDRRange²"`
	Region       string             `json:"region"`
	VMs          VMCollection       `json:"vms"`
}

func (subnet Subnetwork) ToDataMapper(resourceInput resource.IResource) resource.IResource {
	subnetInput := resourceInput.(*resource.Subnetwork)
	subnetInput.Identifier.SubnetworkID = subnet.Name
	subnetInput.Metadata.Managed = subnet.Monitored
	subnetInput.IPCIDRRange = subnet.IPCIDRRange
	subnetInput.Region = subnet.Region
	return subnetInput
}

func FromSubnetworkDataMapper(subnet *resource.Subnetwork) Subnetwork {
	return Subnetwork{
		Name:         subnet.Identifier.SubnetworkID,
		ProviderType: subnet.GetPluginReference().ResourceReference.ProviderType,
		Monitored:    subnet.Metadata.Managed,
		IPCIDRRange:  subnet.IPCIDRRange,
		Region:       subnet.Region,
		VMs:          FromVMCollectionDataMapper(subnet.VMs),
	}
}

type SubnetworkCollection map[string]Subnetwork

func FromSubnetworkCollectionDataMapper(subnets resource.SubnetworkCollection) SubnetworkCollection {
	result := make(SubnetworkCollection)
	for _, subnet := range subnets {
		result[subnet.Identifier.SubnetworkID] = FromSubnetworkDataMapper(&subnet)
	}
	return result
}
