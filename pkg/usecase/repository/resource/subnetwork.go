package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/network"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
)

type SubnetworkChannel struct {
	Channel      chan *network.SubnetworkCollection
	ErrorChannel chan errors.Error
}

type Subnetwork interface {
	FindSubnetwork(context.Context, option.Option) (*network.Subnetwork, errors.Error)
	FindAllSubnetworks(context.Context, option.Option) (*network.SubnetworkCollection, errors.Error)
	FindAllRecursiveSubnetworks(context.Context, option.Option, *SubnetworkChannel)
	CreateSubnetwork(context.Context, *network.Subnetwork) errors.Error
	UpdateSubnetwork(context.Context, *network.Subnetwork) errors.Error
	DeleteSubnetwork(context.Context, *network.Subnetwork) errors.Error
	DeleteSubnetworkCascade(context.Context, *network.Subnetwork) errors.Error
	SubnetworkExists(context.Context, *network.Subnetwork) (bool, errors.Error)
}
