package pipeline

import (
	"context"
	"github.com/PaulBarrie/infra-worker/pkg/infrastructure/runtime"
)

type IProcess interface {
	Build() error
	Run(ctx context.Context) error
	Stop() error
	Status() error
	GetOutput() (string, error)
}

type Process struct {
	Image   string
	Runtime runtime.IRuntime
	Script  string
	Volumes map[string]string
	Secrets map[string]string
}

func NewProcess(image string, script string, volumes map[string]string, secrets map[string]string) Process {
	return Process{
		Image:   image,
		Script:  script,
		Volumes: volumes,
		Secrets: secrets,
	}
}

func (i *Process) Build() error {
	i.Runtime = i.Runtime.AddProcess(i.Image, []string{i.Script}, i.Volumes, nil, i.Secrets)
	return nil
}

func (i *Process) Run(ctx context.Context) error {
	err := i.Runtime.Run(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (i *Process) Stop() error {
	//TODO implement me
	panic("implement me")
}

func (i *Process) Status() error {
	//TODO implement me
	panic("implement me")
}

func (i *Process) GetOutput() (string, error) {
	//TODO implement me
	panic("implement me")
}
