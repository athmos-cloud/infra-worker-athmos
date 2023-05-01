package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/logger"
)

func Factory(resourceType types.ResourceType) IResource {
	switch resourceType {
	case types.Provider:
		return &Provider{}
	case types.Network:
		return &Network{}
	case types.Firewall:
		return &Firewall{}
	case types.VPC:
		return &VPC{}
	case types.Subnetwork:
		return &Subnetwork{}
	case types.VM:
		return &VM{}
	default:
		logger.Error.Printf("Resource type %s not supported", resourceType)
		return nil
	}
}
