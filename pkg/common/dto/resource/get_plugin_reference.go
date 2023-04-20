package resource

import "github.com/athmos-cloud/infra-worker-athmos/pkg/common"

type GetPluginReferenceRequest struct {
	ProviderType common.ProviderType `json:"providerType"`
}

type GetPluginReferenceResponse struct {
	ChartName    string `json:"chartName"`
	ChartVersion string `json:"chartVersion"`
	ReleaseName  string `json:"releaseName"`
}
