package resource

import (
	"github.com/PaulBarrie/infra-worker/pkg/kernel/logger"
	"github.com/PaulBarrie/infra-worker/pkg/resource/firewall"
	"github.com/PaulBarrie/infra-worker/pkg/resource/network"
	"github.com/PaulBarrie/infra-worker/pkg/resource/provider"
	"github.com/PaulBarrie/infra-worker/pkg/resource/subnetwork"
	"github.com/PaulBarrie/infra-worker/pkg/resource/vm"
	"github.com/PaulBarrie/infra-worker/pkg/resource/vpc"
)

func ResourceFactory(resourceType ResourceType) IResource {
	switch resourceType {
	case Provider:
		return &provider.Provider{}
	case Network:
		return &network.Network{}
	case Firewall:
		return &firewall.Firewall{}
	case VPC:
		return &vpc.VPC{}
	case Subnetwork:
		return &subnetwork.Subnetwork{}
	case VM:
		return &vm.VM{}
	default:
		logger.Error.Printf("Resource type %s not supported", resourceType)
		return nil
	}
}
