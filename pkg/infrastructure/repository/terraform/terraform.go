package terraform

import (
	"context"
	"fmt"
	"github.com/PaulBarrie/infra-worker/pkg/kernel"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/logger"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/option"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/types/workdir"
	"github.com/PaulBarrie/infra-worker/pkg/plugin/terraform/config"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"os"
	"reflect"
	"strings"
)

type Repository struct{}

type RetrieveRequestPayload struct {
	ConfigWorkdir workdir.Workdir
}

func (r *Repository) Create(ctx context.Context, option option.Option) (interface{}, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) Get(ctx context.Context, option option.Option) (interface{}, errors.Error) {
	if !option.SetType(reflect.TypeOf(RetrieveRequestPayload{}).Kind()).Validate() {
		return nil, errors.InvalidArgument.WithMessage("Invalid option type")
	}
	var result config.Config
	payload := option.Value.(RetrieveRequestPayload)
	for _, file := range option.Value.(workdir.Workdir).Files {
		if !strings.HasSuffix(file.Name, ".tf") {
			continue
		}
		var tmp config.Config
		if err := hclsimple.DecodeFile(file.Name, nil, &tmp); err != nil {
			return nil, errors.IOError.WithMessage(
				fmt.Sprintf("Error decoding file %s", file.Name),
			)
		}

	}
	file, err := os.CreateTemp(kernel.DefaultTmpDir, "*.tf")
	if err != nil {
		return nil, errors.IOError.WithMessage("Error creating temp file")
	}
	defer func(name string) {
		err2 := os.Remove(name)
		if err2 != nil {
			logger.Error.Println(err2)
		}
	}(file.Name())

}

func (r *Repository) GetAll(ctx context.Context, option option.Option) ([]interface{}, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) Update(ctx context.Context, option option.Option) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) Delete(ctx context.Context, option option.Option) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) Close(context context.Context) errors.Error {
	//TODO implement me
	panic("implement me")
}
