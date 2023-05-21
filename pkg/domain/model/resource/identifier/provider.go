package identifier

import "github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"

const (
	ProviderLabelKey  = "name.provider"
	ProjectIDLabelKey = "athmos/project-id"
)

type Provider struct {
	ProjectID string `json:"projectID"`
	Provider  string `json:"id"`
	VPC       string `json:"vpc,omitempty"`
}

func (provider *Provider) Equals(other ID) bool {
	otherProviderID, ok := other.(*Provider)
	if !ok {
		return false
	}
	return provider.Provider == otherProviderID.Provider && provider.VPC == otherProviderID.VPC
}

func (provider *Provider) ToLabels() map[string]string {
	return map[string]string{
		ProviderIdentifierKey: provider.Provider,
	}
}

func (provider *Provider) FromLabels(labels map[string]string) errors.Error {
	providerID, ok := labels[ProviderIdentifierKey]
	if !ok {
		return errors.InternalError.WithMessage("missing provider identifier")
	}
	*provider = Provider{
		Provider: providerID,
	}
	return errors.OK
}

func (provider *Provider) GetLabelName() map[string]string {
	return map[string]string{
		ProjectIDLabelKey: provider.ProjectID,
		ProviderLabelKey:  provider.Provider,
	}
}
