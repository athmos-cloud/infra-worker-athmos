package domain

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
)

type Firewall struct {
	Name      string             `json:"name"`
	Monitored bool               `json:"monitored"`
	Type      types.ProviderType `json:"type"`
	Allow     RuleList           `json:"allow"`
	Deny      RuleList           `json:"deny"`
}

func (firewall Firewall) ToDataMapper(resourceInput resource.IResource) resource.IResource {
	firewallInput := resourceInput.(*resource.Firewall)
	firewallInput.Identifier.ID = firewall.Name
	firewallInput.Metadata.Managed = firewall.Monitored
	firewallInput.Status.PluginReference.ResourceReference.ProviderType = firewall.Type
	newAllow := make(resource.RuleList, len(firewall.Allow))
	for i, rule := range firewall.Allow {
		newAllow[i] = resource.Rule{
			Protocol: rule.Protocol,
			Ports:    rule.Ports,
		}
	}
	firewallInput.Allow = newAllow
	newDeny := make(resource.RuleList, len(firewall.Deny))
	for i, rule := range firewall.Deny {
		newDeny[i] = resource.Rule{
			Protocol: rule.Protocol,
			Ports:    rule.Ports,
		}
	}
	firewallInput.Deny = newDeny
	return firewallInput
}

func FromFirewallDataMapper(firewall *resource.Firewall) Firewall {
	return Firewall{
		Name:      firewall.Identifier.ID,
		Monitored: firewall.Metadata.Managed,
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
		result[firewall.Identifier.ID] = FromFirewallDataMapper(&firewall)
	}
	return result
}
