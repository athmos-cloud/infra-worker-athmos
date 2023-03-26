package vm

import (
	resource2 "github.com/PaulBarrie/infra-worker/pkg/common/dto/resource"
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

func (vm *VM) Create(request resource2.CreateResourceRequest) (resource2.CreateResourceResponse, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (vm *VM) Update(request resource2.UpdateResourceRequest) (resource2.UpdateResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (vm *VM) Get(request resource2.GetResourceRequest) (resource2.GetResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (vm *VM) Watch(request resource2.GetResourceRequest) (resource2.GetResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (vm *VM) List(request resource2.GetListResourceRequest) (resource2.GetListResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (vm *VM) Delete(request resource2.DeleteResourceRequest) (resource2.DeleteResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}
