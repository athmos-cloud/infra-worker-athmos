package network

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/metadata"
)

type Firewall struct {
	Metadata       metadata.Metadata   `json:"metadata"`
	IdentifierID   identifier.Firewall `json:"identifier_id"`
	IdentifierName identifier.Firewall `json:"identifier_name"`
	Allow          FirewallRuleList    `json:"allow"`
	Deny           FirewallRuleList    `json:"deny"`
}

type FirewallCollection map[string]Firewall

type FirewallRule struct {
	Protocol string
	Ports    []string
}

type FirewallRuleList []FirewallRule
