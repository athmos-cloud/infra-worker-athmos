package firewall

type HelmApplication struct {
	Name    string `yaml:"name"`
	Managed bool   `yaml:"managed"`
	Network string `bson:"network"`
	Allow   []Rule `bson:"allow"`
	Deny    []Rule `bson:"deny"`
}
