package resource

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/common"
	dto "github.com/athmos-cloud/infra-worker-athmos/pkg/common/dto/resource"
	resourcePlugin "github.com/athmos-cloud/infra-worker-athmos/pkg/data/plugin"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/types"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/config"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"reflect"
)

type VM struct {
	Metadata    metadata.Metadata `bson:"metadata"`
	Identifier  identifier.VM     `bson:"identifier"`
	VPC         string            `bson:"vpc" plugin:"vpc"`
	Network     string            `bson:"network" plugin:"network"`
	Subnetwork  string            `bson:"subnetwork" plugin:"subnetwork"`
	Zone        string            `bson:"zone" plugin:"zone"`
	MachineType string            `bson:"machineType" plugin:"machineType"`
	Auths       VMAuthList        `bson:"auths" plugin:"auths"`
	Disks       DiskList          `bson:"disks" plugin:"disks"`
	OS          OS                `bson:"os" plugin:"os"`
}

func NewVM(id identifier.VM) VM {
	return VM{
		Metadata:   metadata.New(metadata.CreateMetadataRequest{Name: id.ID}),
		Identifier: id,
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
	Type    string `bson:"type" plugin:"osType"`
	Version string `bson:"version" plugin:"version"`
}

func (os *OS) Equals(other OS) bool {
	return os.Type == other.Type && os.Version == other.Version
}

func (vm *VM) New(id identifier.ID) (IResource, errors.Error) {
	if reflect.TypeOf(id) != reflect.TypeOf(identifier.VM{}) {
		return nil, errors.InvalidArgument.WithMessage("identifier is not of type VM")
	}
	res := NewVM(id.(identifier.VM))
	return &res, errors.OK
}

func (vm *VM) Equals(other VM) bool {
	return vm.Metadata.Equals(other.Metadata) &&
		vm.Identifier.Equals(other.Identifier) &&
		vm.VPC == other.VPC &&
		vm.Network == other.Network &&
		vm.Subnetwork == other.Subnetwork &&
		vm.Zone == other.Zone &&
		vm.MachineType == other.MachineType &&
		vm.Auths.Equals(other.Auths) &&
		vm.Disks.Equals(other.Disks) &&
		vm.OS.Equals(other.OS)
}

func (vm *VM) GetMetadata() metadata.Metadata {
	return vm.Metadata
}

func (vm *VM) WithMetadata(request metadata.CreateMetadataRequest) {
	vm.Metadata = metadata.New(request)
}

func (vm *VM) GetPluginReference(request dto.GetPluginReferenceRequest) (dto.GetPluginReferenceResponse, errors.Error) {
	switch request.ProviderType {
	case common.GCP:
		return dto.GetPluginReferenceResponse{
			ChartName:    config.Current.Plugins.Crossplane.GCP.VM.Chart,
			ChartVersion: config.Current.Plugins.Crossplane.GCP.VM.Version,
		}, errors.Error{}
	}
	return dto.GetPluginReferenceResponse{}, errors.InvalidArgument.WithMessage(fmt.Sprintf("provider type %s not supported", request.ProviderType))
}

func (vm *VM) FromMap(data map[string]interface{}) errors.Error {
	return resourcePlugin.InjectMapIntoStruct(data, vm)
}

func (vm *VM) Insert(project Project, update ...bool) errors.Error {
	shouldUpdate := false
	if len(update) > 0 {
		shouldUpdate = update[0]
	}
	id := vm.Identifier
	_, ok := project.Resources[id.ProviderID].VPCs[id.VPCID].Networks[id.NetworkID].Subnetworks[id.SubnetID].VMs[id.ID]
	if !ok && shouldUpdate {
		return errors.NotFound.WithMessage(fmt.Sprintf("vm %s not found in subnet %s", id.ID, id.SubnetID))
	}
	if ok && !shouldUpdate {
		return errors.Conflict.WithMessage(fmt.Sprintf("vm %s already exists in subnet %s", id.ID, id.SubnetID))
	}
	project.Resources[id.ProviderID].VPCs[id.VPCID].Networks[id.NetworkID].Subnetworks[id.SubnetID].VMs[id.ID] = *vm
	return errors.OK
}

func (vm *VM) Remove(project Project) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (vm *VM) ToDomain() (interface{}, errors.Error) {
	//TODO implement me
	panic("implement me")
}
