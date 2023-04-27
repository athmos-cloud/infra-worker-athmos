package resource

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/common"
	dto "github.com/athmos-cloud/infra-worker-athmos/pkg/common/dto/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/kubernetes"
	resourcePlugin "github.com/athmos-cloud/infra-worker-athmos/pkg/data/plugin"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/config"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/utils"
	"reflect"
)

type Firewall struct {
	Metadata            metadata.Metadata       `bson:"metadata"`
	Identifier          identifier.Firewall     `bson:"identifier"`
	KubernetesResources kubernetes.ResourceList `bson:"kubernetesResources"`
	Network             string                  `bson:"network" plugin:"network"`
	Allow               RuleList                `bson:"allow" plugin:"allow"`
	Deny                RuleList                `bson:"deny" plugin:"deny"`
}

func NewFirewall(id identifier.Firewall) Firewall {
	return Firewall{
		Metadata: metadata.New(metadata.CreateMetadataRequest{
			Name: id.ID,
		}),
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

func (firewall *Firewall) New(id identifier.ID) (IResource, errors.Error) {
	if reflect.TypeOf(id) != reflect.TypeOf(identifier.Firewall{}) {
		return nil, errors.InvalidArgument.WithMessage("id type is not FirewallID")
	}
	res := NewFirewall(id.(identifier.Firewall))
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

func (firewall *Firewall) WithMetadata(request metadata.CreateMetadataRequest) {
	firewall.Metadata = metadata.New(request)
}

func (firewall *Firewall) GetPluginReference(request dto.GetPluginReferenceRequest) (dto.GetPluginReferenceResponse, errors.Error) {
	switch request.ProviderType {
	case common.GCP:
		return dto.GetPluginReferenceResponse{
			ChartName:    config.Current.Plugins.Crossplane.GCP.Firewall.Chart,
			ChartVersion: config.Current.Plugins.Crossplane.GCP.Firewall.Version,
		}, errors.Error{}
	}
	return dto.GetPluginReferenceResponse{}, errors.InvalidArgument.WithMessage(fmt.Sprintf("provider type %s not supported", request.ProviderType))
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
	//TODO implement me
	panic("implement me")
}

func (firewall *Firewall) ToDomain() (interface{}, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (firewall *Firewall) Equals(other Firewall) bool {
	return firewall.Metadata.Equals(other.Metadata) &&
		firewall.Identifier.Equals(other.Identifier) &&
		firewall.Network == other.Network &&
		firewall.Allow.Equals(other.Allow) &&
		firewall.Deny.Equals(other.Deny)
}
