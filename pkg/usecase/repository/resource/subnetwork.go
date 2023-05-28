package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	"sync"
)

type SubnetworkChannel struct {
	WaitGroup    *sync.WaitGroup
	Channel      chan *resource.SubnetworkCollection
	ErrorChannel chan errors.Error
}

type Subnetwork interface {
	FindSubnetwork(context.Context, option.Option) (*resource.Subnetwork, errors.Error)
	FindAllSubnetworks(context.Context, option.Option) (*resource.SubnetworkCollection, errors.Error)
	FindAllRecursiveSubnetworks(context.Context, option.Option, *SubnetworkChannel)
	CreateSubnetwork(context.Context, *resource.Subnetwork) errors.Error
	UpdateSubnetwork(context.Context, *resource.Subnetwork) errors.Error
	DeleteSubnetwork(context.Context, *resource.Subnetwork) errors.Error
	DeleteSubnetworkCascade(context.Context, *resource.Subnetwork) errors.Error
	SubnetworkExists(context.Context, option.Option) (bool, errors.Error)
}
