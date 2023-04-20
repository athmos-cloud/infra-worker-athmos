package plugin

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/common"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/plugin"
)

type GetPluginRequest struct {
	ProviderType common.ProviderType
	ResourceType common.ResourceType
}

type GetPluginResponse struct {
	Plugin plugin.Plugin
}
