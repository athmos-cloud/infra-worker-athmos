package identifier

import "github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"

type VM struct {
	VM         string `json:"vm"`
	Provider   string `json:"provider"`
	VPC        string `json:"vpc"`
	Network    string `json:"network"`
	Subnetwork string `json:"subnetwork"`
}

func (id *VM) Equals(other ID) bool {
	otherVMID, ok := other.(*VM)
	if !ok {
		return false
	}
	return id.VM == otherVMID.VM &&
		id.Provider == otherVMID.Provider &&
		id.VPC == otherVMID.VPC &&
		id.Network == otherVMID.Network &&
		id.Subnetwork == otherVMID.Subnetwork
}

func (id *VM) ToLabels() map[string]string {
	return map[string]string{
		VMIdentifierKey:         id.VM,
		ProviderIdentifierKey:   id.Provider,
		VpcIdentifierKey:        id.VPC,
		NetworkIdentifierKey:    id.Network,
		SubnetworkIdentifierKey: id.Subnetwork,
	}
}

func (id *VM) FromLabels(labels map[string]string) errors.Error {
	vmID, ok := labels[VMIdentifierKey]
	if !ok {
		return errors.InternalError.WithMessage("missing vm identifier")
	}
	providerID, ok := labels[ProviderIdentifierKey]
	if !ok {
		return errors.InternalError.WithMessage("missing provider identifier")
	}
	vpcID, ok := labels[VpcIdentifierKey]
	if !ok {
		return errors.InternalError.WithMessage("missing vpc identifier")
	}
	networkID, ok := labels[NetworkIdentifierKey]
	if !ok {
		return errors.InternalError.WithMessage("missing network identifier")
	}
	subnetworkID, ok := labels[SubnetworkIdentifierKey]
	if !ok {
		return errors.InternalError.WithMessage("missing subnetwork identifier")
	}
	*id = VM{
		VM:         vmID,
		Provider:   providerID,
		VPC:        vpcID,
		Network:    networkID,
		Subnetwork: subnetworkID,
	}
	return errors.OK
}
