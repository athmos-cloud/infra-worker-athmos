package provider

import "github.com/PaulBarrie/infra-worker/pkg/resource/provider/auth"

type HelmApplication struct {
	Auth auth.HelmApplication `yaml:"auth"`
	VPC  string               `yaml:"vpc"`
}
