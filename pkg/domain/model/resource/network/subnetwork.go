package network

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/instance"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/metadata"
)

type Subnetwork struct {
	Metadata       metadata.Metadata     `json:"metadata"`
	IdentifierID   identifier.Subnetwork `json:"identifier_id"`
	IdentifierName identifier.Subnetwork `json:"identifier_name"`
	Region         string                `json:"region"`
	IPCIDRRange    string                `json:"ipcidr_range"`
	VMs            instance.VMCollection `json:"vms,omitempty"`
}

type SubnetworkCollection map[string]Subnetwork
