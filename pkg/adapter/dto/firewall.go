package dto

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
)

type GetFirewallRequest struct {
	IdentifierID identifier.Firewall `json:"identifier_id"`
}

type GetFirewallResponse struct {
	ProjectID string            `json:"project_id"`
	Payload   resource.Firewall `json:"payload"`
}

type CreateFirewallRequest struct {
	ParentID   identifier.Network        `json:"parent_id"`
	Name       string                    `json:"name"`
	AllowRules resource.FirewallRuleList `json:"allow_rules,omitempty"`
	DenyRules  resource.FirewallRuleList `json:"deny_rules,omitempty"`
	Managed    bool                      `json:"managed" default:"true"`
	Tags       map[string]string         `json:"tags"`
}

type CreateFirewallResponse struct {
	ProjectID string            `json:"projectID"`
	Payload   resource.Firewall `json:"payload"`
}

type UpdateFirewallRequest struct {
	IdentifierID identifier.Firewall        `json:"identifierID"`
	Name         *string                    `json:"name,omitempty"`
	AllowRules   *resource.FirewallRuleList `json:"allowRules,omitempty"`
	DenyRules    *resource.FirewallRuleList `json:"denyRules,omitempty"`
	Tags         *map[string]string         `json:"tags,omitempty"`
	Managed      *bool                      `json:"managed"`
}

type DeleteFirewallRequest struct {
	IdentifierID identifier.Firewall `json:"identifierID"`
}
