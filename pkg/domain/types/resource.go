package types

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

type Resource string

const (
	ProviderResource   Resource = "provider"
	VPCResource        Resource = "vpc"
	NetworkResource    Resource = "network"
	SubnetworkResource Resource = "subnetwork"
	FirewallResource   Resource = "firewall"
	VMResource         Resource = "vm"
	SqlDBResource      Resource = "sqldb"
)

var resourcesMapping = map[string]Resource{
	"provider":   ProviderResource,
	"vpc":        VPCResource,
	"network":    NetworkResource,
	"subnetwork": SubnetworkResource,
	"firewall":   FirewallResource,
	"vm":         VMResource,
	"sqldb":      SqlDBResource,
}

func ResourceFromString(s string) (Resource, errors.Error) {
	if val, ok := resourcesMapping[s]; ok {
		return val, errors.OK
	}
	return "", errors.BadRequest.WithMessage(fmt.Sprintf("resource %s is not supported", s))
}
