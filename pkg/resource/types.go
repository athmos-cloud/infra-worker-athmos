package resource

import (
	"github.com/PaulBarrie/infra-worker/pkg/resource/provider"
)

type ResourceType string

const (
	Provider   ResourceType = "provider"
	VPC        ResourceType = "vpc"
	Subnetwork ResourceType = "subnetwork"
	Network    ResourceType = "network"
	VM         ResourceType = "vm"
	Firewall   ResourceType = "firewall"
)

type HelmChartResource struct {
	ChartName    string `bson:"chartName"`
	ChartVersion string `bson:"chartVersion"`
}

type Reference struct {
	Name              string                 `bson:"name""`
	Monitored         bool                   `bson:"monitored,default=true"`
	Provider          provider.Provider      `bson:"provider"`
	HelmChartResource HelmChartResource      `bson:"helmChartResource"`
	Values            map[string]interface{} `bson:"values"`
	Tags              map[string]string      `bson:"tags"`
}
