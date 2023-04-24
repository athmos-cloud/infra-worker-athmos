package resources

import (
	dto "github.com/athmos-cloud/infra-worker-athmos/pkg/common/dto/resource"
	domain2 "github.com/athmos-cloud/infra-worker-athmos/pkg/data/project"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

type VM struct {
	Metadata    domain.Metadata `bson:"metadata"`
	VPC         string          `bson:"vpc"`
	Network     string          `bson:"network"`
	Subnetwork  string          `bson:"subnetwork"`
	Zone        string          `bson:"zone"`
	MachineType string          `bson:"machineType"`
	Auths       []VMAuth        `bson:"auths"`
	Disks       []Disk          `bson:"disks"`
	OS          OS              `bson:"os"`
}

type VMHierarchyLocation struct {
	ProviderID string
	VPCID      string
	NetworkID  string
	SubnetID   string
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

func (vm *VM) GetMetadata() domain.Metadata {
	return vm.Metadata
}

func (vm *VM) WithMetadata(request domain.CreateMetadataRequest) {
	vm.Metadata = domain.New(request)
}

func (vm *VM) GetPluginReference(request dto.GetPluginReferenceRequest) (dto.GetPluginReferenceResponse, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (vm *VM) FromMap(data map[string]interface{}) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (vm *VM) InsertIntoProject(project domain2.Project, upsert bool) errors.Error {
	//TODO implement me
	panic("implement me")
}
