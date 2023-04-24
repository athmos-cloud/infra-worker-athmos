package domain

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/common"

	firewall "github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/firewall"
	network "github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/network"
	provider "github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/provider"
	subnetwork "github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/subnetwork"
	vm "github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/vm"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/vpc"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/logger"
)

func Factory(resourceType common.ResourceType) IResource {
	switch resourceType {
	case common.Provider:
		return &provider.Provider{}
	case common.Network:
		return &network.Network{}
	case common.Firewall:
		return &firewall.Firewall{}
	case common.VPC:
		return &resources.VPC{}
	case common.Subnetwork:
		return &subnetwork.Subnetwork{}
	case common.VM:
		return &vm.VM{}
	default:
		logger.Error.Printf("Resource type %s not supported", resourceType)
		return nil
	}
}
