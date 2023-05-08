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

func NewNetwork(payload NewResourcePayload) Network {
	payload.Validate()
	var id identifier.Network
	if reflect.TypeOf(payload.ParentIdentifier) == reflect.TypeOf(identifier.Provider{}) {
		parentID := payload.ParentIdentifier.(identifier.Provider)
		id = identifier.Network{
			ProviderID: parentID.ProviderID,
			NetworkID:  fmt.Sprintf("%s-%s", payload.Name, utils.RandomString(resourceIDSuffixLength)),
		}
	} else if reflect.TypeOf(payload.ParentIdentifier) == reflect.TypeOf(identifier.VPC{}) {
		parentID := payload.ParentIdentifier.(identifier.VPC)
		id = identifier.Network{
			ProviderID: parentID.ProviderID,
			VPCID:      parentID.VPCID,
			NetworkID:  fmt.Sprintf("%s-%s", payload.Name, utils.RandomString(resourceIDSuffixLength)),
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
		Status:      status.New(id.NetworkID, types.Network, payload.Provider),
		Subnetworks: map[string]Subnetwork{},
		Firewalls:   map[string]Firewall{},
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

type Network struct {
	mgm.DefaultModel `bson:",inline"`
	Metadata         metadata.Metadata     `bson:"metadata"`
	Identifier       identifier.Network    `bson:"identifier"`
	Status           status.ResourceStatus `bson:"status"`
	Subnetworks      SubnetworkCollection  `bson:"subnetworks,omitempty"`
	Firewalls        FirewallCollection    `bson:"firewalls,omitempty"`
}

func (network *Network) GetIdentifier() identifier.ID {
	return network.Identifier
}

func (network *Network) New(payload NewResourcePayload) IResource {
	id := payload.ParentIdentifier
	if reflect.TypeOf(id) != reflect.TypeOf(identifier.Network{}) {
		panic(errors.InvalidArgument.WithMessage("id type is not NetworkID"))
	}
	res := NewNetwork(payload)
	return &res
}

func (network *Network) GetMetadata() metadata.Metadata {
	return network.Metadata
}

func (network *Network) SetMetadata(request metadata.CreateMetadataRequest) {
	network.Metadata = metadata.New(request)
}

func (network *Network) SetStatus(resourceStatus status.ResourceStatus) {
	network.Status = resourceStatus
}

func (network *Network) GetStatus() status.ResourceStatus {
	return network.Status
}

func (network *Network) GetPluginReference() resourcePlugin.Reference {
	if !network.Status.PluginReference.ChartReference.Empty() {
		return network.Status.PluginReference
	}
	switch network.Status.PluginReference.ResourceReference.ProviderType {
	case types.GCP:
		network.Status.PluginReference.ChartReference = resourcePlugin.HelmChartReference{
			ChartName:    config.Current.Plugins.Crossplane.GCP.Network.Chart,
			ChartVersion: config.Current.Plugins.Crossplane.GCP.Network.Version,
		}
		return network.Status.PluginReference
	}
	panic(errors.InvalidArgument.WithMessage(fmt.Sprintf("network type %s not supported", network.Status.PluginReference.ResourceReference.ProviderType)))
}

func (network *Network) FromMap(data map[string]interface{}) {
	if err := resourcePlugin.InjectMapIntoStruct(data, network); !err.IsOk() {
		panic(err)
	}
}

func (network *Network) Insert(resource IResource, update ...bool) {
	shouldUpdate := false
	if len(update) > 0 {
		shouldUpdate = update[0]
	}
	idPayload := identifier.IDToPayload(resource.GetIdentifier())

	if idPayload.FirewallID != "" {
		if _, ok := network.Firewalls[idPayload.FirewallID]; ok && !shouldUpdate {
			panic(errors.Conflict.WithMessage(fmt.Sprintf("firewall %s already exists in network %s", idPayload.FirewallID, network.Identifier.NetworkID)))
		} else if !ok && shouldUpdate {
			panic(errors.NotFound.WithMessage(fmt.Sprintf("firewall %s not found in network %s", idPayload.FirewallID, network.Identifier.NetworkID)))
		}
		firewall := resource.(*Firewall)
		if network.Firewalls == nil {
			network.Firewalls = map[string]Firewall{
				idPayload.FirewallID: *firewall,
			}
		} else {
			network.Firewalls[idPayload.FirewallID] = *firewall
		}
		return
	} else if reflect.TypeOf(resource) == reflect.TypeOf(&Subnetwork{}) {
		if _, ok := network.Subnetworks[idPayload.SubnetID]; ok && !shouldUpdate {
			panic(errors.Conflict.WithMessage(fmt.Sprintf("firewall %s already exists in network %s", idPayload.FirewallID, network.Identifier.NetworkID)))
		} else if !ok && shouldUpdate {
			panic(errors.NotFound.WithMessage(fmt.Sprintf("firewall %s not found in network %s", idPayload.FirewallID, network.Identifier.NetworkID)))
		}
		subnet := resource.(*Subnetwork)
		if network.Subnetworks == nil {
			network.Subnetworks = map[string]Subnetwork{
				idPayload.SubnetID: *subnet,
			}
		} else {
			network.Subnetworks[idPayload.SubnetID] = *subnet
		}
		return
	} else if reflect.TypeOf(resource) != reflect.TypeOf(&Subnetwork{}) {
		network.insertInSubnetwork(resource, shouldUpdate)
		return
	}
	panic(errors.InvalidArgument.WithMessage(fmt.Sprintf("resource %v not supported for network insertion", resource)))
}

func (network *Network) insertInSubnetwork(resource IResource, update ...bool) {
	shouldUpdate := false
	if len(update) > 0 {
		shouldUpdate = update[0]
	}
	idPayload := identifier.IDToPayload(resource.GetIdentifier())
	_, ok := network.Subnetworks[idPayload.NetworkID]
	if reflect.TypeOf(resource) == reflect.TypeOf(&Subnetwork{}) && !ok && shouldUpdate {
		panic(errors.NotFound.WithMessage(fmt.Sprintf("subnetwork %s not found in network %s", idPayload.SubnetID, network.Identifier.NetworkID)))
	}
	if reflect.TypeOf(resource) == reflect.TypeOf(&Subnetwork{}) && ok && !shouldUpdate {
		panic(errors.Conflict.WithMessage(fmt.Sprintf("subnetwork %s already exists in network %s", idPayload.SubnetID, network.Identifier.NetworkID)))
	}
	if reflect.TypeOf(resource) == reflect.TypeOf(&Subnetwork{}) && (ok && shouldUpdate || !ok && !shouldUpdate) {
		subnet := resource.(*Subnetwork)
		network.Subnetworks[idPayload.NetworkID] = *subnet
		return
	} else if reflect.TypeOf(resource) != reflect.TypeOf(&Subnetwork{}) && idPayload.NetworkID != "" {
		subnet := network.Subnetworks[idPayload.NetworkID]
		subnet.Insert(resource, update...)
		return
	}
	panic(errors.InternalError.WithMessage(fmt.Sprintf("Invalid network insertion %v", resource)))
}

func (network *Network) Remove(resource IResource) {
	idPayload := identifier.IDToPayload(resource.GetIdentifier())
	if idPayload.FirewallID != "" {
		delete(network.Firewalls, idPayload.FirewallID)
	} else {
		subnet := network.Subnetworks[idPayload.NetworkID]
		subnet.Remove(resource)
	}
}

func (network *Network) Equals(other Network) bool {
	return network.Metadata.Equals(other.Metadata) &&
		network.Identifier.Equals(other.Identifier) &&
		network.Status.Equals(other.Status) &&
		network.Subnetworks.Equals(other.Subnetworks) &&
		network.Firewalls.Equals(other.Firewalls)
}
