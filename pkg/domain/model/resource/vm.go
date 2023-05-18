package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/metadata"
	"github.com/kamva/mgm/v3"
)

type VM struct {
	mgm.DefaultModel `bson:",inline"`
	Metadata         metadata.Metadata `bson:"metadata"`
	Identifier       identifier.VM     `bson:"identifier" plugin:"identifier"`
	PublicIP         bool              `bson:"publicIP" plugin:"publicIP"`
	Zone             string            `bson:"zone" plugin:"zone"`
	MachineType      string            `bson:"machineType" plugin:"machineType"`
	Auths            VMAuthList        `bson:"auths" plugin:"auths"`
	Disks            DiskList          `bson:"disk" plugin:"disks" yaml:"disk"`
	OS               OS                `bson:"os" plugin:"os"`
}

type VMCollection map[string]VM

type Disk struct {
	Type string `bson:"type" plugin:"type"`
	//Mode       types2.DiskMode `bson:"mode" plugin:"diskMode"`
	SizeGib    int  `bson:"sizeGib" plugin:"sizeGib"`
	AutoDelete bool `bson:"autoDelete" plugin:"autoDelete"`
}

type OS struct {
	Type    string `bson:"type" plugin:"osType"`
	Version string `bson:"version" plugin:"version"`
}

func NewVM(payload NewResourcePayload) VM {
	payload.Validate()
	parentID := payload.ParentIdentifier.(identifier.Subnetwork)
	id := identifier.VM{
		ProviderID: parentID.ProviderID,
		VPCID:      parentID.ProviderID,
		NetworkID:  parentID.NetworkID,
		SubnetID:   parentID.SubnetworkID,
		VMID:       formatResourceName(payload.Name),
	}
	return VM{
		Metadata: metadata.New(metadata.CreateMetadataRequest{
			Name:         id.VMID,
			NotMonitored: !payload.Managed,
			Tags:         payload.Tags,
		}),
		Identifier: id,
		Auths:      make(VMAuthList, 0),
		Disks:      make(DiskList, 0),
	}
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
	Username     string `bson:"username" plugin:"username"`
	SSHPublicKey string `bson:"sshPublicKey" plugin:"sshPublicKey"`
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

func (os *OS) Equals(other OS) bool {
	return os.Type == other.Type && os.Version == other.Version
}

func (vm *VM) Equals(other VM) bool {
	return vm.Metadata.Equals(other.Metadata) &&
		vm.Identifier.Equals(other.Identifier) &&
		vm.Zone == other.Zone &&
		vm.MachineType == other.MachineType &&
		vm.Auths.Equals(other.Auths) &&
		vm.Disks.Equals(other.Disks) &&
		vm.OS.Equals(other.OS)
}
