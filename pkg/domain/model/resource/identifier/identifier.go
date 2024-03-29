package identifier

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

const (
	ProviderIdentifierKey   = "identifier.provider"
	VpcIdentifierKey        = "identifier.vpc"
	NetworkIdentifierKey    = "identifier.network"
	SubnetworkIdentifierKey = "identifier.subnetwork"
	VMIdentifierKey         = "identifier.vm"
	SqlDBIdentifierKey      = "identifier.sqldb"
	FirewallIdentifierKey   = "identifier.firewall"
	ProviderNameKey         = "name.provider"
	VpcNameKey              = "name.vpc"
	NetworkNameKey          = "name.network"
	SubnetworkNameKey       = "name.subnetwork"
	VMNameKey               = "name.vm"
	FirewallNameKey         = "name.firewall"
	SqlDBNameKey            = "name.sqldb"
)

type ID interface {
	Equals(other ID) bool
	ToIDLabels() map[string]string
	ToNameLabels() map[string]string
	IDFromLabels(labels map[string]string) errors.Error
	NameFromLabels(labels map[string]string) errors.Error
}

type Payload struct {
	Provider   string `json:"provider"`
	VPC        string `json:"vpc"`
	Network    string `json:"network"`
	Subnetwork string `json:"subnetwork"`
	VM         string `json:"vm"`
	Firewall   string `json:"firewall"`
	SqlDB      string `json:"sqldb"`
}


func FromPayload(payload Payload) ID {
	if payload.Provider != "" && (payload.VPC != "" || payload.Network != "") && payload.Subnetwork != "" && payload.VM != "" {
		return &VM{
			Provider:   payload.Provider,
			VPC:        payload.VPC,
			Network:    payload.Network,
			Subnetwork: payload.Subnetwork,
			VM:         payload.VM,
		}
	}
	if payload.Provider != "" && (payload.VPC != "" || payload.Network != "") && payload.Subnetwork != "" {
		return &Subnetwork{
			Provider:   payload.Provider,
			VPC:        payload.VPC,
			Network:    payload.Network,
			Subnetwork: payload.Subnetwork,
		}
	}
	if payload.Provider != "" && (payload.VPC != "" || payload.Network != "") && payload.Firewall != "" {
		return &Firewall{
			Provider: payload.Provider,
			VPC:      payload.VPC,
			Network:  payload.Network,
			Firewall: payload.Firewall,
		}
	}
	if payload.Provider != "" && payload.Network != "" {
		return &Network{
			Provider: payload.Provider,
			VPC:      payload.VPC,
			Network:  payload.Network,
		}
	}
	if payload.Provider != "" {
		return &Provider{
			Provider: payload.Provider,
		}
	}

	return nil
}
