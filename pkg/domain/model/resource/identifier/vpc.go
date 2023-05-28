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

func (id *VPC) ToIDLabels() map[string]string {
	return map[string]string{
		VpcIdentifierKey:      id.VPC,
		ProviderIdentifierKey: id.Provider,
	}
}

func (id *VPC) IDFromLabels(labels map[string]string) errors.Error {
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

func (id *VPC) ToNameLabels() map[string]string {
	return map[string]string{
		VpcNameKey:      id.VPC,
		ProviderNameKey: id.Provider,
	}
}

func (id *VPC) NameFromLabels(labels map[string]string) errors.Error {
	vpcName, ok := labels[VpcNameKey]
	if !ok {
		return errors.InternalError.WithMessage("missing vpc name")
	}
	providerName, ok := labels[ProviderNameKey]
	if !ok {
		return errors.InternalError.WithMessage("missing provider name")
	}
	*id = VPC{
		VPC:      vpcName,
		Provider: providerName,
	}
	return errors.OK
}
