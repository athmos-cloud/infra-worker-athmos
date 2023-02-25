package terraform_config

import (
	common "github.com/PaulBarrie/infra-worker/pkg/plugin/common/config"
	"github.com/hashicorp/hcl/v2/hclsimple"
)

type TerraformConfig struct {
	Config common.Config `hcl:"terraform"`
}

type Config struct {
	RequiredProviders map[string]RequiredProviders `hcl:"required_providers,optional"`
}

type RequiredProviders struct {
	Version string `hcl:"version"`
	Source  string `hcl:"source"`
}

func New(dependencyList common.DependencyList) *Config {
	requiredProviders := make(map[string]RequiredProviders)
	for _, dependency := range dependencyList {
		requiredProviders[dependency.Name] = RequiredProviders{
			Version: dependency.Version,
			Source:  dependency.Path,
		}
	}
	return &Config{requiredProviders}
}

func (c *Config) ToString() string {
	hclsimple.EncodeToString(c, hcl.NewBody)
}
