package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/instance"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
)

type VMChannel struct {
	Channel      chan *instance.VMCollection
	ErrorChannel chan errors.Error
}

type VM interface {
	FindVM(context.Context, option.Option) (*instance.VM, errors.Error)
	FindAllVMs(context.Context, option.Option) (*instance.VMCollection, errors.Error)
	FindAllRecursiveVMs(context.Context, option.Option, *VMChannel)
	CreateVM(context.Context, *instance.VM) errors.Error
	UpdateVM(context.Context, *instance.VM) errors.Error
	DeleteVM(context.Context, *instance.VM) errors.Error
	VMExists(context.Context, *instance.VM) (bool, errors.Error)
}
