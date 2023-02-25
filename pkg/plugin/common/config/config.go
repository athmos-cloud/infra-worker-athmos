package config

import (
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
	comPlugin "github.com/PaulBarrie/infra-worker/pkg/plugin/common"
)

type IPluginConfigStrategy interface {
	Build() errors.Error
}

type Config struct {
	PluginLocation comPlugin.Location `bson:"location"`
	Inputs         InputList          `bson:"inputs"`
	Outputs        OutputPayloadList  `bson:"outputs"`
	Dependencies   DependencyList     `bson:"dependencies"`
}
