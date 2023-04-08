package domain

import (
	"github.com/PaulBarrie/infra-worker/pkg/common"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/logger"
)

func Factory(resourceType common.ResourceType) IResource {
	switch resourceType {
	case common.Provider:
		return &Provider{}
	case common.Network:
		return &Network{}
	case common.Firewall:
		return &Firewall{}
	case common.VPC:
		return &VPC{}
	case common.Subnetwork:
		return &Subnetwork{}
	case common.VM:
		return &VM{}
	default:
		logger.Error.Printf("Resource type %s not supported", resourceType)
		return nil
	}
}
