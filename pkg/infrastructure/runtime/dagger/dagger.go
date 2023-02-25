package dagger

import (
	"context"
	"dagger.io/dagger"
	"errors"
	kernel_errors "github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/logger"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/option"
	"os"
	"reflect"
)

const (
	DefaultDaggerWorkdirPermission = 0755
)

type Dagger struct {
	Workdir string
	Client  *dagger.Client
	Process *[]dagger.Container
}

func (d *Dagger) Get() interface{} {
	return d
}

func (d *Dagger) AddProcess(image string, commands []string, volumes map[string]string, envs map[string]string, secrets map[string]string) *Dagger {
	var process *dagger.Container

	if d.Process != nil {
		process = d.Client.Container()
	}
	process.From(image)
	for from, to := range volumes {
		process.WithMountedDirectory(to, d.Client.Host().Directory(from))
	}
	for key, value := range envs {
		process.WithEnvVariable(key, value)
	}
	//for dest, value := range secrets {
	//	d.Client.Address().Secret(dest, value)
	//}

	process.WithExec(commands)
	newProcessList := append(*d.Process, *process)
	d.Process = &newProcessList
	return d
}

func New(ctx context.Context, args ...string) (Dagger, error) {
	optionList := option.NewList(reflect.String, args)
	if !optionList.Validate() {
		return Dagger{}, kernel_errors.InvalidArgument.WithMessage("Args should be a tuple of strings").Error()
	}
	workdir := optionList.Options[0].Get().(string)

	if _, err := os.Stat(workdir); errors.Is(err, os.ErrNotExist) {
		logger.Info.Printf("Dagger local workdir %s does not exists", workdir)
		logger.Info.Printf("Create workdir %s", workdir)
		if err = os.MkdirAll(workdir, DefaultDaggerWorkdirPermission); err != nil {
			logger.Error.Printf("Error creating folder %s", workdir)
			return Dagger{}, err
		}
	}
	client, err := dagger.Connect(ctx, dagger.WithWorkdir(workdir))
	if err != nil {
		return Dagger{}, err
	}
	return Dagger{
		Workdir: workdir,
		Client:  client,
		Process: &[]dagger.Container{},
	}, nil
}

func (d *Dagger) Run(ctx context.Context) error {
	for _, process := range *d.Process {
		if _, err := process.Export(ctx, d.Workdir); err != nil {
			logger.Error.Printf("Error running pipeline %s", err)
			return err
		}
	}
	return nil
}

func (d *Dagger) Clear() {
	d.Process = nil
}
