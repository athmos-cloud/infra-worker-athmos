package domain

import "github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource"

type Firewall struct {
	Name      string   `json:"name"`
	Monitored bool     `json:"monitored"`
	Allow     RuleList `json:"allow"`
	Deny      RuleList `json:"deny"`
}

func FromFirewallDataMapper(firewall resource.Firewall) Firewall {
	return Firewall{
		Name:      firewall.Identifier.ID,
		Monitored: firewall.Metadata.Monitored,
		Allow:     FromRuleListDataMapper(firewall.Allow),
		Deny:      FromRuleListDataMapper(firewall.Deny),
	}
}

type Rule struct {
	Protocol string `json:"protocol"`
	Ports    []int  `json:"ports"`
}

func FromRuleDataMapper(rule resource.Rule) Rule {
	return Rule{
		Protocol: rule.Protocol,
		Ports:    rule.Ports,
	}
}

type RuleList []Rule

func FromRuleListDataMapper(rules resource.RuleList) RuleList {
	result := make(RuleList, len(rules))
	for i, rule := range rules {
		result[i] = FromRuleDataMapper(rule)
	}
	return result
}

type FirewallCollection map[string]Firewall

func FromFirewallCollectionDataMapper(firewalls resource.FirewallCollection) FirewallCollection {
	result := make(FirewallCollection)
	for _, firewall := range firewalls {
		result[firewall.Identifier.ID] = FromFirewallDataMapper(firewall)
	}
	return result
}
