package plugin

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"github.com/kamva/mgm/v3"
)

type Reference struct {
	mgm.DefaultModel  `bson:",inline"`
	ResourceReference ResourceReference  `bson:"resourceReference"`
	ChartReference    HelmChartReference `bson:"chartReference"`
	Plugin            Plugin             `bson:"plugin"`
}

type ResourceReference struct {
	mgm.DefaultModel `bson:",inline"`
	ResourceType     types.ResourceType `bson:"resourceType"`
	ProviderType     types.ProviderType `bson:"providerType"`
}

type HelmChartReference struct {
	mgm.DefaultModel `bson:",inline"`
	ChartName        string `bson:"chartName"`
	ChartVersion     string `bson:"chartVersion"`
}

func (chart *HelmChartReference) Empty() bool {
	return chart.ChartName == "" || chart.ChartVersion == ""
}
func NewReference(resourceType types.ResourceType, providerType types.ProviderType) Reference {
	resourceReference := ResourceReference{ResourceType: resourceType, ProviderType: providerType}
	return Reference{
		ResourceReference: resourceReference,
		Plugin:            Get(resourceReference),
	}
}
