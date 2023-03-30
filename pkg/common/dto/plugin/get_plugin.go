package plugin

import (
	"github.com/PaulBarrie/infra-worker/pkg/common"
	"github.com/PaulBarrie/infra-worker/pkg/plugin"
)

type GetPluginRequest struct {
	ProviderType common.ProviderType
	ResourceType common.ResourceType
}

type GetPluginResponse struct {
	Plugin plugin.Plugin
}
