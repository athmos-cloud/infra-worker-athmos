package identifier

import (
	"encoding/json"
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/kamva/mgm/v3"
	"reflect"
)

const (
	IdentifierKey = "identifier"
)

type ID interface {
	Equals(other ID) bool
}

type IdPayload struct {
	ProviderID string `json:"providerID"`
	VPCID      string `json:"vpcID"`
	NetworkID  string `json:"networkID"`
	SubnetID   string `json:"subnetID"`
	VMID       string `json:"vmID"`
	FirewallID string `json:"firewallID"`
}

func IDToPayload(input ID) IdPayload {
	switch reflect.TypeOf(input) {
	case reflect.TypeOf(Provider{}):
		id := input.(Provider)
		return IdPayload{
			ProviderID: id.ProviderID,
		}
	case reflect.TypeOf(VPC{}):
		id := input.(VPC)
		return IdPayload{
			ProviderID: id.ProviderID,
			VPCID:      id.VPCID,
		}
	case reflect.TypeOf(Network{}):
		id := input.(Network)
		return IdPayload{
			ProviderID: id.ProviderID,
			VPCID:      id.VPCID,
			NetworkID:  id.NetworkID,
		}
	case reflect.TypeOf(Subnetwork{}):
		id := input.(Subnetwork)
		return IdPayload{
			ProviderID: id.ProviderID,
			VPCID:      id.VPCID,
			NetworkID:  id.NetworkID,
			SubnetID:   id.SubnetworkID,
		}
	case reflect.TypeOf(Firewall{}):
		id := input.(Firewall)
		return IdPayload{
			ProviderID: id.ProviderID,
			VPCID:      id.VPCID,
			NetworkID:  id.NetworkID,
			FirewallID: id.FirewallID,
		}
	case reflect.TypeOf(VM{}):
		id := input.(VM)
		return IdPayload{
			ProviderID: id.ProviderID,
			VPCID:      id.VPCID,
			NetworkID:  id.NetworkID,
			SubnetID:   id.SubnetID,
			VMID:       id.VMID,
		}
	}
	return IdPayload{}
}

func BuildFromMap(payload map[string]interface{}) ID {
	body, err := json.Marshal(payload)
	if err != nil {
		panic(errors.InternalError.WithMessage(fmt.Sprintf("Error while marshalling payload: %s", err.Error())))
	}
	var idPayload IdPayload
	if errUnmarshal := json.Unmarshal(body, &idPayload); errUnmarshal != nil {
		panic(errors.InvalidArgument.WithMessage(fmt.Sprintf("%v is not a valid FirewallID payload", payload)))
	}
	return Build(idPayload)
}

func Build(payload IdPayload) ID {
	if payload.VMID != "" && payload.ProviderID != "" && payload.VPCID != "" && payload.NetworkID != "" && payload.SubnetID != "" {
		return VM{
			VMID:       payload.VMID,
			ProviderID: payload.ProviderID,
			VPCID:      payload.VPCID,
			NetworkID:  payload.NetworkID,
			SubnetID:   payload.SubnetID,
		}
	} else if payload.FirewallID != "" && payload.ProviderID != "" && payload.VPCID != "" && payload.NetworkID != "" {
		return Firewall{
			FirewallID: payload.FirewallID,
			ProviderID: payload.ProviderID,
			VPCID:      payload.VPCID,
			NetworkID:  payload.NetworkID,
		}
	} else if payload.ProviderID != "" && payload.NetworkID != "" && payload.SubnetID != "" {
		return Subnetwork{
			SubnetworkID: payload.SubnetID,
			ProviderID:   payload.ProviderID,
			VPCID:        payload.VPCID,
			NetworkID:    payload.NetworkID,
		}
	} else if payload.ProviderID != "" && payload.NetworkID != "" {
		return Network{
			NetworkID:  payload.NetworkID,
			ProviderID: payload.ProviderID,
			VPCID:      payload.VPCID,
		}
	} else if payload.ProviderID != "" && payload.VPCID != "" {
		return VPC{
			VPCID:      payload.VPCID,
			ProviderID: payload.ProviderID,
		}
	} else if payload.ProviderID != "" {
		return Provider{
			ProviderID: payload.ProviderID,
		}
	} else {
		return Empty{}
	}
}

type Empty struct{}

func (e Empty) Equals(other ID) bool {
	_, ok := other.(Empty)
	return ok
}

type Provider struct {
	mgm.DefaultModel `bson:",inline"`
	ProviderID       string `bson:"id" plugin:"id"`
}

func (provider Provider) Equals(other ID) bool {
	otherProviderID, ok := other.(Provider)
	if !ok {
		panic(errors.InvalidArgument.WithMessage(fmt.Sprintf("%v is not a Provider FirewallID", other)))
	}
	return provider.ProviderID == otherProviderID.ProviderID
}

type VPC struct {
	mgm.DefaultModel `bson:",inline"`
	VPCID            string `bson:"id" json:"id"`
	ProviderID       string `bson:"providerId" json:"providerID"`
}

func (id VPC) Equals(other ID) bool {
	otherVPCID, ok := other.(VPC)
	if !ok {
		panic(errors.InvalidArgument.WithMessage(fmt.Sprintf("%v is not a VPC VMID", other)))
	}
	return id.VPCID == otherVPCID.VPCID &&
		id.ProviderID == otherVPCID.ProviderID
}

type Network struct {
	mgm.DefaultModel `bson:",inline"`
	NetworkID        string `bson:"id" json:"id" plugin:"networkID"`
	ProviderID       string `bson:"providerId" json:"providerID" plugin:"providerID"`
	VPCID            string `bson:"vpcId" json:"vpcID" plugin:"vpcID"`
}

func (id Network) Equals(other ID) bool {
	otherNetworkID, ok := other.(Network)
	if !ok {
		panic(errors.InvalidArgument.WithMessage(fmt.Sprintf("%v is not a Network VMID", other)))
	}
	return id.NetworkID == otherNetworkID.NetworkID &&
		id.ProviderID == otherNetworkID.ProviderID &&
		id.VPCID == otherNetworkID.VPCID
}

type Subnetwork struct {
	mgm.DefaultModel `bson:",inline"`
	SubnetworkID     string `bson:"id" json:"subnetID" plugin:"subnetID"`
	ProviderID       string `bson:"providerId" json:"providerID" plugin:"providerID"`
	VPCID            string `bson:"vpcId" json:"vpcID" plugin:"vpcID"`
	NetworkID        string `bson:"networkId" json:"networkID" plugin:"networkID"`
}

func (id Subnetwork) Equals(other ID) bool {
	otherSubneworkID, ok := other.(Subnetwork)
	if !ok {
		panic(errors.InvalidArgument.WithMessage(fmt.Sprintf("%v is not a Subnetwork FirewallID", other)))
	}
	return id.SubnetworkID == otherSubneworkID.SubnetworkID &&
		id.ProviderID == otherSubneworkID.ProviderID &&
		id.VPCID == otherSubneworkID.VPCID &&
		id.NetworkID == otherSubneworkID.NetworkID
}

type Firewall struct {
	mgm.DefaultModel `bson:",inline"`
	FirewallID       string `bson:"id" json:"id" plugin:"firewallID"`
	ProviderID       string `bson:"providerId" json:"providerID" plugin:"providerID"`
	VPCID            string `bson:"vpcId" json:"vpcID" plugin:"vpcID"`
	NetworkID        string `bson:"networkId" json:"networkID" plugin:"networkID"`
}

func (id Firewall) Equals(other ID) bool {
	otherFirewallID, ok := other.(Firewall)
	if !ok {
		panic(errors.InvalidArgument.WithMessage(fmt.Sprintf("%v is not a Firewall FirewallID", other)))
	}
	return id.FirewallID == otherFirewallID.FirewallID &&
		id.ProviderID == otherFirewallID.ProviderID &&
		id.VPCID == otherFirewallID.VPCID &&
		id.NetworkID == otherFirewallID.NetworkID
}

type VM struct {
	mgm.DefaultModel `bson:",inline"`
	VMID             string `bson:"id" json:"vmID" plugin:"vmID"`
	ProviderID       string `bson:"providerId" json:"providerID" plugin:"providerID"`
	VPCID            string `bson:"vpcId" json:"vpcID" plugin:"vpcID"`
	NetworkID        string `bson:"networkId" json:"networkID" plugin:"networkID"`
	SubnetID         string `bson:"subnetId" json:"subnetID" plugin:"subnetID"`
}

func (id VM) Equals(other ID) bool {
	otherVMID, ok := other.(VM)
	if !ok {
		panic(errors.InvalidArgument.WithMessage(fmt.Sprintf("%v is not a VM FirewallID", other)))
	}
	return id.VMID == otherVMID.VMID &&
		id.ProviderID == otherVMID.ProviderID &&
		id.VPCID == otherVMID.VPCID &&
		id.NetworkID == otherVMID.NetworkID &&
		id.SubnetID == otherVMID.SubnetID
}

func IDParentMatchesWithResource(idParent ID, resourceType types.ResourceType) bool {
	switch reflect.TypeOf(idParent) {
	case nil:
		return resourceType == types.Provider
	case reflect.TypeOf(Provider{}):
		return resourceType == types.VPC || resourceType == types.Network
	case reflect.TypeOf(VPC{}):
		return resourceType == types.Network
	case reflect.TypeOf(Network{}):
		return resourceType == types.Subnetwork || resourceType == types.Firewall
	case reflect.TypeOf(Subnetwork{}):
		return resourceType == types.VM
	default:
		return false
	}
}
