package plugin

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/common"
)

type Reference struct {
	ResourceReference ResourceReference  `bson:"resourceReference"`
	ChartReference    HelmChartReference `bson:"chartReference"`
	Plugin            Plugin             `bson:"plugin"`
}

type ResourceReference struct {
	ResourceType common.ResourceType `bson:"resourceType"`
	ProviderType common.ProviderType `bson:"providerType"`
}

type HelmChartReference struct {
	ChartName    string `bson:"chartName"`
	ChartVersion string `bson:"chartVersion"`
}

func (chart *HelmChartReference) Empty() bool {
	return chart.ChartName == "" || chart.ChartVersion == ""
}
func NewReference(resourceType common.ResourceType, providerType common.ProviderType) Reference {
	resourceReference := ResourceReference{ResourceType: resourceType, ProviderType: providerType}
	return Reference{
		ResourceReference: resourceReference,
		Plugin:            Get(resourceReference),
	}
}
