package runtime

import (
	"context"
	"errors"
	"fmt"
	"github.com/PaulBarrie/infra-worker/pkg/infrastructure/runtime/dagger"
)

var (
	CurrentRuntime IRuntime
)

type RuntimeType string

const (
	DAGGER RuntimeType = "dagger"
)

type IRuntime interface {
	Get() interface{}
	AddProcess(image string, strings []string, volumes map[string]string, envs map[string]string, secrets map[string]string) IRuntime
	Run(context.Context) error
	Clear()
}

func New(ctx context.Context, workdir string, runtimeType RuntimeType) (interface{}, error) {
	switch runtimeType {
	case DAGGER:
		daggerRuntime, err := dagger.New(ctx, workdir)
		if err != nil {
			return nil, err
		}
		return &daggerRuntime, nil
	default:
		return nil, errors.New(fmt.Sprintf("IRuntime of type %s is not handled", runtimeType))
	}
}
