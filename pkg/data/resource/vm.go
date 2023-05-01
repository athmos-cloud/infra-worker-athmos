package resource

import (
	"fmt"
	resourcePlugin "github.com/athmos-cloud/infra-worker-athmos/pkg/data/plugin"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/types"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/status"
	commonTypes "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/config"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"reflect"
)

type VM struct {
	Metadata    metadata.Metadata     `bson:"metadata"`
	Identifier  identifier.VM         `bson:"identifier" plugin:"identifier"`
	Status      status.ResourceStatus `bson:"status"`
	VPC         string                `bson:"vpc" plugin:"vpc"`
	Provider    string                `bson:"provider" plugin:"providerName"`
	Network     string                `bson:"network" plugin:"network"`
	Subnetwork  string                `bson:"vmwork" plugin:"subnet"`
	PublicIP    bool                  `bson:"publicIP" plugin:"publicIP"`
	Zone        string                `bson:"zone" plugin:"zone"`
	MachineType string                `bson:"machineType" plugin:"machineType"`
	Auths       VMAuthList            `bson:"auths" plugin:"auth"`
	Disk        Disk                  `bson:"disk" plugin:"disk"`
	OS          OS                    `bson:"os" plugin:"os"`
}

func NewVM(id identifier.VM, providerType commonTypes.ProviderType) VM {
	return VM{
		Metadata:   metadata.New(metadata.CreateMetadataRequest{Name: id.ID}),
		Identifier: id,
		Status:     status.New(id.ID, commonTypes.VM, providerType),
		Auths:      make(VMAuthList, 0),
		Disk:       Disk{},
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
	Type       string         `bson:"type" plugin:"type"`
	Mode       types.DiskMode `bson:"mode" plugin:"diskMode"`
	SizeGib    int            `bson:"sizeGib" plugin:"sizeGib"`
	AutoDelete bool           `bson:"autoDelete" plugin:"autoDelete"`
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
	Username     string `bson:"username" plugin:"user"`
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
	Type    string `bson:"type" plugin:"osType"`
	Version string `bson:"version" plugin:"version"`
}

func (os *OS) Equals(other OS) bool {
	return os.Type == other.Type && os.Version == other.Version
}

func (vm *VM) New(id identifier.ID, providerType commonTypes.ProviderType) IResource {
	if reflect.TypeOf(id) != reflect.TypeOf(identifier.VM{}) {
		panic(errors.InvalidArgument.WithMessage("identifier is not of type VM"))
	}
	res := NewVM(id.(identifier.VM), providerType)
	return &res
}

func (vm *VM) Equals(other VM) bool {
	return vm.Metadata.Equals(other.Metadata) &&
		vm.Identifier.Equals(other.Identifier) &&
		vm.Status.Equals(other.Status) &&
		vm.VPC == other.VPC &&
		vm.Network == other.Network &&
		vm.Subnetwork == other.Subnetwork &&
		vm.Zone == other.Zone &&
		vm.MachineType == other.MachineType &&
		vm.Auths.Equals(other.Auths) &&
		vm.Disk.Equals(other.Disk) &&
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
	if err := resourcePlugin.InjectMapIntoStruct(data, vm); err.IsOk() {
		panic(err)
	}
}

func (vm *VM) Insert(project Project, update ...bool) {
	shouldUpdate := false
	if len(update) > 0 {
		shouldUpdate = update[0]
	}
	id := vm.Identifier
	_, ok := project.Resources[id.ProviderID].VPCs[id.VPCID].Networks[id.NetworkID].Subnetworks[id.SubnetID].VMs[id.ID]
	if !ok && shouldUpdate {
		panic(errors.NotFound.WithMessage(fmt.Sprintf("vm %s not found in vm %s", id.ID, id.SubnetID)))
	}
	if ok && !shouldUpdate {
		panic(errors.Conflict.WithMessage(fmt.Sprintf("vm %s already exists in vm %s", id.ID, id.SubnetID)))
	}
	project.Resources[id.ProviderID].VPCs[id.VPCID].Networks[id.NetworkID].Subnetworks[id.SubnetID].VMs[id.ID] = *vm
}

func (vm *VM) Remove(project Project) {
	id := vm.Identifier
	_, ok := project.Resources[id.ProviderID].VPCs[id.VPCID].Networks[id.NetworkID].Subnetworks[id.SubnetID].VMs[id.ID]
	if !ok {
		panic(errors.NotFound.WithMessage(fmt.Sprintf("vm %s not found in vm %s", id.ID, id.SubnetID)))
	}
	delete(project.Resources[id.ProviderID].VPCs[id.VPCID].Networks[id.NetworkID].Subnetworks[id.SubnetID].VMs, id.ID)
}
