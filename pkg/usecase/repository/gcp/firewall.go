package gcpRepo

import (
	resource2 "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
)

type Firewall interface {
	Find(option.Option) *resource2.Firewall
	FindAll(option.Option) []*resource2.Firewall
	Create(*resource2.Provider) *resource2.Firewall
	Update(*resource2.Provider) *resource2.Firewall
	Delete(*resource2.Firewall)
}
