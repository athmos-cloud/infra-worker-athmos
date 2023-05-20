package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/utils"
	"github.com/kamva/mgm/v3"
)

type Firewall struct {
	Metadata       metadata.Metadata   `json:"metadata"`
	IdentifierID   identifier.Firewall `json:"identifierID"`
	IdentifierName identifier.Firewall `json:"identifierName"`
	Allow          RuleList            `json:"allow"`
	Deny           RuleList            `json:"deny"`
}

type FirewallCollection map[string]Firewall

func (collection *FirewallCollection) Equals(other FirewallCollection) bool {
	if len(*collection) != len(other) {
		return false
	}
	for key, value := range *collection {
		if !value.Equals(other[key]) {
			return false
		}
	}
	return true
}

type Rule struct {
	mgm.DefaultModel
	Protocol string
	Ports    []string
}

func (rule *Rule) Equals(other Rule) bool {
	return rule.Protocol == other.Protocol && utils.SliceEquals(rule.Ports, other.Ports)
}

type RuleList []Rule

func (list RuleList) Equals(other RuleList) bool {
	if len(list) != len(other) {
		return false
	}
	for _, value := range list {
		equals := false
		for _, otherValue := range other {
			if value.Equals(otherValue) {
				equals = true
			}
		}
		if !equals {
			return false
		}
	}
	return true
}

func (firewall *Firewall) Equals(other Firewall) bool {
	return firewall.Metadata.Equals(other.Metadata) &&
		firewall.Allow.Equals(other.Allow) &&
		firewall.Deny.Equals(other.Deny)
}
