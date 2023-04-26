package resource

import (
	"fmt"
	dto "github.com/athmos-cloud/infra-worker-athmos/pkg/common/dto/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

type VM struct {
	Metadata    metadata.Metadata `bson:"metadata"`
	Identifier  identifier.VM     `bson:"identifier"`
	VPC         string            `bson:"vpc"`
	Network     string            `bson:"network"`
	Subnetwork  string            `bson:"subnetwork"`
	Zone        string            `bson:"zone"`
	MachineType string            `bson:"machineType"`
	Auths       []VMAuth          `bson:"auths"`
	Disks       []Disk            `bson:"disks"`
	OS          OS                `bson:"os"`
}

func NewVM(id identifier.VM) VM {
	return VM{
		Metadata:   metadata.New(metadata.CreateMetadataRequest{Name: id.ID}),
		Identifier: id,
	}
}

type VMCollection map[string]VM

type Disk struct {
	Type       string   `bson:"type"`
	Mode       DiskMode `bson:"mode"`
	SizeGib    int      `bson:"sizeGib"`
	AutoDelete bool     `bson:"autoDelete"`
}

type DiskMode string

const (
	READ_ONLY  DiskMode = "READ_ONLY"
	READ_WRITE DiskMode = "READ_WRITE"
)

type VMAuth struct {
	Username     string `bson:"username"`
	SSHPublicKey string `bson:"sshPublicKey"`
}

type OS struct {
	Type    string `bson:"type"`
	Version string `bson:"version"`
}

func (vm *VM) GetMetadata() metadata.Metadata {
	return vm.Metadata
}

func (vm *VM) WithMetadata(request metadata.CreateMetadataRequest) {
	vm.Metadata = metadata.New(request)
}

func (vm *VM) GetPluginReference(request dto.GetPluginReferenceRequest) (dto.GetPluginReferenceResponse, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (vm *VM) FromMap(data map[string]interface{}) errors.Error {
	//TODO implement me
	panic("implement me")
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

func (vm *VM) ToDomain() (interface{}, errors.Error) {
	//TODO implement me
	panic("implement me")
}
