package subnetwork

type HelmApplication struct {
	Name        string `yaml:"name"`
	Managed     bool   `yaml:"managed"`
	VPC         string `yaml:"vpc"`
	Network     string `bson:"network"`
	IPCIDRRange string `yaml:"ipCidrRange"`
	Region      string `bson:"region"`
}
