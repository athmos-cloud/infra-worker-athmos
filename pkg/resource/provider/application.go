package provider

import "github.com/PaulBarrie/infra-worker/pkg/resource/provider/auth"

type HelmApplication struct {
	Secret auth.SecretHelmApplication `yaml:"secret"`
	VPC    string                     `yaml:"vpc"`
}
