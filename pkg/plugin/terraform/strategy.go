package terraform

import (
	"github.com/PaulBarrie/infra-worker/pkg/kernel/option"
	"github.com/PaulBarrie/infra-worker/pkg/plugin/common"
	pipeline2 "github.com/PaulBarrie/infra-worker/pkg/plugin/common/pipeline"
)

type Strategy struct {
	Plugin   common.Plugin
	Workdir  string
	Pipeline pipeline2.Pipeline
}

func (s *Strategy) Init() pipeline2.Pipeline {
	//TODO implement me
	panic("implement me")
}

func (s *Strategy) Create(option option.Option) pipeline2.Pipeline {
	//TODO implement me
	panic("implement me")
}

func (s *Strategy) Get(option option.Option) pipeline2.Pipeline {
	//TODO implement me
	panic("implement me")
}

func (s *Strategy) Update(option option.Option) pipeline2.Pipeline {
	//TODO implement me
	panic("implement me")
}

func (s *Strategy) Delete(option option.Option) pipeline2.Pipeline {
	//TODO implement me
	panic("implement me")
}

func NewStrategy(plugin common.Plugin) *Strategy {
	return &Strategy{
		Plugin: plugin,
	}
}
