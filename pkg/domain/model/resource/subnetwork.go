package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/metadata"
)

type Subnetwork struct {
	Metadata       metadata.Metadata     `json:"metadata"`
	IdentifierID   identifier.Subnetwork `json:"identifierID"`
	IdentifierName identifier.Subnetwork `json:"identifierName"`
	Region         string                `json:"region"`
	IPCIDRRange    string                `json:"ipCIDRRange"`
	VMs            VMCollection          `json:"vms,omitempty"`
}

type SubnetworkCollection map[string]Subnetwork

func (collection *SubnetworkCollection) Equals(other SubnetworkCollection) bool {
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

func (subnet *Subnetwork) Equals(other Subnetwork) bool {
	return subnet.Metadata.Equals(other.Metadata) &&
		subnet.IdentifierID.Equals(&other.IdentifierID) &&
		subnet.IdentifierName.Equals(&other.IdentifierName) &&
		subnet.Region == other.Region &&
		subnet.IPCIDRRange == other.IPCIDRRange &&
		subnet.VMs.Equals(other.VMs)
}
