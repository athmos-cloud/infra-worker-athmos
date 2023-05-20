package identifier

import "github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"

type VM struct {
	VM       string `json:"vm"`
	Provider string `json:"provider"`
	VPC      string `json:"vpc"`
	Network  string `json:"network"`
	Subnet   string `json:"subnet"`
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
		id.Subnet == otherVMID.Subnet
}

func (id *VM) ToLabels() map[string]string {
	return map[string]string{
		vmIdentifierKey:         id.VM,
		providerIdentifierKey:   id.Provider,
		vpcIdentifierKey:        id.VPC,
		networkIdentifierKey:    id.Network,
		subnetworkIdentifierKey: id.Subnet,
	}
}

func (id *VM) FromLabels(labels map[string]string) errors.Error {
	vmID, ok := labels[vmIdentifierKey]
	if !ok {
		return errors.InternalError.WithMessage("missing vm identifier")
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
	subnetworkID, ok := labels[subnetworkIdentifierKey]
	if !ok {
		return errors.InternalError.WithMessage("missing subnetwork identifier")
	}
	*id = VM{
		VM:       vmID,
		Provider: providerID,
		VPC:      vpcID,
		Network:  networkID,
		Subnet:   subnetworkID,
	}
	return errors.OK
}
