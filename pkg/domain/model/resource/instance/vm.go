package instance

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/metadata"
)

type VM struct {
	Metadata       metadata.Metadata `json:"metadata"`
	IdentifierID   identifier.VM     `json:"identifier_id"`
	IdentifierName identifier.VM     `json:"identifier_name"`
	AssignPublicIP bool              `json:"assign_public_ip"`
	PublicIP       string            `json:"public_ip,omitempty"`
	Zone           string            `json:"zone"`
	MachineType    string            `json:"machine_type"`
	Auths          model.SSHKeyList  `json:"auths,omitempty"`
	Disks          VMDiskList        `json:"disks"`
	OS             VMOS              `json:"os"`
}

type VMCollection map[string]VM

type VMOS struct {
	Name string `json:"name"`
	ID   string `json:"id,omitempty"`
}
