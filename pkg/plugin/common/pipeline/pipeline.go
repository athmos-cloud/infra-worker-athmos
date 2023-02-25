package pipeline

import (
	"context"
	"github.com/PaulBarrie/infra-worker/pkg/plugin/common/config"
)

type Pipeline struct {
	Processes []Process
	Outputs   []config.Output
	Context   context.Context
}

func (p *Pipeline) Build() error {
	for _, instruction_ := range p.Processes {
		err := instruction_.Build()
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Pipeline) Run() error {
	for _, instruction_ := range p.Processes {
		err := instruction_.Run(p.Context)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Pipeline) Status() error {
	//TODO implement me
	panic("implement me")
}
func (p *Pipeline) Stop() error {
	//TODO implement me
	panic("implement me")
}

func (p *Pipeline) AddProcess(process Process) {

	p.Processes = append(p.Processes, process)
}
