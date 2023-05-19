package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/utils"
	"github.com/kamva/mgm/v3"
	"reflect"
)

type Firewall struct {
	Metadata   metadata.Metadata   `bson:"metadata"`
	Identifier identifier.Firewall `bson:"identifier"`
	Allow      RuleList            `bson:"allow" plugin:"allow" yaml:"allow"`
	Deny       RuleList            `bson:"deny" plugin:"deny" yaml:"deny"`
}

func NewFirewall(payload NewResourcePayload) Firewall {
	payload.Validate()
	if reflect.TypeOf(payload.ParentIdentifier) != reflect.TypeOf(identifier.Network{}) {
		panic(errors.InvalidArgument.WithMessage("ID type must be network ID"))
	}
	parentID := payload.ParentIdentifier.(identifier.Network)
	id := identifier.Firewall{
		ProviderID: parentID.ProviderID,
		VPCID:      parentID.VPCID,
		NetworkID:  parentID.NetworkID,
		FirewallID: formatResourceName(payload.Name),
	}
	return Firewall{
		Metadata: metadata.New(metadata.CreateMetadataRequest{
			Name:         id.FirewallID,
			NotMonitored: !payload.Managed,
			Tags:         payload.Tags,
		}),
		Identifier: id,
		Allow:      make(RuleList, 0),
		Deny:       make(RuleList, 0),
	}
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
	mgm.DefaultModel `bson:",inline"`
	Protocol         string   `bson:"protocol" plugin:"protocol" yaml:"protocol"`
	Ports            []string `bson:"ports" plugin:"ports" yaml:"ports"`
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
		firewall.Identifier.Equals(other.Identifier) &&
		firewall.Allow.Equals(other.Allow) &&
		firewall.Deny.Equals(other.Deny)
}
