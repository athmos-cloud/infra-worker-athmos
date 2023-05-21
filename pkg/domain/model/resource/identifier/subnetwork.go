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

func (id *Subnetwork) ToLabels() map[string]string {
	return map[string]string{
		subnetworkIdentifierKey: id.Subnetwork,
		providerIdentifierKey:   id.Provider,
		vpcIdentifierKey:        id.VPC,
		networkIdentifierKey:    id.Network,
	}
}

func (id *Subnetwork) FromLabels(labels map[string]string) errors.Error {
	subnetworkID, ok := labels[subnetworkIdentifierKey]
	if !ok {
		return errors.InternalError.WithMessage("missing subnetwork identifier")
	}
	providerID, ok := labels[providerIdentifierKey]
	if !ok {
		return errors.InternalError.WithMessage("missing provider identifier")
	}
	vpcID, ok := labels[vpcIdentifierKey]
	if !ok {
		return errors.InternalError.WithMessage("missing vpc identifier")
	}
	networkID, ok := labels[networkIdentifierKey]
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
