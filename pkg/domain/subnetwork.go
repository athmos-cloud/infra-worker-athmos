package domain

type Subnetwork struct {
	ID          string
	Name        string
	IPCIDRRange string
	Region      string `bson:"region"`
	VMs         map[string]VM
}
