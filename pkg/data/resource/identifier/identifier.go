package identifier

import (
	"encoding/json"
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/kamva/mgm/v3"
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

func BuildFromMap(payload map[string]interface{}) ID {
	body, err := json.Marshal(payload)
	if err != nil {
		panic(errors.InternalError.WithMessage(fmt.Sprintf("Error while marshalling payload: %s", err.Error())))
	}
	var idPayload IdPayload
	if errUnmarshal := json.Unmarshal(body, &idPayload); errUnmarshal != nil {
		panic(errors.InvalidArgument.WithMessage(fmt.Sprintf("%v is not a valid ID payload", payload)))
	}
	return Build(idPayload)
}

func Build(payload IdPayload) ID {
	if payload.VMID != "" && payload.ProviderID != "" && payload.VPCID != "" && payload.NetworkID != "" && payload.SubnetID != "" {
		return VM{
			ID:         payload.VMID,
			ProviderID: payload.ProviderID,
			VPCID:      payload.VPCID,
			NetworkID:  payload.NetworkID,
			SubnetID:   payload.SubnetID,
		}
	} else if payload.FirewallID != "" && payload.ProviderID != "" && payload.VPCID != "" && payload.NetworkID != "" {
		return Firewall{
			ID:         payload.FirewallID,
			ProviderID: payload.ProviderID,
			VPCID:      payload.VPCID,
			NetworkID:  payload.NetworkID,
		}
	} else if payload.ProviderID != "" && payload.VPCID != "" && payload.NetworkID != "" && payload.SubnetID != "" {
		return Subnetwork{
			ID:         payload.SubnetID,
			ProviderID: payload.ProviderID,
			VPCID:      payload.VPCID,
			NetworkID:  payload.NetworkID,
		}
	} else if payload.ProviderID != "" && payload.VPCID != "" && payload.NetworkID != "" {
		return Network{
			ID:         payload.NetworkID,
			ProviderID: payload.ProviderID,
			VPCID:      payload.VPCID,
		}
	} else if payload.ProviderID != "" && payload.VPCID != "" {
		return VPC{
			ID:         payload.VPCID,
			ProviderID: payload.ProviderID,
		}
	} else if payload.ProviderID != "" {
		return Provider{
			ID: payload.ProviderID,
		}
	} else {
		panic(errors.InvalidArgument.WithMessage(fmt.Sprintf("Invalid ID payload: %v", payload)))
	}
}

type Provider struct {
	mgm.DefaultModel `bson:",inline"`
	ID               string `bson:"id" plugin:"name"`
}

func (provider Provider) Equals(other ID) bool {
	otherProviderID, ok := other.(Provider)
	if !ok {
		panic(errors.InvalidArgument.WithMessage(fmt.Sprintf("%v is not a Provider ID", other)))
	}
	return provider.ID == otherProviderID.ID
}

type VPC struct {
	mgm.DefaultModel `bson:",inline"`
	ID               string `bson:"id" json:"id"`
	ProviderID       string `bson:"providerId" json:"providerID"`
}

func (id VPC) Equals(other ID) bool {
	otherVPCID, ok := other.(VPC)
	if !ok {
		panic(errors.InvalidArgument.WithMessage(fmt.Sprintf("%v is not a VPC VMID", other)))
	}
	return id.ID == otherVPCID.ID &&
		id.ProviderID == otherVPCID.ProviderID
}

type Network struct {
	mgm.DefaultModel `bson:",inline"`
	ID               string `bson:"id" json:"id" plugin:"networkID"`
	ProviderID       string `bson:"providerId" json:"providerID" plugin:"providerID"`
	VPCID            string `bson:"vpcId" json:"vpcID" plugin:"vpcID"`
}

func (id Network) Equals(other ID) bool {
	otherNetworkID, ok := other.(Network)
	if !ok {
		panic(errors.InvalidArgument.WithMessage(fmt.Sprintf("%v is not a Network VMID", other)))
	}
	return id.ID == otherNetworkID.ID &&
		id.ProviderID == otherNetworkID.ProviderID &&
		id.VPCID == otherNetworkID.VPCID
}

type Subnetwork struct {
	mgm.DefaultModel `bson:",inline"`
	ID               string `bson:"id" json:"subnetID" plugin:"subnetID"`
	ProviderID       string `bson:"providerId" json:"providerID" plugin:"providerID"`
	VPCID            string `bson:"vpcId" json:"vpcID" plugin:"vpcID"`
	NetworkID        string `bson:"networkId" json:"networkID" plugin:"networkID"`
}

func (id Subnetwork) Equals(other ID) bool {
	otherSubneworkID, ok := other.(Subnetwork)
	if !ok {
		panic(errors.InvalidArgument.WithMessage(fmt.Sprintf("%v is not a Subnetwork ID", other)))
	}
	return id.ID == otherSubneworkID.ID &&
		id.ProviderID == otherSubneworkID.ProviderID &&
		id.VPCID == otherSubneworkID.VPCID &&
		id.NetworkID == otherSubneworkID.NetworkID
}

type Firewall struct {
	mgm.DefaultModel `bson:",inline"`
	ID               string `bson:"id" json:"id" plugin:"firewallID"`
	ProviderID       string `bson:"providerId" json:"providerID" plugin:"providerID"`
	VPCID            string `bson:"vpcId" json:"vpcID" plugin:"vpcID"`
	NetworkID        string `bson:"networkId" json:"networkID" plugin:"networkID"`
}

func (id Firewall) Equals(other ID) bool {
	otherFirewallID, ok := other.(Firewall)
	if !ok {
		panic(errors.InvalidArgument.WithMessage(fmt.Sprintf("%v is not a Firewall ID", other)))
	}
	return id.ID == otherFirewallID.ID &&
		id.ProviderID == otherFirewallID.ProviderID &&
		id.VPCID == otherFirewallID.VPCID &&
		id.NetworkID == otherFirewallID.NetworkID
}

type VM struct {
	mgm.DefaultModel `bson:",inline"`
	ID               string `bson:"id" json:"vmID" plugin:"vmID"`
	ProviderID       string `bson:"providerId" json:"providerID" plugin:"providerID"`
	VPCID            string `bson:"vpcId" json:"vpcID" plugin:"vpcID"`
	NetworkID        string `bson:"networkId" json:"networkID" plugin:"networkID"`
	SubnetID         string `bson:"subnetId" json:"subnetID" plugin:"subnetID"`
}

func (id VM) Equals(other ID) bool {
	otherVMID, ok := other.(VM)
	if !ok {
		panic(errors.InvalidArgument.WithMessage(fmt.Sprintf("%v is not a VM ID", other)))
	}
	return id.ID == otherVMID.ID &&
		id.ProviderID == otherVMID.ProviderID &&
		id.VPCID == otherVMID.VPCID &&
		id.NetworkID == otherVMID.NetworkID &&
		id.SubnetID == otherVMID.SubnetID
}
