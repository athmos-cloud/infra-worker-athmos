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

func (id *Firewall) ToLabels() map[string]string {
	return map[string]string{
		FirewallIdentifierKey: id.Firewall,
		ProviderIdentifierKey: id.Provider,
		VpcIdentifierKey:      id.VPC,
		NetworkIdentifierKey:  id.Network,
	}
}

func (id *Firewall) FromLabels(labels map[string]string) errors.Error {
	firewallID, ok := labels[FirewallIdentifierKey]
	if !ok {
		return errors.InternalError.WithMessage(fmt.Sprintf("missing firewall identifier for firewall ID: %v", labels))
	}
	providerID, ok := labels[ProviderIdentifierKey]
	if !ok {
		return errors.InternalError.WithMessage(fmt.Sprintf("missing provider identifier for firewall ID: %v", labels))
	}
	vpcID, ok := labels[VpcIdentifierKey]
	if !ok {
		return errors.InternalError.WithMessage(fmt.Sprintf("missing vpc identifier for firewall ID: %v", labels))
	}
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
