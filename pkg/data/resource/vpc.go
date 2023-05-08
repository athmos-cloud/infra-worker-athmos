package resource

import (
	"fmt"
	resourcePlugin "github.com/athmos-cloud/infra-worker-athmos/pkg/data/plugin"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/status"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/config"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/utils"
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
		VPCID:      fmt.Sprintf("%s-%s", payload.Name, utils.RandomString(resourceIDSuffixLength)),
	}
	return VPC{
		Metadata: metadata.New(metadata.CreateMetadataRequest{
			Name:         id.VPCID,
			NotMonitored: !payload.Managed,
			Tags:         payload.Tags,
		}),
		Identifier: id,
		Status:     status.New(id.VPCID, types.VPC, payload.Provider),
		Networks:   make(NetworkCollection),
	}
}

type VPC struct {
	mgm.DefaultModel `bson:",inline"`
	Metadata         metadata.Metadata     `bson:"metadata"`
	Identifier       identifier.VPC        `bson:"identifier"`
	Status           status.ResourceStatus `bson:"status"`
	Networks         NetworkCollection     `bson:"networks,omitempty"`
}

func (vpc *VPC) GetIdentifier() identifier.ID {
	return vpc.Identifier
}

func (vpc *VPC) New(payload NewResourcePayload) IResource {
	if reflect.TypeOf(payload.ParentIdentifier) != reflect.TypeOf(identifier.VPC{}) {
		panic(errors.InvalidArgument.WithMessage("invalid id type"))
	}
	res := NewVPC(payload)
	return &res
}

func (vpc *VPC) SetMetadata(request metadata.CreateMetadataRequest) {
	vpc.Metadata = metadata.New(request)
}

func (vpc *VPC) GetMetadata() metadata.Metadata {
	return vpc.Metadata
}

func (vpc *VPC) SetStatus(resourceStatus status.ResourceStatus) {
	vpc.Status = resourceStatus
}

func (vpc *VPC) GetStatus() status.ResourceStatus {
	return vpc.Status
}

func (vpc *VPC) GetPluginReference() resourcePlugin.Reference {
	if !vpc.Status.PluginReference.ChartReference.Empty() {
		return vpc.Status.PluginReference
	}
	switch vpc.Status.PluginReference.ResourceReference.ProviderType {
	case types.GCP:
		vpc.Status.PluginReference.ChartReference = resourcePlugin.HelmChartReference{
			ChartName:    config.Current.Plugins.Crossplane.GCP.VPC.Chart,
			ChartVersion: config.Current.Plugins.Crossplane.GCP.VPC.Version,
		}
		return vpc.Status.PluginReference
	}
	panic(errors.InvalidArgument.WithMessage(fmt.Sprintf("provider type %s not supported", vpc.Status.PluginReference.ResourceReference.ProviderType)))
}

func (vpc *VPC) FromMap(data map[string]interface{}) {
	err := resourcePlugin.InjectMapIntoStruct(data, vpc)
	if !err.IsOk() {
		panic(err)
	}
}

func (vpc *VPC) Insert(resource IResource, update ...bool) {
	shouldUpdate := false
	if len(update) > 0 {
		shouldUpdate = update[0]
	}
	idPayload := identifier.IDToPayload(resource.GetIdentifier())
	_, ok := vpc.Networks[idPayload.NetworkID]
	if reflect.TypeOf(resource) == reflect.TypeOf(&Network{}) && !ok && shouldUpdate {
		panic(errors.NotFound.WithMessage(fmt.Sprintf("network %s not found in vpc %s", idPayload.NetworkID, vpc.Identifier.VPCID)))
	}
	if reflect.TypeOf(resource) == reflect.TypeOf(&Network{}) && ok && !shouldUpdate {
		panic(errors.Conflict.WithMessage(fmt.Sprintf("network %s already exists in vpc %s", idPayload.NetworkID, vpc.Identifier.VPCID)))
	}
	if reflect.TypeOf(resource) == reflect.TypeOf(&Network{}) && (ok && shouldUpdate || !ok && !shouldUpdate) {
		network := resource.(*Network)
		if vpc.Networks == nil {
			vpc.Networks = map[string]Network{
				idPayload.NetworkID: *network,
			}
		} else {
			vpc.Networks[idPayload.NetworkID] = *network
		}
		return
	} else if reflect.TypeOf(resource) != reflect.TypeOf(&Network{}) && idPayload.NetworkID != "" {
		network := vpc.Networks[idPayload.NetworkID]
		network.Insert(resource, update...)
		return
	}
	panic(errors.InternalError.WithMessage(fmt.Sprintf("Invalid vpc insertion %v", resource)))
}

func (vpc *VPC) Remove(resource IResource) {
	idPayload := identifier.IDToPayload(resource.GetIdentifier())
	if reflect.TypeOf(resource) == reflect.TypeOf(&Network{}) {
		delete(vpc.Networks, idPayload.NetworkID)
	} else {
		network := vpc.Networks[idPayload.NetworkID]
		network.Remove(resource)
	}
}

func (vpc *VPC) Equals(other VPC) bool {
	return vpc.Metadata.Equals(other.Metadata) &&
		vpc.Identifier.Equals(other.Identifier) &&
		vpc.Status.Equals(other.Status) &&
		vpc.Networks.Equals(other.Networks)
}
