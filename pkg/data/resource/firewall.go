package resource

import (
	"fmt"
	resourcePlugin "github.com/athmos-cloud/infra-worker-athmos/pkg/data/plugin"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/status"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/config"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/utils"
	"github.com/kamva/mgm/v3"
	"reflect"
)

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
		FirewallID: fmt.Sprintf("%s-%s", payload.Name, utils.RandomString(resourceIDSuffixLength)),
	}
	return Firewall{
		Metadata: metadata.New(metadata.CreateMetadataRequest{
			Name:         id.FirewallID,
			NotMonitored: !payload.Monitored,
			Tags:         payload.Tags,
		}),
		Status:     status.New(id.FirewallID, types.Firewall, payload.Provider),
		Identifier: id,
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
	Protocol         string `bson:"protocol" plugin:"protocol"`
	Ports            []int  `bson:"ports" plugin:"ports"`
}

func (rule *Rule) Equals(other Rule) bool {
	return rule.Protocol == other.Protocol && utils.IntSliceEquals(rule.Ports, other.Ports)
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

type Firewall struct {
	mgm.DefaultModel `bson:",inline"`
	Metadata         metadata.Metadata     `bson:"metadata"`
	Identifier       identifier.Firewall   `bson:"identifier"`
	Status           status.ResourceStatus `bson:"status"`
	Allow            RuleList              `bson:"allow" plugin:"allow"`
	Deny             RuleList              `bson:"deny" plugin:"deny"`
}

func (firewall *Firewall) GetIdentifier() identifier.ID {
	return firewall.Identifier
}

func (firewall *Firewall) New(payload NewResourcePayload) IResource {
	if reflect.TypeOf(payload.ParentIdentifier) != reflect.TypeOf(identifier.Firewall{}) {
		panic(errors.InvalidArgument.WithMessage("id type is not FirewallID"))
	}
	res := NewFirewall(payload)
	return &res
}

func (firewall *Firewall) GetMetadata() metadata.Metadata {
	return firewall.Metadata
}

func (firewall *Firewall) SetMetadata(request metadata.CreateMetadataRequest) {
	firewall.Metadata = metadata.New(request)
}

func (firewall *Firewall) SetStatus(resourceStatus status.ResourceStatus) {
	firewall.Status = resourceStatus
}

func (firewall *Firewall) GetStatus() status.ResourceStatus {
	return firewall.Status
}

func (firewall *Firewall) GetPluginReference() resourcePlugin.Reference {
	if !firewall.Status.PluginReference.ChartReference.Empty() {
		return firewall.Status.PluginReference
	}
	switch firewall.Status.PluginReference.ResourceReference.ProviderType {
	case types.GCP:
		firewall.Status.PluginReference.ChartReference = resourcePlugin.HelmChartReference{
			ChartName:    config.Current.Plugins.Crossplane.GCP.Firewall.Chart,
			ChartVersion: config.Current.Plugins.Crossplane.GCP.Firewall.Version,
		}
		return firewall.Status.PluginReference
	}
	panic(errors.InvalidArgument.WithMessage(fmt.Sprintf("firewall type %s not supported", firewall.Status.PluginReference.ResourceReference.ProviderType)))
}

func (firewall *Firewall) FromMap(data map[string]interface{}) {
	if err := resourcePlugin.InjectMapIntoStruct(data, firewall); !err.IsOk() {
		panic(err)
	}
}

func (firewall *Firewall) Insert(_ IResource, _ ...bool) {
	return
}

func (firewall *Firewall) Remove(_ IResource) {
	return
}

func (firewall *Firewall) Equals(other Firewall) bool {
	return firewall.Metadata.Equals(other.Metadata) &&
		firewall.Identifier.Equals(other.Identifier) &&
		firewall.Status.Equals(other.Status) &&
		firewall.Allow.Equals(other.Allow) &&
		firewall.Deny.Equals(other.Deny)
}
