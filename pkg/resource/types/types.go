package types

type HelmChartResource struct {
	ChartName    string `bson:"chartName"`
	ChartVersion string `bson:"chartVersion"`
}

type Reference struct {
	Name              string                 `bson:"name"`
	Monitored         bool                   `bson:"monitored,default=true"`
	HelmChartResource HelmChartResource      `bson:"helmChartResource"`
	Values            map[string]interface{} `bson:"values"`
	Tags              map[string]string      `bson:"tags"`
}
