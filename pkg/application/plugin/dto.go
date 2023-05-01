package plugin

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/plugin"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
)

type GetPluginRequest struct {
	ProviderType types.ProviderType
	ResourceType types.ResourceType
}

type GetPluginResponse struct {
	Plugin plugin.Plugin
}
