package domain

import "github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource"

type VPC struct {
	Name      string            `json:"name"`
	Monitored bool              `json:"monitored"`
	Networks  NetworkCollection `json:"networks"`
}

func FromVPCDataMapper(vpc resource.VPC) VPC {
	return VPC{
		Name:      vpc.Identifier.ID,
		Monitored: vpc.Metadata.Monitored,
		Networks:  FromNetworkCollectionDataMapper(vpc.Networks),
	}
}

type VPCCollection map[string]VPC

func FromVPCCollectionDataMapper(vpcs resource.VPCCollection) VPCCollection {
	result := make(VPCCollection)
	for _, vpc := range vpcs {
		result[vpc.Identifier.ID] = FromVPCDataMapper(vpc)
	}
	return result
}
