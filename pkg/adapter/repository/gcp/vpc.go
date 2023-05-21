package gcp

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
)

func (gcp *gcpRepository) FindVPC(ctx context.Context, opt option.Option) (*resource.VPC, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (gcp *gcpRepository) FindAllRecursiveVPCs(ctx context.Context, opt option.Option) (*resource.VPCCollection, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (gcp *gcpRepository) FindAllVPCs(ctx context.Context, opt option.Option) (*resource.VPCCollection, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (gcp *gcpRepository) CreateVPC(ctx context.Context, vpc *resource.VPC) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (gcp *gcpRepository) UpdateVPC(ctx context.Context, vpc *resource.VPC) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (gcp *gcpRepository) DeleteVPC(ctx context.Context, vpc *resource.VPC) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (gcp *gcpRepository) DeleteVPCCascade(ctx context.Context, vpc *resource.VPC) errors.Error {
	//TODO implement me
	panic("implement me")
}
