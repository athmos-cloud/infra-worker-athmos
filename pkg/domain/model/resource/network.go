package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/metadata"
)

type Network struct {
	Metadata       metadata.Metadata    `json:"metadata"`
	IdentifierID   identifier.Network   `json:"identifierID"`
	IdentifierName identifier.Network   `json:"identifierName"`
	Subnetworks    SubnetworkCollection `json:"subnetworks,omitempty"`
	Firewalls      FirewallCollection   `json:"firewalls,omitempty"`
}

type NetworkCollection map[string]Network

func (netCollection *NetworkCollection) Equals(other NetworkCollection) bool {
	if len(*netCollection) != len(other) {
		return false
	}
	for key, value := range *netCollection {
		if !value.Equals(other[key]) {
			return false
		}
	}
	return true
}

func (network *Network) Equals(other Network) bool {
	return network.Metadata.Equals(other.Metadata) &&
		network.Subnetworks.Equals(other.Subnetworks) &&
		network.Firewalls.Equals(other.Firewalls)
}
