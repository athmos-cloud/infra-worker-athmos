package project

import (
	"github.com/PaulBarrie/infra-worker/pkg/plugin/common"
	comConfig "github.com/PaulBarrie/infra-worker/pkg/plugin/common/config"
)

type PluginInstance struct {
	PluginModel common.Plugin
	Variables   comConfig.InputPayload
	Outputs     comConfig.OutputPayload
}
