package gcpRepo

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
)

type Provider interface {
	Find(option.Option) *resource.Provider
	FindAll(option.Option) []*resource.Provider
	Create(*resource.Provider) *resource.Provider
	Update(*resource.Provider) *resource.Provider
	Delete(*resource.Provider)
}
