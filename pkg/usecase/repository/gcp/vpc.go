package gcpRepo

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
)

type VPC interface {
	Find(option.Option) *resource.VPC
	FindAll(option.Option) []*resource.VPC
	Create(*resource.VPC) *resource.VPC
	Update(*resource.VPC) *resource.VPC
	Delete(*resource.VPC)
}
