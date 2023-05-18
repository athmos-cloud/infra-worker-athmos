package gcpRepo

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
)

type Subnetwork interface {
	Find(option.Option) *resource.Subnetwork
	FindAll(option.Option) []*resource.Subnetwork
	Create(*resource.Subnetwork) *resource.Subnetwork
	Update(*resource.Subnetwork) *resource.Subnetwork
	Delete(*resource.Subnetwork)
}
