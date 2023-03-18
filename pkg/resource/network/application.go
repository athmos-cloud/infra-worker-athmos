package network

type HelmApplication struct {
	Name    string `yaml:"name"`
	Managed bool   `yaml:"managed"`
	VPC     string `yaml:"vpc"`
}
