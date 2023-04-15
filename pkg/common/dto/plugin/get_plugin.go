package plugin

import (
	"github.com/PaulBarrie/infra-worker/pkg/common"
	"github.com/PaulBarrie/infra-worker/pkg/domain/plugin"
)

type GetPluginRequest struct {
	ProviderType common.ProviderType
	ResourceType common.ResourceType
}

type GetPluginResponse struct {
	Plugin plugin.Plugin
}
