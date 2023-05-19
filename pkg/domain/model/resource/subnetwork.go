package resource

import (
	identifier2 "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"reflect"
)

type Subnetwork struct {
	Metadata    metadata.Metadata      `bson:"metadata"`
	Identifier  identifier2.Subnetwork `bson:"hierarchyLocation"`
	Region      string                 `bson:"region" plugin:"region"`
	IPCIDRRange string                 `bson:"ipCidrRange" plugin:"ipCidrRange"`
	VMs         VMCollection           `bson:"vmList"`
}

func NewSubnetwork(payload NewResourcePayload) Subnetwork {
	payload.Validate()
	if reflect.TypeOf(payload.ParentIdentifier) != reflect.TypeOf(identifier2.Network{}) {
		panic(errors.InvalidArgument.WithMessage("ID type must be network ID"))
	}
	parentID := payload.ParentIdentifier.(identifier2.Network)
	id := identifier2.Subnetwork{
		ProviderID:   parentID.ProviderID,
		NetworkID:    parentID.NetworkID,
		VPCID:        parentID.VPCID,
		SubnetworkID: formatResourceName(payload.Name),
	}

	return Subnetwork{
		Metadata: metadata.New(metadata.CreateMetadataRequest{
			Name:         id.ProviderID,
			NotMonitored: !payload.Managed,
			Tags:         payload.Tags,
		}),
		Identifier: id,
		VMs:        make(VMCollection),
	}
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
		subnet.Identifier.Equals(other.Identifier) &&
		subnet.Region == other.Region &&
		subnet.IPCIDRRange == other.IPCIDRRange &&
		subnet.VMs.Equals(other.VMs)
}
