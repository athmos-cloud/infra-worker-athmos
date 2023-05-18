package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/kamva/mgm/v3"
	"reflect"
)

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

func NewVPC(payload NewResourcePayload) VPC {
	payload.Validate()
	if reflect.TypeOf(payload.ParentIdentifier) != reflect.TypeOf(identifier.Provider{}) {
		panic(errors.InvalidArgument.WithMessage("invalid id type"))
	}
	parentID := payload.ParentIdentifier.(identifier.Provider)
	id := identifier.VPC{
		ProviderID: parentID.ProviderID,
		VPCID:      formatResourceName(payload.Name),
	}
	return VPC{
		Metadata: metadata.New(metadata.CreateMetadataRequest{
			Name:         id.VPCID,
			NotMonitored: !payload.Managed,
			Tags:         payload.Tags,
		}),
		Identifier: id,
		Networks:   make(NetworkCollection),
	}
}

type VPC struct {
	mgm.DefaultModel `bson:",inline"`
	Metadata         metadata.Metadata `bson:"metadata"`
	Identifier       identifier.VPC    `bson:"identifier"`
	Networks         NetworkCollection `bson:"networks,omitempty"`
}

func (vpc *VPC) Equals(other VPC) bool {
	return vpc.Metadata.Equals(other.Metadata) &&
		vpc.Identifier.Equals(other.Identifier) &&
		vpc.Networks.Equals(other.Networks)
}
