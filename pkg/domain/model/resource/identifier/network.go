package identifier

import "github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"

type Network struct {
	Network  string `json:"id"`
	Provider string `json:"provider"`
	VPC      string `json:"vpc"`
}

func (id *Network) Equals(other ID) bool {
	otherNetworkID, ok := other.(*Network)
	if !ok {
		return false
	}
	return id.Network == otherNetworkID.Network &&
		id.Provider == otherNetworkID.Provider &&
		id.VPC == otherNetworkID.VPC
}

func (id *Network) ToLabels() map[string]string {
	return map[string]string{
		networkIdentifierKey:  id.Network,
		providerIdentifierKey: id.Provider,
		vpcIdentifierKey:      id.VPC,
	}
}

func (id *Network) FromLabels(labels map[string]string) errors.Error {
	networkID, ok := labels[networkIdentifierKey]
	if !ok {
		return errors.InternalError.WithMessage("missing network identifier")
	}
	providerID, ok := labels[providerIdentifierKey]
	if !ok {
		return errors.InternalError.WithMessage("missing provider identifier")
	}
	vpcID, ok := labels[vpcIdentifierKey]
	if !ok {
		return errors.InternalError.WithMessage("missing vpc identifier")
	}
	*id = Network{
		Network:  networkID,
		Provider: providerID,
		VPC:      vpcID,
	}
	return errors.OK
}
