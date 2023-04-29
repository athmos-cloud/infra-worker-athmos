package resource

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/common"
	resourcePlugin "github.com/athmos-cloud/infra-worker-athmos/pkg/data/plugin"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/status"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/config"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/utils"
	"reflect"
)

type Firewall struct {
	Metadata   metadata.Metadata     `bson:"metadata"`
	Identifier identifier.Firewall   `bson:"identifier"`
	Status     status.ResourceStatus `bson:"status"`
	Network    string                `bson:"network" plugin:"network"`
	Allow      RuleList              `bson:"allow" plugin:"allow"`
	Deny       RuleList              `bson:"deny" plugin:"deny"`
}

func NewFirewall(id identifier.Firewall, providerType common.ProviderType) Firewall {
	return Firewall{
		Metadata: metadata.New(metadata.CreateMetadataRequest{
			Name: id.ID,
		}),
		Status:     status.New(id.ID, common.Firewall, providerType),
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

func (firewall *Firewall) New(id identifier.ID, providerType common.ProviderType) (IResource, errors.Error) {
	if reflect.TypeOf(id) != reflect.TypeOf(identifier.Firewall{}) {
		return nil, errors.InvalidArgument.WithMessage("id type is not FirewallID")
	}
	res := NewFirewall(id.(identifier.Firewall), providerType)
	return &res, errors.OK
}

type Rule struct {
	Protocol string `bson:"protocol" plugin:"protocol"`
	Ports    []int  `bson:"ports" plugin:"ports"`
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

func (firewall *Firewall) GetPluginReference() (resourcePlugin.Reference, errors.Error) {
	if !firewall.Status.PluginReference.ChartReference.Empty() {
		return firewall.Status.PluginReference, errors.OK
	}
	switch firewall.Status.PluginReference.ResourceReference.ProviderType {
	case common.GCP:
		firewall.Status.PluginReference.ChartReference = resourcePlugin.HelmChartReference{
			ChartName:    config.Current.Plugins.Crossplane.GCP.Subnet.Chart,
			ChartVersion: config.Current.Plugins.Crossplane.GCP.Subnet.Version,
		}
		return firewall.Status.PluginReference, errors.OK
	}
	return resourcePlugin.Reference{}, errors.InvalidArgument.WithMessage(fmt.Sprintf("firewall type %s not supported", firewall.Status.PluginReference.ResourceReference.ProviderType))
}

func (firewall *Firewall) FromMap(data map[string]interface{}) errors.Error {
	return resourcePlugin.InjectMapIntoStruct(data, firewall)
}

func (firewall *Firewall) Insert(project Project, update ...bool) errors.Error {
	shouldUpdate := false
	if len(update) > 0 {
		shouldUpdate = update[0]
	}
	id := firewall.Identifier
	_, ok := project.Resources[id.ProviderID].VPCs[id.VPCID].Networks[id.NetworkID].Firewalls[id.ID]
	if !ok && shouldUpdate {
		return errors.NotFound.WithMessage(fmt.Sprintf("network %s not found in vpc %s", id.ID, id.VPCID))
	}
	if ok && !shouldUpdate {
		return errors.Conflict.WithMessage(fmt.Sprintf("network %s already exists in vpc %s", id.ID, id.VPCID))
	}
	project.Resources[id.ProviderID].VPCs[id.VPCID].Networks[id.NetworkID].Firewalls[id.ID] = *firewall
	return errors.OK
}

func (firewall *Firewall) Remove(project Project) errors.Error {
	id := firewall.Identifier
	_, ok := project.Resources[id.ProviderID].VPCs[id.VPCID].Networks[id.NetworkID].Firewalls[id.ID]
	if !ok {
		return errors.NotFound.WithMessage(fmt.Sprintf("network %s not found in vpc %s", id.ID, id.VPCID))
	}
	delete(project.Resources[id.ProviderID].VPCs[id.VPCID].Networks[id.NetworkID].Firewalls, id.ID)
	return errors.NoContent
}

func (firewall *Firewall) Equals(other Firewall) bool {
	return firewall.Metadata.Equals(other.Metadata) &&
		firewall.Identifier.Equals(other.Identifier) &&
		firewall.Status.Equals(other.Status) &&
		firewall.Network == other.Network &&
		firewall.Allow.Equals(other.Allow) &&
		firewall.Deny.Equals(other.Deny)
}
