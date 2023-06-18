package identifier

import "github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"

type Subnetwork struct {
	Subnetwork string `json:"subnetwork"`
	Provider   string `json:"provider"`
	VPC        string `json:"vpc"`
	Network    string `json:"network"`
}

func (id *Subnetwork) Equals(other ID) bool {
	otherSubnetworkID, ok := other.(*Subnetwork)
	if !ok {
		return false
	}
	return id.Subnetwork == otherSubnetworkID.Subnetwork &&
		id.Provider == otherSubnetworkID.Provider &&
		id.VPC == otherSubnetworkID.VPC &&
		id.Network == otherSubnetworkID.Network
}

func (id *Subnetwork) ToIDLabels() map[string]string {
	return map[string]string{
		SubnetworkIdentifierKey: id.Subnetwork,
		ProviderIdentifierKey:   id.Provider,
		VpcIdentifierKey:        id.VPC,
		NetworkIdentifierKey:    id.Network,
	}
}

func (id *Subnetwork) IDFromLabels(labels map[string]string) errors.Error {
	subnetworkID, ok := labels[SubnetworkIdentifierKey]
	if !ok {
		return errors.InternalError.WithMessage("missing subnetwork identifier")
	}
	providerID, ok := labels[ProviderIdentifierKey]
	if !ok {
		return errors.InternalError.WithMessage("missing provider identifier")
	}
	vpcID := labels[VpcIdentifierKey]
	networkID, ok := labels[NetworkIdentifierKey]
	if !ok {
		return errors.InternalError.WithMessage("missing network identifier")
	}
	*id = Subnetwork{
		Subnetwork: subnetworkID,
		Provider:   providerID,
		VPC:        vpcID,
		Network:    networkID,
	}
	return errors.OK
}

func (id *Subnetwork) ToNameLabels() map[string]string {
	return map[string]string{
		SubnetworkNameKey: id.Subnetwork,
		ProviderNameKey:   id.Provider,
		VpcNameKey:        id.VPC,
		NetworkNameKey:    id.Network,
	}
}

func (id *Subnetwork) NameFromLabels(labels map[string]string) errors.Error {
	subnetworkName, ok := labels[SubnetworkNameKey]
	if !ok {
		return errors.InternalError.WithMessage("missing subnetwork name")
	}
	providerName, ok := labels[ProviderNameKey]
	if !ok {
		return errors.InternalError.WithMessage("missing provider name")
	}
	vpcName, ok := labels[VpcNameKey]
	networkName, ok := labels[NetworkNameKey]
	if !ok {
		return errors.InternalError.WithMessage("missing network name")
	}
	*id = Subnetwork{
		Subnetwork: subnetworkName,
		Provider:   providerName,
		VPC:        vpcName,
		Network:    networkName,
	}
	return errors.OK
}
