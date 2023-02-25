package plugin

import (
	"github.com/PaulBarrie/infra-worker/pkg/kernel/option"
	"github.com/PaulBarrie/infra-worker/pkg/plugin/common"
	"github.com/PaulBarrie/infra-worker/pkg/plugin/common/pipeline"
	"github.com/PaulBarrie/infra-worker/pkg/plugin/terraform"
)

type IPluginStrategy interface {
	Init() pipeline.Pipeline
	Create(option.Option) pipeline.Pipeline
	Get(option.Option) pipeline.Pipeline
	Update(option.Option) pipeline.Pipeline
	Delete(option.Option) pipeline.Pipeline
}

func StrategyBuilder(plugin common.Plugin) IPluginStrategy {
	switch plugin.Type {
	case common.Terraform:
		return terraform.NewStrategy(plugin)
	}
	return nil
}
