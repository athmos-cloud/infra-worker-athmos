package common

import "github.com/PaulBarrie/infra-worker/pkg/plugin/common/config"

type Type string

const (
	Terraform Type = "terraform"
	Ansible   Type = "ansible"
)

type Location struct {
	Bucket string
	Folder string
}

type Plugin struct {
	Id          string            `bson:"id"`
	Name        string            `bson:"name"`
	Type        Type              `bson:"type"`
	Location    Location          `bson:"location"`
	Config      config.Config     `bson:"config"`
	Environment map[string]string `bson:"environment"`
}
