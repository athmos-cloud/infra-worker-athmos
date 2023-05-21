package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/metadata"
	"github.com/kamva/mgm/v3"
)

type VPC struct {
	mgm.DefaultModel
	Metadata       metadata.Metadata `json:"metadata"`
	IdentifierID   identifier.VPC    `json:"identifier"`
	IdentifierName identifier.VPC    `json:"identifierName"`
	Networks       NetworkCollection `json:"networks,omitempty"`
}

func (vpc *VPC) Equals(other VPC) bool {
	return vpc.Metadata.Equals(other.Metadata) &&
		vpc.IdentifierID.Equals(&other.IdentifierID) &&
		vpc.IdentifierName == other.IdentifierName &&
		vpc.Networks.Equals(other.Networks)
}

type VPCCollection map[string]VPC

func (collection *VPCCollection) Equals(other VPCCollection) bool {
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
