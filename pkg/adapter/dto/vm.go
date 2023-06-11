package dto

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/instance"
)

type GetVMRequest struct {
	IdentifierID identifier.VM `json:"identifier_id"`
}

type GetVMResponse struct {
	ProjectID string      `json:"project_id"`
	Payload   instance.VM `json:"payload"`
}

type CreateVMRequest struct {
	ParentID       identifier.Subnetwork `json:"parent_id"`
	Name           string                `json:"name"`
	AssignPublicIP bool                  `json:"assign_public_ip" default:"false"`
	Zone           string                `json:"zone"`
	MachineType    string                `json:"machine_type"`
	Auths          []VMAuth              `json:"auths"`
	Disks          instance.VMDiskList   `json:"disks"`
	OS             instance.VMOS         `json:"os"`
	Managed        bool                  `json:"managed" default:"true"`
	Tags           map[string]string     `json:"tags"`
}

type VMAuth struct {
	Username     string `json:"username"`
	RSAKeyLength int    `json:"rsa_key_length"`
}

type CreateVMResponse struct {
	ProjectID string      `json:"project_id"`
	Payload   instance.VM `json:"payload"`
}

type UpdateVMRequest struct {
	IdentifierID   identifier.VM        `json:"identifier_id"`
	Name           *string              `json:"name"`
	AssignPublicIP *bool                `json:"assign_public_ip"`
	Zone           *string              `json:"zone"`
	MachineType    *string              `json:"machine_type"`
	Auths          *[]VMAuth            `json:"auths"`
	UpdateSSHKeys  bool                 `json:"update_ssh_keys" default:"false"`
	Disks          *instance.VMDiskList `json:"disks"`
	OS             *instance.VMOS       `json:"os"`
	Managed        *bool                `json:"managed"`
	Tags           *map[string]string   `json:"tags"`
}

type DeleteVMRequest struct {
	IdentifierID identifier.VM `json:"identifier_id"`
	Cascade      *bool         `json:"cascade" default:"false"`
}
