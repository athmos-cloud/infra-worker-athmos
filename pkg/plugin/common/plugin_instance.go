package common

import "github.com/PaulBarrie/infra-worker/pkg/plugin/common/config"

type PluginInstance struct {
	Name    string
	Plugin  Plugin
	Inputs  config.InputPayloadList
	Outputs config.OutputPayloadList
}
