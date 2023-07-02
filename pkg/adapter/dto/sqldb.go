package dto

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/instance"
)

type GetSqlDBResponse struct {
	ProjectID string         `json:"projectID"`
	Payload   instance.SqlDB `json:"payload"`
}

type CreateSqlDBRequest struct {
	ParentID       identifier.Network `json:"parent_id"`
	Name           string             `json:"name"`
	Region         string             `json:"region"`
	MachineType    string             `json:"machine_type"`
	SQLType        instance.SQLType   `json:"sql_type"`
	SQLVersion     string             `json:"sql_version"`
	RootPassword   string             `json:"root_password"`
	Disk           instance.SqlDbDisk `json:"sql_db_disk"`
	Subnet1IpRange string             `json:"subnet1_ip_range,omitempty"`
	Subnet2IpRange string             `json:"subnet2_ip_range,omitempty"`
	Managed        bool               `json:"managed" default:"true"`
	Tags           map[string]string  `json:"tags"`
}

type CreateSqlDBResponse struct {
	ProjectID string         `json:"projectID"`
	Payload   instance.SqlDB `json:"payload"`
}

type UpdateSqlDBRequest struct {
	IdentifierID identifier.SqlDB         `json:"identifier_id"`
	Name         *string                  `json:"name"`
	Region       *string                  `json:"region"`
	MachineType  *string                  `json:"machine_type"`
	SQLType      *instance.SQLTypeVersion `json:"sql_type_version"`
	RootPassword *string                  `json:"root_password"`
	Disk         *instance.SqlDbDisk      `json:"sql_db_disk"`
	Managed      *bool                    `json:"managed"`
	Tags         *map[string]string       `json:"tags"`
}

type DeleteSqlDBRequest struct {
	IdentifierID string `json:"identifier_id"`
	Cascade      *bool  `json:"cascade" default:"true"`
}
