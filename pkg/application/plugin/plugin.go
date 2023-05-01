package plugin

import (
	plugin2 "github.com/athmos-cloud/infra-worker-athmos/pkg/data/plugin"
)

type Service struct{}

func (service *Service) GetPlugin(request GetPluginRequest) GetPluginResponse {
	plugin := plugin2.Get(plugin2.ResourceReference{ResourceType: request.ResourceType, ProviderType: request.ProviderType})

	return GetPluginResponse{Plugin: plugin}
}
