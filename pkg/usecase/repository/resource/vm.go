package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	"sync"
)

type VMChannel struct {
	WaitGroup    *sync.WaitGroup
	Channel      chan *resource.VM
	ErrorChannel chan errors.Error
}

type VM interface {
	FindVM(context.Context, option.Option) (*resource.VM, errors.Error)
	FindAllVMs(context.Context, option.Option) (*resource.VMCollection, errors.Error)
	FindAllRecursiveVMs(context.Context, option.Option, *VMChannel)
	CreateVM(context.Context, *resource.VM) errors.Error
	UpdateVM(context.Context, *resource.VM) errors.Error
	DeleteVM(context.Context, *resource.VM) errors.Error
}
