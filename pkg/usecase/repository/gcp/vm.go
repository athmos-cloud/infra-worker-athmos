package gcpRepo

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
)

type VM interface {
	Find(option.Option) *resource.VM
	FindAll(option.Option) []*resource.VM
	Create(*resource.VM) *resource.VM
	Update(*resource.VM) *resource.VM
	Delete(*resource.VM)
}
