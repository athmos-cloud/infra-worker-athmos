package identifier

type ID interface{}

type Provider struct {
	ID string `bson:"id"`
}

type VPC struct {
	ID         string
	ProviderID string
}

type Network struct {
	ID         string
	ProviderID string
	VPCID      string
}

type Subnetwork struct {
	ID         string
	ProviderID string
	VPCID      string
	NetworkID  string
}

type Firewall struct {
	ID         string `bson:"id"`
	ProviderID string `bson:"providerId"`
	VPCID      string `bson:"vpcId"`
	NetworkID  string `bson:"networkId"`
}

type VM struct {
	ID         string
	ProviderID string
	VPCID      string
	NetworkID  string
	SubnetID   string
}
