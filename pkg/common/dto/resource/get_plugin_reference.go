package resource

import "github.com/athmos-cloud/infra-worker-athmos/pkg/common"

type GetPluginReferenceRequest struct {
	ProviderType common.ProviderType `json:"providerType"`
}

type GetPluginReferenceResponse struct {
}
