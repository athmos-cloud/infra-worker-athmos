package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
)

type VPC interface {
	FindVPC(context.Context, option.Option) (*resource.VPC, errors.Error)
	FindAllVPCs(context.Context, option.Option) (*resource.VPCCollection, errors.Error)
	CreateVPC(context.Context, *resource.VPC) errors.Error
	UpdateVPC(context.Context, *resource.VPC) errors.Error
	DeleteVPC(context.Context, *resource.VPC) errors.Error
}
