package identifier

import (
	"encoding/json"
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"reflect"
)

const (
	providerIDKey = "providerID"
	vpcIDKey      = "vpcID"
	networkIDKey  = "networkID"
	subnetIDKey   = "subnetID"
	vmIDKey       = "vmID"
	firewallIDKey = "firewallID"
)

func ToMap(id ID) map[string]interface{} {
	switch reflect.TypeOf(id) {
	case reflect.TypeOf(Provider{}):
		provider := id.(Provider)
		return map[string]interface{}{
			providerIDKey: provider.ProviderID,
		}
	case reflect.TypeOf(VPC{}):
		vpc := id.(VPC)
		return map[string]interface{}{
			providerIDKey: vpc.ProviderID,
			vpcIDKey:      vpc.VPCID,
		}
	case reflect.TypeOf(Network{}):
		network := id.(Network)
		return map[string]interface{}{
			providerIDKey: network.ProviderID,
			vpcIDKey:      network.VPCID,
			networkIDKey:  network.NetworkID,
		}
	case reflect.TypeOf(Subnetwork{}):
		subnetwork := id.(Subnetwork)
		return map[string]interface{}{
			providerIDKey: subnetwork.ProviderID,
			vpcIDKey:      subnetwork.VPCID,
			networkIDKey:  subnetwork.NetworkID,
			"subnetID":    subnetwork.SubnetworkID,
		}
	case reflect.TypeOf(Firewall{}):
		firewall := id.(Firewall)
		return map[string]interface{}{
			providerIDKey: firewall.ProviderID,
			vpcIDKey:      firewall.VPCID,
			networkIDKey:  firewall.NetworkID,
			firewallIDKey: firewall.FirewallID,
		}
	case reflect.TypeOf(VM{}):
		vm := id.(VM)
		return map[string]interface{}{
			providerIDKey: vm.ProviderID,
			vpcIDKey:      vm.VPCID,
			networkIDKey:  vm.NetworkID,
			subnetIDKey:   vm.SubnetID,
			vmIDKey:       vm.VMID,
		}
	}
	panic(errors.InternalError.WithMessage(fmt.Sprintf("Unknown ID type: %v", id)))
}

func FromMap(payload map[string]interface{}) ID {
	body, err := json.Marshal(payload)
	if err != nil {
		panic(errors.InternalError.WithMessage(fmt.Sprintf("Error while marshalling payload: %s", err.Error())))
	}
	var idPayload IdPayload
	if errUnmarshal := json.Unmarshal(body, &idPayload); errUnmarshal != nil {
		panic(errors.InvalidArgument.WithMessage(fmt.Sprintf("%v is not a valid FirewallID payload", payload)))
	}
	return FromPayload(idPayload)
}

func ToPayload(input ID) IdPayload {
	switch reflect.TypeOf(input) {
	case reflect.TypeOf(Provider{}):
		id := input.(Provider)
		return IdPayload{
			ProviderID: id.ProviderID,
		}
	case reflect.TypeOf(VPC{}):
		id := input.(VPC)
		return IdPayload{
			ProviderID: id.ProviderID,
			VPCID:      id.VPCID,
		}
	case reflect.TypeOf(Network{}):
		id := input.(Network)
		return IdPayload{
			ProviderID: id.ProviderID,
			VPCID:      id.VPCID,
			NetworkID:  id.NetworkID,
		}
	case reflect.TypeOf(Subnetwork{}):
		id := input.(Subnetwork)
		return IdPayload{
			ProviderID: id.ProviderID,
			VPCID:      id.VPCID,
			NetworkID:  id.NetworkID,
			SubnetID:   id.SubnetworkID,
		}
	case reflect.TypeOf(Firewall{}):
		id := input.(Firewall)
		return IdPayload{
			ProviderID: id.ProviderID,
			VPCID:      id.VPCID,
			NetworkID:  id.NetworkID,
			FirewallID: id.FirewallID,
		}
	case reflect.TypeOf(VM{}):
		id := input.(VM)
		return IdPayload{
			ProviderID: id.ProviderID,
			VPCID:      id.VPCID,
			NetworkID:  id.NetworkID,
			SubnetID:   id.SubnetID,
			VMID:       id.VMID,
		}
	}
	return IdPayload{}
}

func FromPayload(payload IdPayload) ID {
	if payload.VMID != "" && payload.ProviderID != "" && payload.VPCID != "" && payload.NetworkID != "" && payload.SubnetID != "" {
		return VM{
			VMID:       payload.VMID,
			ProviderID: payload.ProviderID,
			VPCID:      payload.VPCID,
			NetworkID:  payload.NetworkID,
			SubnetID:   payload.SubnetID,
		}
	} else if payload.FirewallID != "" && payload.ProviderID != "" && payload.VPCID != "" && payload.NetworkID != "" {
		return Firewall{
			FirewallID: payload.FirewallID,
			ProviderID: payload.ProviderID,
			VPCID:      payload.VPCID,
			NetworkID:  payload.NetworkID,
		}
	} else if payload.ProviderID != "" && payload.NetworkID != "" && payload.SubnetID != "" {
		return Subnetwork{
			SubnetworkID: payload.SubnetID,
			ProviderID:   payload.ProviderID,
			VPCID:        payload.VPCID,
			NetworkID:    payload.NetworkID,
		}
	} else if payload.ProviderID != "" && payload.NetworkID != "" {
		return Network{
			NetworkID:  payload.NetworkID,
			ProviderID: payload.ProviderID,
			VPCID:      payload.VPCID,
		}
	} else if payload.ProviderID != "" && payload.VPCID != "" {
		return VPC{
			VPCID:      payload.VPCID,
			ProviderID: payload.ProviderID,
		}
	} else if payload.ProviderID != "" {
		return Provider{
			ProviderID: payload.ProviderID,
		}
	} else {
		return Empty{}
	}
}
