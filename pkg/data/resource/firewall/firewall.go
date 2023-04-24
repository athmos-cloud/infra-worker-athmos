package resources

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/common"
	dto "github.com/athmos-cloud/infra-worker-athmos/pkg/common/dto/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/kubernetes"
	domain2 "github.com/athmos-cloud/infra-worker-athmos/pkg/data/project"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/config"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/utils"
)

type Firewall struct {
	Metadata            domain.Metadata         `bson:"metadata"`
	Identifier          Identifier              `bson:"identifier"`
	KubernetesResources kubernetes.ResourceList `bson:"kubernetesResources"`
	Network             string                  `bson:"network"`
	Allow               RuleList                `bson:"allow"`
	Deny                RuleList                `bson:"deny"`
}

type Identifier struct {
	ID         string `bson:"id"`
	ProviderID string `bson:"providerId"`
	VPCID      string `bson:"vpcId"`
	NetworkID  string `bson:"networkId"`
}

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

func (firewall *Firewall) GetMetadata() domain.Metadata {
	return firewall.Metadata
}

func (firewall *Firewall) WithMetadata(request domain.CreateMetadataRequest) {
	firewall.Metadata = domain.New(request)
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

func (firewall *Firewall) InsertIntoProject(project domain2.Project, upsert bool) errors.Error {
	panic("implement me")
}

func (firewall *Firewall) ToDomain() (interface{}, errors.Error) {
	//TODO implement me
	panic("implement me")
}
