package resource

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/common"
	dto "github.com/athmos-cloud/infra-worker-athmos/pkg/common/dto/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/config"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/utils"
)

type Firewall struct {
	Metadata            metadata.Metadata       `bson:"metadata"`
	Identifier          identifier.Firewall     `bson:"identifier"`
	KubernetesResources kubernetes.ResourceList `bson:"kubernetesResources"`
	Network             string                  `bson:"network"`
	Allow               RuleList                `bson:"allow"`
	Deny                RuleList                `bson:"deny"`
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

type Rule struct {
	Protocol string `bson:"protocol"`
	Ports    []int  `bson:"ports"`
}

type RuleList []Rule

func (rules *RuleList) FromMap(data []interface{}) errors.Error {
	*rules = []Rule{}
	for _, rule := range data {
		ruleMap := rule.(map[string]interface{})
		if ruleMap["protocol"] == nil {
			return errors.InvalidArgument.WithMessage("protocol is required")
		}
		if ruleMap["ports"] == nil {
			return errors.InvalidArgument.WithMessage("ports is required")
		}
		*rules = append(*rules, Rule{
			Protocol: ruleMap["protocol"].(string),
			Ports:    ruleMap["ports"].([]int),
		})
	}
	return errors.OK
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
	*firewall = Firewall{}
	if data["id"] == nil {
		firewall.Identifier.ID = utils.GenerateUUID()
	} else {
		firewall.Identifier.ID = data["id"].(string)
	}
	if data["name"] == nil {
		return errors.InvalidArgument.WithMessage("name is required")
	}
	if data["allow"] == nil && data["deny"] == nil {
		return errors.InvalidArgument.WithMessage("allow or deny field is required")
	}
	if data["allow"] != nil {
		firewall.Allow = RuleList{}
		if err := firewall.Allow.FromMap(data["allow"].([]interface{})); !err.IsOk() {
			return err
		}
	}
	if data["deny"] != nil {
		firewall.Deny = RuleList{}
		if err := firewall.Deny.FromMap(data["deny"].([]interface{})); !err.IsOk() {
			return err
		}
	}
	return errors.OK
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

func (firewall *Firewall) ToDomain() (interface{}, errors.Error) {
	//TODO implement me
	panic("implement me")
}
