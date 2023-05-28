package identifier

import "github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"

type Provider struct {
	Provider string `json:"id"`
	VPC      string `json:"vpc,omitempty"`
}

func (provider *Provider) Equals(other ID) bool {
	otherProviderID, ok := other.(*Provider)
	if !ok {
		return false
	}
	return provider.Provider == otherProviderID.Provider && provider.VPC == otherProviderID.VPC
}

func (provider *Provider) ToIDLabels() map[string]string {
	return map[string]string{
		ProviderIdentifierKey: provider.Provider,
		VpcIdentifierKey:      provider.VPC,
	}
}

func (provider *Provider) ToNameLabels() map[string]string {
	return map[string]string{
		ProviderNameKey: provider.Provider,
		VpcNameKey:      provider.VPC,
	}
}

func (provider *Provider) IDFromLabels(labels map[string]string) errors.Error {
	providerID, ok := labels[ProviderIdentifierKey]
	if !ok {
		return errors.InternalError.WithMessage("missing provider identifier")
	}
	vpcID := labels[VpcIdentifierKey]
	*provider = Provider{
		Provider: providerID,
		VPC:      vpcID,
	}
	return errors.OK
}

func (provider *Provider) NameFromLabels(labels map[string]string) errors.Error {
	providerName, ok := labels[ProviderNameKey]
	if !ok {
		return errors.InternalError.WithMessage("missing provider name")
	}
	vpcName := labels[VpcNameKey]
	*provider = Provider{
		Provider: providerName,
		VPC:      vpcName,
	}
	return errors.OK
}
