package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	"sync"
)

type VPCChannel struct {
	WaitGroup    *sync.WaitGroup
	Channel      chan *resource.VPCCollection
	ErrorChannel chan errors.Error
}

type VPC interface {
	FindVPC(context.Context, option.Option) (*resource.VPC, errors.Error)
	FindAllVPCs(context.Context, option.Option) (*resource.VPCCollection, errors.Error)
	FindAllRecursiveVPCs(context.Context, option.Option, *VPCChannel)
	CreateVPC(context.Context, *resource.VPC) errors.Error
	UpdateVPC(context.Context, *resource.VPC) errors.Error
	DeleteVPC(context.Context, *resource.VPC) errors.Error
	DeleteVPCCascade(context.Context, *resource.VPC) errors.Error
}