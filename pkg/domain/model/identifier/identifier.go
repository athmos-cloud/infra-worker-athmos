package identifier

import (
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
