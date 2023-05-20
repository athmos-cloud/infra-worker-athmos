package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
)

type VM struct {
	Metadata       metadata.Metadata `json:"metadata"`
	IdentifierID   identifier.VM     `json:"identifierID"`
	IdentifierName identifier.VM     `json:"identifierName"`
	PublicIP       bool              `json:"publicIP,omitempty"`
	Zone           string            `json:"zone"`
	MachineType    string            `json:"machineType"`
	Auths          VMAuthList        `json:"auths,omitempty"`
	Disks          DiskList          `json:"disks"`
	OS             OS                `json:"os"`
}

type VMCollection map[string]VM

type Disk struct {
	Type       string         `json:"type"`
	Mode       types.DiskMode `json:"mode"`
	SizeGib    int            `json:"sizeGib"`
	AutoDelete bool           `json:"autoDelete"`
}

func (vm *VM) Equals(other VM) bool {
	return vm.Metadata.Equals(other.Metadata) &&
		vm.IdentifierID.Equals(&other.IdentifierID) &&
		vm.Zone == other.Zone &&
		vm.MachineType == other.MachineType &&
		vm.Auths.Equals(other.Auths) &&
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

func (disk *Disk) Equals(other Disk) bool {
	return disk.Type == other.Type && disk.SizeGib == other.SizeGib && disk.AutoDelete == other.AutoDelete
}

type DiskList []Disk

func (diskList *DiskList) Equals(other DiskList) bool {
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

type VMAuth struct {
	Username     string
	SSHPublicKey string
}

func (auth *VMAuth) Equals(other VMAuth) bool {
	return auth.Username == other.Username && auth.SSHPublicKey == other.SSHPublicKey
}

type VMAuthList []VMAuth

func (authList *VMAuthList) Equals(other VMAuthList) bool {
	if len(*authList) != len(other) {
		return false
	}
	for i, auth := range *authList {
		if !auth.Equals(other[i]) {
			return false
		}
	}
	return true
}

type OS struct {
	Type    string `json:"type"`
	Version string `json:"version"`
}

func (os *OS) Equals(other OS) bool {
	return os.Type == other.Type && os.Version == other.Version
}
