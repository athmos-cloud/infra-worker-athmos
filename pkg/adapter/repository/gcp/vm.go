package gcp

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
)

func (gcp *gcpRepository) FindVM(ctx context.Context, opt option.Option) (*resource.VM, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (gcp *gcpRepository) FindAllVMs(ctx context.Context, opt option.Option) (*resource.VMCollection, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (gcp *gcpRepository) CreateVM(ctx context.Context, vm *resource.VM) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (gcp *gcpRepository) UpdateVM(ctx context.Context, vm *resource.VM) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (gcp *gcpRepository) DeleteVM(ctx context.Context, vm *resource.VM) errors.Error {
	//TODO implement me
	panic("implement me")
}
