package domain

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"reflect"
)

type IResource interface {
	ToDataMapper(resource.IResource) resource.IResource
}

func FromDataMapper(resourceInput resource.IResource) IResource {
	switch reflect.TypeOf(resourceInput) {
	case reflect.TypeOf(resource.NewProvider(identifier.Provider{}, "")):
		provider := resourceInput.(*resource.Provider)
		return FromProviderDataMapper(provider)
	case reflect.TypeOf(resource.Network{}):
		network := resourceInput.(*resource.Network)
		return FromNetworkDataMapper(network)
	case reflect.TypeOf(resource.Firewall{}):
		firewall := resourceInput.(*resource.Firewall)
		return FromFirewallDataMapper(firewall)
	case reflect.TypeOf(resource.VPC{}):
		vpc := resourceInput.(*resource.VPC)
		return FromVPCDataMapper(vpc)
	case reflect.TypeOf(resource.Subnetwork{}):
		subnet := resourceInput.(*resource.Subnetwork)
		return FromSubnetworkDataMapper(subnet)
	case reflect.TypeOf(resource.VM{}):
		vm := resourceInput.(*resource.VM)
		return FromVMDataMapper(vm)
	default:
		panic(errors.InvalidArgument.WithMessage(
			fmt.Sprintf("resource type %s not supported", reflect.TypeOf(resourceInput).Kind()),
		))
	}
}
