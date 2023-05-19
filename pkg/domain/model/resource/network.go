package resource

import (
	identifier2 "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"reflect"
)

type Network struct {
	Metadata    metadata.Metadata    `bson:"metadata"`
	Identifier  identifier2.Network  `bson:"identifier"`
	Subnetworks SubnetworkCollection `bson:"subnetworks"`
	Firewalls   FirewallCollection   `bson:"firewalls"`
}

func NewNetwork(payload NewResourcePayload) Network {
	payload.Validate()
	var id identifier2.Network
	if reflect.TypeOf(payload.ParentIdentifier) == reflect.TypeOf(identifier2.Provider{}) {
		parentID := payload.ParentIdentifier.(identifier2.Provider)
		id = identifier2.Network{
			ProviderID: parentID.ProviderID,
			NetworkID:  formatResourceName(payload.Name),
		}
	} else if reflect.TypeOf(payload.ParentIdentifier) == reflect.TypeOf(identifier2.VPC{}) {
		parentID := payload.ParentIdentifier.(identifier2.VPC)
		id = identifier2.Network{
			ProviderID: parentID.ProviderID,
			VPCID:      parentID.VPCID,
			NetworkID:  formatResourceName(payload.Name),
		}
	} else {
		errors.InvalidArgument.WithMessage("ID type must be provider or VPC ID")
	}

	return Network{
		Metadata: metadata.New(metadata.CreateMetadataRequest{
			Name:         payload.Name,
			NotMonitored: !payload.Managed,
			Tags:         payload.Tags,
		}),
		Identifier:  id,
		Subnetworks: make(SubnetworkCollection),
		Firewalls:   make(FirewallCollection),
	}
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
		network.Identifier.Equals(other.Identifier) &&
		network.Subnetworks.Equals(other.Subnetworks) &&
		network.Firewalls.Equals(other.Firewalls)
}
