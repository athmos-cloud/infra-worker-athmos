package terraform

import (
	"github.com/PaulBarrie/infra-worker/pkg/kernel/logger"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/option"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/types/workdir"
	tfConf "github.com/PaulBarrie/infra-worker/pkg/plugin/terraform/config"
)

func GetPluginConfig(workdir workdir.Workdir) tfConf.Config {
	var config tfConf.Config
	parseMerge := func(content string) {
		curConfig, err := tfConf.Marshall(content)
		if !err.IsOk() {
			logger.Warning.Printf("error parsing module content: %v", err)
			return
		}
		config.Merge(*curConfig, option.Null())
	}
	for _, moduleFolder := range workdir.Folders {
		for _, file := range moduleFolder.Files {
			parseMerge(file.Content)
		}
	}
	for _, moduleFile := range workdir.Files {
		parseMerge(moduleFile.Content)
	}
	return config
}
