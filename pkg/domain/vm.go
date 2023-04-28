package domain

import "github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource"

type VM struct {
	Name        string   `json:"name"`
	Monitored   bool     `json:"monitored"`
	VPC         string   `json:"vpc"`
	Network     string   `json:"network"`
	Subnetwork  string   `json:"subnetwork"`
	Zone        string   `json:"zone"`
	MachineType string   `json:"machineType"`
	Auths       []VMAuth `json:"auths"`
	Disks       []Disk   `json:"disks"`
	OS          OS       `json:"os"`
}

type VMCollection map[string]VM

func FromVMDataMapper(vm resource.VM) VM {
	auths := make([]VMAuth, len(vm.Auths))
	for i, auth := range vm.Auths {
		auths[i] = VMAuth{
			Username:     auth.Username,
			SSHPublicKey: auth.SSHPublicKey,
		}
	}
	disks := make([]Disk, len(vm.Disks))
	for i, disk := range vm.Disks {
		disks[i] = Disk{
			Type:       disk.Type,
			Mode:       DiskModeFromString(string(disk.Mode)),
			SizeGib:    disk.SizeGib,
			AutoDelete: disk.AutoDelete,
		}
	}
	return VM{
		Name:        vm.Identifier.ID,
		Monitored:   vm.Metadata.Monitored,
		VPC:         vm.VPC,
		Network:     vm.Network,
		Subnetwork:  vm.Subnetwork,
		Zone:        vm.Zone,
		MachineType: vm.MachineType,
		Auths:       auths,
		Disks:       disks,
	}
}

func FromVMCollectionDataMapper(vmCollection resource.VMCollection) VMCollection {
	vms := VMCollection{}
	for _, vm := range vmCollection {
		vms[vm.Identifier.ID] = FromVMDataMapper(vm)
	}
	return vms
}

type Disk struct {
	Type       string   `json:"type"`
	Mode       DiskMode `json:"mode"`
	SizeGib    int      `json:"sizeGib"`
	AutoDelete bool     `json:"autoDelete"`
}

type DiskMode string

const (
	ReadOnly  DiskMode = "READ_ONLY"
	ReadWrite DiskMode = "READ_WRITE"
)

func DiskModeFromString(input string) DiskMode {
	switch input {
	case "READ_ONLY":
		return ReadOnly
	case "READ_WRITE":
		return ReadWrite
	default:
		return ""
	}
}

type VMAuth struct {
	Username     string `json:"username"`
	SSHPublicKey string `json:"sshPublicKey"`
}

type OS struct {
	Type    string `json:"type"`
	Version string `json:"version"`
}
