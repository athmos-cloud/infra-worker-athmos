package resource

import (
	"fmt"
	resourcePlugin "github.com/athmos-cloud/infra-worker-athmos/pkg/data/plugin"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/status"
	commonTypes "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/config"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/utils"
	"github.com/kamva/mgm/v3"
	"reflect"
)

func NewVM(payload NewResourcePayload) VM {
	payload.Validate()
	if reflect.TypeOf(payload.ParentIdentifier) != reflect.TypeOf(identifier.Subnetwork{}) {
		panic(errors.InvalidArgument.WithMessage(fmt.Sprintf("Invalid parent identifier type: %s", reflect.TypeOf(payload.ParentIdentifier))))
	}
	parentID := payload.ParentIdentifier.(identifier.Subnetwork)
	id := identifier.VM{
		ProviderID: parentID.ProviderID,
		VPCID:      parentID.ProviderID,
		NetworkID:  parentID.NetworkID,
		SubnetID:   parentID.SubnetworkID,
		VMID:       fmt.Sprintf("%s-%s", payload.Name, utils.RandomString(resourceIDSuffixLength)),
	}
	return VM{
		Metadata: metadata.New(metadata.CreateMetadataRequest{
			Name:         id.VMID,
			NotMonitored: !payload.Managed,
			Tags:         payload.Tags,
		}),
		Identifier: id,
		Status:     status.New(id.VMID, commonTypes.VM, payload.Provider),
		Auths:      make(VMAuthList, 0),
		Disks:      DiskList{},
	}
}

type VMCollection map[string]VM

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

type Disk struct {
	mgm.DefaultModel `bson:",inline"`
	Type             string               `bson:"type" plugin:"type"`
	Mode             commonTypes.DiskMode `bson:"mode" plugin:"diskMode"`
	SizeGib          int                  `bson:"sizeGib" plugin:"sizeGib"`
	AutoDelete       bool                 `bson:"autoDelete" plugin:"autoDelete"`
}

func (disk *Disk) Equals(other Disk) bool {
	return disk.Type == other.Type && disk.Mode == other.Mode && disk.SizeGib == other.SizeGib && disk.AutoDelete == other.AutoDelete
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

type OS struct {
	mgm.DefaultModel `bson:",inline"`
	Type             string `bson:"type" plugin:"osType"`
	Version          string `bson:"version" plugin:"version"`
}

func (os *OS) Equals(other OS) bool {
	return os.Type == other.Type && os.Version == other.Version
}

type VM struct {
	mgm.DefaultModel `bson:",inline"`
	Metadata         metadata.Metadata     `bson:"metadata"`
	Identifier       identifier.VM         `bson:"identifier" plugin:"identifier"`
	Status           status.ResourceStatus `bson:"status"`
	PublicIP         bool                  `bson:"publicIP" plugin:"publicIP"`
	Zone             string                `bson:"zone" plugin:"zone"`
	MachineType      string                `bson:"machineType" plugin:"machineType"`
	Auths            VMAuthList            `bson:"auths" plugin:"auths"`
	Disks            DiskList              `bson:"disk" plugin:"disk"`
	OS               OS                    `bson:"os" plugin:"os"`
}

func (vm *VM) GetIdentifier() identifier.ID {
	return vm.Identifier
}

func (vm *VM) New(payload NewResourcePayload) IResource {
	if reflect.TypeOf(payload.ParentIdentifier) != reflect.TypeOf(identifier.VM{}) {
		panic(errors.InvalidArgument.WithMessage("identifier is not of type VM"))
	}
	res := NewVM(payload)
	return &res
}

func (vm *VM) Equals(other VM) bool {
	return vm.Metadata.Equals(other.Metadata) &&
		vm.Identifier.Equals(other.Identifier) &&
		vm.Status.Equals(other.Status) &&
		vm.Zone == other.Zone &&
		vm.MachineType == other.MachineType &&
		vm.Auths.Equals(other.Auths) &&
		vm.Disks.Equals(other.Disks) &&
		vm.OS.Equals(other.OS)
}

func (vm *VM) GetMetadata() metadata.Metadata {
	return vm.Metadata
}

func (vm *VM) SetMetadata(request metadata.CreateMetadataRequest) {
	vm.Metadata = metadata.New(request)
}

func (vm *VM) SetStatus(resourceStatus status.ResourceStatus) {
	vm.Status = resourceStatus
}

func (vm *VM) GetStatus() status.ResourceStatus {
	return vm.Status
}

func (vm *VM) GetPluginReference() resourcePlugin.Reference {
	if !vm.Status.PluginReference.ChartReference.Empty() {
		return vm.Status.PluginReference
	}
	switch vm.Status.PluginReference.ResourceReference.ProviderType {
	case commonTypes.GCP:
		vm.Status.PluginReference.ChartReference = resourcePlugin.HelmChartReference{
			ChartName:    config.Current.Plugins.Crossplane.GCP.Subnet.Chart,
			ChartVersion: config.Current.Plugins.Crossplane.GCP.Subnet.Version,
		}
		return vm.Status.PluginReference
	}
	panic(errors.InvalidArgument.WithMessage(fmt.Sprintf("provider type %s not supported", vm.Status.PluginReference.ResourceReference.ProviderType)))
}

func (vm *VM) FromMap(data map[string]interface{}) {
	err := resourcePlugin.InjectMapIntoStruct(data, vm)
	if !err.IsOk() {
		panic(err)
	}
}

func (vm *VM) Insert(_ IResource, _ ...bool) {
	return
}

func (vm *VM) Remove(_ IResource) {
	return
}
