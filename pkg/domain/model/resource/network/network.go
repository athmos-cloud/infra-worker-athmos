package network

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/instance"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/metadata"
)

type Network struct {
	Metadata       metadata.Metadata        `json:"metadata"`
	IdentifierID   identifier.Network       `json:"identifier_id"`
	IdentifierName identifier.Network       `json:"identifier_name"`
	Subnetworks    SubnetworkCollection     `json:"subnetworks,omitempty"`
	Firewalls      FirewallCollection       `json:"firewalls,omitempty"`
	SqlDbs         instance.SqlDBCollection `json:"sqlDbs,omitempty"`
}

type NetworkCollection map[string]Network
