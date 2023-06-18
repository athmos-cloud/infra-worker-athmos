package identifier

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

type Firewall struct {
	Firewall string `json:"firewall"`
	Provider string `json:"provider"`
	VPC      string `json:"vpc"`
	Network  string `json:"network"`
}

func (id *Firewall) Equals(other ID) bool {
	otherFirewallID, ok := other.(*Firewall)
	if !ok {
		return false
	}
	return id.Firewall == otherFirewallID.Firewall &&
		id.Provider == otherFirewallID.Provider &&
		id.VPC == otherFirewallID.VPC &&
		id.Network == otherFirewallID.Network
}

func (id *Firewall) ToIDLabels() map[string]string {
	return map[string]string{
		FirewallIdentifierKey: id.Firewall,
		ProviderIdentifierKey: id.Provider,
		VpcIdentifierKey:      id.VPC,
		NetworkIdentifierKey:  id.Network,
	}
}

func (id *Firewall) IDFromLabels(labels map[string]string) errors.Error {
	firewallID, ok := labels[FirewallIdentifierKey]
	if !ok {
		return errors.InternalError.WithMessage(fmt.Sprintf("missing firewall identifier for firewall ID: %v", labels))
	}
	providerID, ok := labels[ProviderIdentifierKey]
	if !ok {
		return errors.InternalError.WithMessage(fmt.Sprintf("missing provider identifier for firewall ID: %v", labels))
	}
	vpcID := labels[VpcIdentifierKey]
	networkID, ok := labels[NetworkIdentifierKey]
	if !ok {
		return errors.InternalError.WithMessage(fmt.Sprintf("missing network identifier for firewall ID: %v", labels))
	}
	*id = Firewall{
		Firewall: firewallID,
		Provider: providerID,
		VPC:      vpcID,
		Network:  networkID,
	}
	return errors.OK
}

func (id *Firewall) ToNameLabels() map[string]string {
	return map[string]string{
		FirewallNameKey: id.Firewall,
		ProviderNameKey: id.Provider,
		VpcNameKey:      id.VPC,
		NetworkNameKey:  id.Network,
	}
}

func (id *Firewall) NameFromLabels(labels map[string]string) errors.Error {
	firewall, ok := labels[FirewallNameKey]
	if !ok {
		return errors.InternalError.WithMessage(fmt.Sprintf("missing firewall identifier for firewall ID: %v", labels))
	}
	provider, ok := labels[ProviderNameKey]
	if !ok {
		return errors.InternalError.WithMessage(fmt.Sprintf("missing provider identifier for firewall ID: %v", labels))
	}
	vpc := labels[VpcNameKey]
	network, ok := labels[NetworkNameKey]
	if !ok {
		return errors.InternalError.WithMessage(fmt.Sprintf("missing network identifier for firewall ID: %v", labels))
	}
	*id = Firewall{
		Firewall: firewall,
		Provider: provider,
		VPC:      vpc,
		Network:  network,
	}
	return errors.OK
}
