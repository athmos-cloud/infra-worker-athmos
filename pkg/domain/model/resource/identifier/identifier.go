package identifier

import "github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"

const (
	providerIdentifierKey   = "identifier.provider"
	vpcIdentifierKey        = "identifier.vpc"
	networkIdentifierKey    = "identifier.network"
	subnetworkIdentifierKey = "identifier.subnetwork"
	vmIdentifierKey         = "identifier.vm"
	firewallIdentifierKey   = "identifier.firewall"
)

type ID interface {
	Equals(other ID) bool
	ToLabels() map[string]string
	FromLabels(labels map[string]string) errors.Error
}

type IdPayload struct {
	ProviderID string `json:"providerID"`
	VPCID      string `json:"vpcID"`
	NetworkID  string `json:"networkID"`
	SubnetID   string `json:"subnetID"`
	VMID       string `json:"vmID"`
	FirewallID string `json:"firewallID"`
}

type Empty struct{}

func (e *Empty) ToLabels() map[string]string {
	//TODO implement me
	panic("implement me")
}

func (e *Empty) FromLabels(_ map[string]string) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (e *Empty) Equals(other ID) bool {
	_, ok := other.(*Empty)
	return ok
}
