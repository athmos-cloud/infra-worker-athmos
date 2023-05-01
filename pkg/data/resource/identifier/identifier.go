package identifier

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
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

func NewID(payload IdPayload) ID {
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
	ID string `bson:"id"`
}

func (provider Provider) Equals(other ID) bool {
	otherProviderID, ok := other.(Provider)
	if !ok {
		panic(errors.InvalidArgument.WithMessage(fmt.Sprintf("%v is not a Provider ID", other)))
	}
	return provider.ID == otherProviderID.ID
}

type VPC struct {
	ID         string `bson:"id" json:"id"`
	ProviderID string `bson:"providerId" json:"providerID"`
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
	ID         string `bson:"id" json:"id"`
	ProviderID string `bson:"providerId" json:"providerID"`
	VPCID      string `bson:"vpcId" json:"vpcID"`
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
	ID         string `bson:"id" json:"id"`
	ProviderID string `bson:"providerId" json:"providerID"`
	VPCID      string `bson:"vpcId" json:"vpcID"`
	NetworkID  string `bson:"networkId" json:"networkID"`
}

func (id Subnetwork) Equals(other ID) bool {
	otherSubneworkID, ok := other.(Subnetwork)
	if !ok {
		panic(errors.InvalidArgument.WithMessage(fmt.Sprintf("%v is not a Subentework VMID", other)))
	}
	return id.ID == otherSubneworkID.ID &&
		id.ProviderID == otherSubneworkID.ProviderID &&
		id.VPCID == otherSubneworkID.VPCID &&
		id.NetworkID == otherSubneworkID.NetworkID
}

type Firewall struct {
	ID         string `bson:"id" json:"id"`
	ProviderID string `bson:"providerId" json:"providerID"`
	VPCID      string `bson:"vpcId" json:"vpcID"`
	NetworkID  string `bson:"networkId" json:"networkID"`
}

func (id Firewall) Equals(other ID) bool {
	otherFirewallID, ok := other.(Firewall)
	if !ok {
		panic(errors.InvalidArgument.WithMessage(fmt.Sprintf("%v is not a Firewall VMID", other)))
	}
	return id.ID == otherFirewallID.ID &&
		id.ProviderID == otherFirewallID.ProviderID &&
		id.VPCID == otherFirewallID.VPCID &&
		id.NetworkID == otherFirewallID.NetworkID
}

type VM struct {
	ID         string `bson:"id" json:"vmID"`
	ProviderID string `bson:"providerId" json:"providerID"`
	VPCID      string `bson:"vpcId" json:"vpcID"`
	NetworkID  string `bson:"networkId" json:"networkID"`
	SubnetID   string `bson:"subnetId" json:"subnetID"`
}

func (id VM) Equals(other ID) bool {
	otherVMID, ok := other.(VM)
	if !ok {
		panic(errors.InvalidArgument.WithMessage(fmt.Sprintf("%v is not a VM VMID", other)))
	}
	return id.ID == otherVMID.ID &&
		id.ProviderID == otherVMID.ProviderID &&
		id.VPCID == otherVMID.VPCID &&
		id.NetworkID == otherVMID.NetworkID &&
		id.SubnetID == otherVMID.SubnetID
}
