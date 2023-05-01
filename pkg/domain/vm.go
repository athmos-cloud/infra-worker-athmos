package domain

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
)

type VM struct {
	Name         string             `json:"name"`
	ProviderType types.ProviderType `json:"providerType"`
	Monitored    bool               `json:"monitored"`
	VPC          string             `json:"vpc"`
	Network      string             `json:"network"`
	Subnetwork   string             `json:"subnetwork"`
	Zone         string             `json:"zone"`
	MachineType  string             `json:"machineType"`
	Auths        []VMAuth           `json:"auths"`
	Disk         Disk               `json:"disk"`
	OS           OS                 `json:"os"`
}

func (vm VM) ToDataMapper(resourceInput resource.IResource) resource.IResource {
	vmInput := resourceInput.(*resource.VM)
	vmInput.Identifier.ID = vm.Name
	vmInput.Identifier.VPCID = vm.VPC
	vmInput.Identifier.NetworkID = vm.Network
	vmInput.Identifier.SubnetID = vm.Subnetwork
	vmInput.Zone = vm.Zone
	vmInput.MachineType = vm.MachineType
	vmInput.Auths = make([]resource.VMAuth, len(vm.Auths))
	for i, auth := range vm.Auths {
		vmInput.Auths[i] = resource.VMAuth{
			Username:     auth.Username,
			SSHPublicKey: auth.SSHPublicKey,
		}
	}
	vmInput.Disk.Type = vm.Disk.Type
	vmInput.Disk.Mode = vm.Disk.Mode
	vmInput.Disk.SizeGib = vm.Disk.SizeGib
	vmInput.OS.Type = vm.OS.Type
	vmInput.OS.Version = vm.OS.Version
	return vmInput
}

type VMCollection map[string]VM

func FromVMDataMapper(vm *resource.VM) VM {
	auths := make([]VMAuth, len(vm.Auths))
	for i, auth := range vm.Auths {
		auths[i] = VMAuth{
			Username:     auth.Username,
			SSHPublicKey: auth.SSHPublicKey,
		}
	}
	return VM{
		Name:         vm.Identifier.ID,
		ProviderType: vm.GetPluginReference().ResourceReference.ProviderType,
		Monitored:    vm.Metadata.Managed,
		VPC:          vm.Identifier.VPCID,
		Network:      vm.Identifier.NetworkID,
		Subnetwork:   vm.Identifier.SubnetID,
		Zone:         vm.Zone,
		MachineType:  vm.MachineType,
		Auths:        auths,
		Disk: Disk{
			Type:       vm.Disk.Type,
			Mode:       vm.Disk.Mode,
			SizeGib:    vm.Disk.SizeGib,
			AutoDelete: vm.Disk.AutoDelete,
		},
	}
}

func FromVMCollectionDataMapper(vmCollection resource.VMCollection) VMCollection {
	vms := VMCollection{}
	for _, vm := range vmCollection {
		vms[vm.Identifier.ID] = FromVMDataMapper(&vm)
	}
	return vms
}

type Disk struct {
	Type       string         `json:"type"`
	Mode       types.DiskMode `json:"mode"`
	SizeGib    int            `json:"sizeGib"`
	AutoDelete bool           `json:"autoDelete"`
}

type VMAuth struct {
	Username     string `json:"username"`
	SSHPublicKey string `json:"sshPublicKey"`
}

type OS struct {
	Type    string `json:"type"`
	Version string `json:"version"`
}
