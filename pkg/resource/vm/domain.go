package vm

import (
	"github.com/PaulBarrie/infra-worker/pkg/application/dto"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
	"github.com/PaulBarrie/infra-worker/pkg/resource"
)

type VM struct {
	ID                string `bson:"_id,omitempty"`
	ResourceReference resource.Reference
	VPC               string   `bson:"vpc"`
	Network           string   `bson:"network"`
	Subnetwork        string   `bson:"subnetwork"`
	Zone              string   `bson:"zone"`
	MachineType       string   `bson:"machineType"`
	Auths             []VMAuth `bson:"auths"`
	Disks             []Disk   `bson:"disks"`
	OS                OS       `bson:"os"`
}

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

func (vm *VM) Create(request dto.CreateResourceRequest) (dto.CreateResourceResponse, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (vm *VM) Update(request dto.UpdateResourceRequest) (dto.UpdateResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (vm *VM) Get(request dto.GetResourceRequest) (dto.GetResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (vm *VM) Watch(request dto.GetResourceRequest) (dto.GetResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (vm *VM) List(request dto.GetListResourceRequest) (dto.GetListResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (vm *VM) Delete(request dto.DeleteResourceRequest) (dto.DeleteResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}
