package identifier

import "github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"

type VPC struct {
	VPC      string `json:"vpc"`
	Provider string `json:"provider"`
}

func (id *VPC) Equals(other ID) bool {
	otherVPCID, ok := other.(*VPC)
	if !ok {
		return false
	}
	return id.VPC == otherVPCID.VPC &&
		id.Provider == otherVPCID.Provider
}

func (id *VPC) ToLabels() map[string]string {
	return map[string]string{
		VpcIdentifierKey:      id.VPC,
		ProviderIdentifierKey: id.Provider,
	}
}

func (id *VPC) FromLabels(labels map[string]string) errors.Error {
	vpcID, ok := labels[VpcIdentifierKey]
	if !ok {
		return errors.InternalError.WithMessage("missing vpc identifier")
	}
	providerID, ok := labels[ProviderIdentifierKey]
	if !ok {
		return errors.InternalError.WithMessage("missing provider identifier")
	}
	*id = VPC{
		VPC:      vpcID,
		Provider: providerID,
	}
	return errors.OK
}
