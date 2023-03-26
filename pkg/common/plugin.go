package common

type Plugin struct {
	Name           string         `bson:"name"`
	Provider       ProviderType   `bson:"provider"`
	ResourceType   ResourceType   `bson:"resourceType"`
	ChartReference ChartReference `bson:"chartReference"`
	Inputs         []Input        `bson:"inputs"`
}

type Input struct {
	Name        string      `bson:"name"`
	Description string      `bson:"description"`
	Type        string      `bson:"type"`
	Default     interface{} `bson:"default"`
	Required    bool        `bson:"required"`
}

type ChartReference struct {
	Name    string `bson:"name"`
	Version string `bson:"version"`
}
