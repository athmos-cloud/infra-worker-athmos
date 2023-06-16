package dto

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/network"
)

type GetFirewallResponse struct {
	ProjectID string           `json:"project_id"`
	Payload   network.Firewall `json:"payload"`
}

type CreateFirewallRequest struct {
	ParentID   identifier.Network       `json:"parent_id"`
	Name       string                   `json:"name"`
	AllowRules network.FirewallRuleList `json:"allow_rules,omitempty"`
	DenyRules  network.FirewallRuleList `json:"deny_rules,omitempty"`
	Managed    bool                     `json:"managed" default:"true"`
	Tags       map[string]string        `json:"tags"`
}

type CreateFirewallResponse struct {
	ProjectID string           `json:"projectID"`
	Payload   network.Firewall `json:"payload"`
}

type UpdateFirewallRequest struct {
	IdentifierID identifier.Firewall       `json:"identifierID"`
	Name         *string                   `json:"name,omitempty"`
	AllowRules   *network.FirewallRuleList `json:"allowRules,omitempty"`
	DenyRules    *network.FirewallRuleList `json:"denyRules,omitempty"`
	Tags         *map[string]string        `json:"tags,omitempty"`
	Managed      *bool                     `json:"managed"`
}

type DeleteFirewallRequest struct {
	IdentifierID identifier.Firewall `json:"identifierID"`
}
