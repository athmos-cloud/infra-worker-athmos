package dto

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
)

type GetVMRequest struct {
	IdentifierID identifier.VM `json:"identifierID"`
}

type GetVMResponse struct {
	ProjectID string      `json:"projectID"`
	Payload   resource.VM `json:"payload"`
}

type CreateVMRequest struct {
	ParentID       identifier.Subnetwork `json:"parentID"`
	Name           string                `json:"name"`
	AssignPublicIP bool                  `json:"assignPublicIP" default:"false"`
	Zone           string                `json:"zone"`
	MachineType    string                `json:"machineType"`
	Auths          resource.VMAuthList   `json:"auths"`
	Disks          resource.VMDiskList   `json:"disks"`
	OS             resource.VMOS         `json:"os"`
	Managed        bool                  `json:"managed" default:"true"`
	Tags           map[string]string     `json:"tags"`
}

type CreateVMResponse struct {
	ProjectID string      `json:"projectID"`
	Payload   resource.VM `json:"payload"`
}

type UpdateVMRequest struct {
	IdentifierID   identifier.VM        `json:"identifierID"`
	Name           *string              `json:"name"`
	AssignPublicIP *bool                `json:"assignPublicIP"`
	Zone           *string              `json:"zone"`
	MachineType    *string              `json:"machineType"`
	Auths          *resource.VMAuthList `json:"auths"`
	Disks          *resource.VMDiskList `json:"disks"`
	OS             *resource.VMOS       `json:"os"`
	Managed        *bool                `json:"managed"`
	Tags           *map[string]string   `json:"tags"`
}

type DeleteVMRequest struct {
	IdentifierID identifier.VM `json:"identifierID"`
	Cascade      *bool         `json:"cascade" default:"false"`
}
