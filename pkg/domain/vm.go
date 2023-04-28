package domain

import "github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource"

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
	Type       string   `bson:"type"`
	Mode       DiskMode `bson:"mode"`
	SizeGib    int      `bson:"sizeGib"`
	AutoDelete bool     `bson:"autoDelete"`
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
	Username     string `bson:"username"`
	SSHPublicKey string `bson:"sshPublicKey"`
}

type OS struct {
	Type    string `bson:"type"`
	Version string `bson:"version"`
}
