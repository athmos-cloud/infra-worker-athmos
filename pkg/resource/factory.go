package resource

import (
	"github.com/PaulBarrie/infra-worker/pkg/common"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/logger"
	"github.com/PaulBarrie/infra-worker/pkg/resource/firewall"
	"github.com/PaulBarrie/infra-worker/pkg/resource/network"
	"github.com/PaulBarrie/infra-worker/pkg/resource/provider"
	"github.com/PaulBarrie/infra-worker/pkg/resource/subnetwork"
	"github.com/PaulBarrie/infra-worker/pkg/resource/vm"
	"github.com/PaulBarrie/infra-worker/pkg/resource/vpc"
)

func ResourceFactory(resourceType common.ResourceType) IResource {
	switch resourceType {
	case common.Provider:
		return &provider.Provider{}
	case common.Network:
		return &network.Network{}
	case common.Firewall:
		return &firewall.Firewall{}
	case common.VPC:
		return &vpc.VPC{}
	case common.Subnetwork:
		return &subnetwork.Subnetwork{}
	case common.VM:
		return &vm.VM{}
	default:
		logger.Error.Printf("Resource type %s not supported", resourceType)
		return nil
	}
}
