package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/metadata"
)

type VM struct {
	Metadata       metadata.Metadata `json:"metadata"`
	IdentifierID   identifier.VM     `json:"identifierID"`
	IdentifierName identifier.VM     `json:"identifierName"`
	AssignPublicIP bool              `json:"assignPublicIP"`
	PublicIP       string            `json:"publicIP,omitempty"`
	Zone           string            `json:"zone"`
	MachineType    string            `json:"machineType"`
	Auths          model.SSHKeyList  `json:"auths,omitempty"`
	Disks          VMDiskList        `json:"disks"`
	OS             VMOS              `json:"os"`
}

type VMCollection map[string]VM

type VMDisk struct {
	Type       DiskType `json:"type"`
	Mode       DiskMode `json:"mode"`
	SizeGib    int      `json:"sizeGib"`
	AutoDelete bool     `json:"autoDelete"`
}

type VMDiskList []VMDisk

type VMOS struct {
	Name string `json:"name"`
	ID   string `json:"id,omitempty"`
}

func (vm *VM) Equals(other VM) bool {
	return vm.Metadata.Equals(other.Metadata) &&
		vm.IdentifierID.Equals(&other.IdentifierID) &&
		vm.IdentifierName.Equals(&other.IdentifierName) &&
		vm.AssignPublicIP == other.AssignPublicIP &&
		vm.PublicIP == other.PublicIP &&
		vm.Zone == other.Zone &&
		vm.MachineType == other.MachineType &&
		vm.Disks.Equals(other.Disks) &&
		vm.OS.Equals(other.OS)
}

func (collection *VMCollection) Equals(other VMCollection) bool {
	if len(*collection) != len(other) {
		return false
	}
	for key, value := range *collection {
		if !value.Equals(other[key]) {
			return false
		}
	}
	return true
}

func (disk *VMDisk) Equals(other VMDisk) bool {
	return disk.Type == other.Type && disk.SizeGib == other.SizeGib && disk.AutoDelete == other.AutoDelete
}

func (diskList *VMDiskList) Equals(other VMDiskList) bool {
	if len(*diskList) != len(other) {
		return false
	}
	for i, disk := range *diskList {
		if !disk.Equals(other[i]) {
			return false
		}
	}
	return true
}

func (os *VMOS) Equals(other VMOS) bool {
	return os.Name == other.Name && os.ID == other.ID
}
