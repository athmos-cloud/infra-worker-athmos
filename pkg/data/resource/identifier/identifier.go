package identifier

type ID interface{}

type Provider struct {
	ID string `bson:"id"`
}

func (provider *Provider) Equals(other Provider) bool {
	return provider.ID == other.ID
}

type VPC struct {
	ID         string
	ProviderID string
}

func (vpc *VPC) Equals(other VPC) bool {
	return vpc.ID == other.ID &&
		vpc.ProviderID == other.ProviderID
}

type Network struct {
	ID         string
	ProviderID string
	VPCID      string
}

func (network *Network) Equals(other Network) bool {
	return network.ID == other.ID &&
		network.ProviderID == other.ProviderID &&
		network.VPCID == other.VPCID
}

type Subnetwork struct {
	ID         string
	ProviderID string
	VPCID      string
	NetworkID  string
}

func (subnetwork *Subnetwork) Equals(other Subnetwork) bool {
	return subnetwork.ID == other.ID &&
		subnetwork.ProviderID == other.ProviderID &&
		subnetwork.VPCID == other.VPCID &&
		subnetwork.NetworkID == other.NetworkID
}

type Firewall struct {
	ID         string `bson:"id"`
	ProviderID string `bson:"providerId"`
	VPCID      string `bson:"vpcId"`
	NetworkID  string `bson:"networkId"`
}

func (firewall *Firewall) Equals(other Firewall) bool {
	return firewall.ID == other.ID &&
		firewall.ProviderID == other.ProviderID &&
		firewall.VPCID == other.VPCID &&
		firewall.NetworkID == other.NetworkID
}

type VM struct {
	ID         string
	ProviderID string
	VPCID      string
	NetworkID  string
	SubnetID   string
}

func (id *VM) Equals(other VM) bool {
	return id.ID == other.ID &&
		id.ProviderID == other.ProviderID &&
		id.VPCID == other.VPCID &&
		id.NetworkID == other.NetworkID &&
		id.SubnetID == other.SubnetID
}
