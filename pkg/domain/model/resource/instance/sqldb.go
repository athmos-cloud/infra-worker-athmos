package instance

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/metadata"
)

type SqlDB struct {
	Metadata       metadata.Metadata `json:"metadata"`
	IdentifierID   identifier.SqlDB  `json:"identifier_id"`
	IdentifierName identifier.SqlDB  `json:"identifier_name"`
	PublicIP       string            `json:"public_ip,omitempty"`
	MachineType    string            `json:"machine_type"`
	SQLTypeVersion SQLTypeVersion    `json:"sql_version"`
	PrivateIp      string            `json:"private_ip,omitempty"`
	Region         string            `json:"region"`
	Auth           SqlDBAuth         `json:"sql_db_auth"`
	Disk           SqlDbDisk         `json:"sql_db_disk"`
	Subnet1IpRange string            `json:"subnet1_ip_range,omitempty"`
	Subnet2IpRange string            `json:"subnet2_ip_range,omitempty"`
}

type SqlDBCollection map[string]SqlDB

type SQLTypeVersion struct {
	Type    SQLType `json:"type"`
	Version string  `json:"version"`
}

type SQLType string

const (
	PostgresSQLType SQLType = "postgres"
	MySqlSQLType    SQLType = "mysql"
	SQLServerType   SQLType = "sqlserver"
)

type SqlDBAuth struct {
	RootPassword string `json:"root_password"`
}
