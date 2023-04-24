package domain

type VM struct {
	Name        string   `bson:"name"`
	VPC         string   `bson:"vpc"`
	Network     string   `bson:"network"`
	Subnetwork  string   `bson:"subnetwork"`
	Zone        string   `bson:"zone"`
	MachineType string   `bson:"machineType"`
	Auths       []VMAuth `bson:"auths"`
	Disks       []Disk   `bson:"disks"`
	OS          OS       `bson:"os"`
}

type VMHierarchyLocation struct {
	ProviderID string
	VPCID      string
	NetworkID  string
	SubnetID   string
}

type Disk struct {
	Type       string   `bson:"type"`
	Mode       DiskMode `bson:"mode"`
	SizeGib    int      `bson:"sizeGib"`
	AutoDelete bool     `bson:"autoDelete"`
}

type DiskMode string

const (
	READ_ONLY  DiskMode = "READ_ONLY"
	READ_WRITE DiskMode = "READ_WRITE"
)

type VMAuth struct {
	Username     string `bson:"username"`
	SSHPublicKey string `bson:"sshPublicKey"`
}

type OS struct {
	Type    string `bson:"type"`
	Version string `bson:"version"`
}
