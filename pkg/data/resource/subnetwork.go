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

func NewSubnetwork(payload NewResourcePayload) Subnetwork {
	payload.Validate()
	if reflect.TypeOf(payload.ParentIdentifier) != reflect.TypeOf(identifier.Network{}) {
		panic(errors.InvalidArgument.WithMessage("ID type must be network ID"))
	}
	parentID := payload.ParentIdentifier.(identifier.Network)
	id := identifier.Subnetwork{
		ProviderID:   parentID.ProviderID,
		NetworkID:    parentID.NetworkID,
		VPCID:        parentID.VPCID,
		SubnetworkID: fmt.Sprintf("%s-%s", payload.Name, utils.RandomString(resourceIDSuffixLength)),
	}

	return Subnetwork{
		Metadata: metadata.New(metadata.CreateMetadataRequest{
			Name:         id.ProviderID,
			NotMonitored: !payload.Monitored,
			Tags:         payload.Tags,
		}),
		Identifier: id,
		Status:     status.New(id.SubnetworkID, types.Subnetwork, payload.Provider),
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

type Subnetwork struct {
	mgm.DefaultModel `bson:",inline"`
	Metadata         metadata.Metadata     `bson:"metadata"`
	Identifier       identifier.Subnetwork `bson:"hierarchyLocation"`
	Status           status.ResourceStatus `bson:"status"`
	Region           string                `bson:"region" plugin:"region"`
	IPCIDRRange      string                `bson:"ipCidrRange" plugin:"ipCidrRange"`
	VMs              VMCollection          `bson:"vmList,omitempty"`
}

func (subnet *Subnetwork) GetIdentifier() identifier.ID {
	return subnet.Identifier
}

func (subnet *Subnetwork) New(payload NewResourcePayload) IResource {
	if reflect.TypeOf(payload.ParentIdentifier) != reflect.TypeOf(identifier.Subnetwork{}) {
		panic(errors.InvalidArgument.WithMessage("id type is not SubnetworkID"))
	}
	res := NewSubnetwork(payload)
	return &res
}

func (subnet *Subnetwork) GetMetadata() metadata.Metadata {
	return subnet.Metadata
}

func (subnet *Subnetwork) SetMetadata(request metadata.CreateMetadataRequest) {
	subnet.Metadata = metadata.New(request)
}

func (subnet *Subnetwork) SetStatus(resourceStatus status.ResourceStatus) {
	subnet.Status = resourceStatus
}

func (subnet *Subnetwork) GetStatus() status.ResourceStatus {
	return subnet.Status
}

func (subnet *Subnetwork) GetPluginReference() resourcePlugin.Reference {
	if !subnet.Status.PluginReference.ChartReference.Empty() {
		return subnet.Status.PluginReference
	}
	switch subnet.Status.PluginReference.ResourceReference.ProviderType {
	case types.GCP:
		subnet.Status.PluginReference.ChartReference = resourcePlugin.HelmChartReference{
			ChartName:    config.Current.Plugins.Crossplane.GCP.Subnet.Chart,
			ChartVersion: config.Current.Plugins.Crossplane.GCP.Subnet.Version,
		}
		return subnet.Status.PluginReference
	}
	panic(errors.InvalidArgument.WithMessage(fmt.Sprintf("provider type %s not supported", subnet.Status.PluginReference.ResourceReference.ProviderType)))
}

func (subnet *Subnetwork) FromMap(data map[string]interface{}) {
	if err := resourcePlugin.InjectMapIntoStruct(data, subnet); !err.IsOk() {
		panic(err)
	}
}

func (subnet *Subnetwork) Insert(resource IResource, update ...bool) {
	shouldUpdate := false
	if len(update) > 0 {
		shouldUpdate = update[0]
	}
	idPayload := identifier.IDToPayload(resource.GetIdentifier())
	_, ok := subnet.VMs[idPayload.VMID]
	if reflect.TypeOf(resource) == reflect.TypeOf(&VM{}) && !ok && shouldUpdate {
		panic(errors.NotFound.WithMessage(fmt.Sprintf("vm %s not found in subnet %s", idPayload.NetworkID, subnet.Identifier.SubnetworkID)))
	}
	if reflect.TypeOf(resource) == reflect.TypeOf(&VM{}) && ok && !shouldUpdate {
		panic(errors.Conflict.WithMessage(fmt.Sprintf("vm %s already exists in subnet %s", idPayload.NetworkID, subnet.Identifier.VPCID)))
	}
	if reflect.TypeOf(resource) == reflect.TypeOf(&VM{}) && (ok && shouldUpdate || !ok && !shouldUpdate) {
		vm := resource.(*VM)
		subnet.VMs[idPayload.VMID] = *vm
		return
	}
	panic(errors.InternalError.WithMessage(fmt.Sprintf("Invalid subnet insertion %v", resource)))
}

func (subnet *Subnetwork) Remove(resource IResource) {
	if reflect.TypeOf(resource) != reflect.TypeOf(&VM{}) {
		return
	}
	delete(subnet.VMs, identifier.IDToPayload(resource.GetIdentifier()).VMID)
}

func (subnet *Subnetwork) Equals(other Subnetwork) bool {
	return subnet.Metadata.Equals(other.Metadata) &&
		subnet.Identifier.Equals(other.Identifier) &&
		subnet.Status.Equals(other.Status) &&
		subnet.Region == other.Region &&
		subnet.IPCIDRRange == other.IPCIDRRange &&
		subnet.VMs.Equals(other.VMs)
}
