package domain

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
)

type VPC struct {
	Name      string             `json:"name"`
	Monitored bool               `json:"monitored"`
	Type      types.ProviderType `json:"type"`
	Networks  NetworkCollection  `json:"networks"`
}

func (vpc VPC) ToDataMapper(resourceInput resource.IResource) resource.IResource {
	panic("implement me")
}

func FromVPCDataMapper(vpc *resource.VPC) VPC {
	return VPC{
		Name:      vpc.Identifier.ID,
		Monitored: vpc.Metadata.Managed,
		Networks:  FromNetworkCollectionDataMapper(vpc.Networks),
	}
}

type VPCCollection map[string]VPC

func FromVPCCollectionDataMapper(vpcs resource.VPCCollection) VPCCollection {
	result := make(VPCCollection)
	for _, vpc := range vpcs {
		result[vpc.Identifier.ID] = FromVPCDataMapper(&vpc)
	}
	return result
}
