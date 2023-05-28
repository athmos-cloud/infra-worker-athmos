package identifier

import "github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"

type Network struct {
	Network  string `json:"id"`
	Provider string `json:"provider"`
	VPC      string `json:"vpc"`
}

func (id *Network) NameFromLabels(labels map[string]string) errors.Error {
	network, ok := labels[NetworkNameKey]
	if !ok {
		return errors.InternalError.WithMessage("missing network identifier")
	}
	provider, ok := labels[ProviderIdentifierKey]
	if !ok {
		return errors.InternalError.WithMessage("missing provider identifier")
	}
	vpc := labels[VpcIdentifierKey]
	*id = Network{
		Network:  network,
		Provider: provider,
		VPC:      vpc,
	}
	return errors.OK
}

func (id *Network) ToNameLabels() map[string]string {
	return map[string]string{
		NetworkNameKey:  id.Network,
		ProviderNameKey: id.Provider,
		VpcNameKey:      id.VPC,
	}
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

func (id *Network) ToIDLabels() map[string]string {
	return map[string]string{
		NetworkIdentifierKey:  id.Network,
		ProviderIdentifierKey: id.Provider,
		VpcIdentifierKey:      id.VPC,
	}
}

func (id *Network) IDFromLabels(labels map[string]string) errors.Error {
	networkID, ok := labels[NetworkIdentifierKey]
	if !ok {
		return errors.InternalError.WithMessage("missing network identifier")
	}
	providerID, ok := labels[ProviderIdentifierKey]
	if !ok {
		return errors.InternalError.WithMessage("missing provider identifier")
	}
	vpcID, ok := labels[VpcIdentifierKey]
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
