package config

import (
	"github.com/PaulBarrie/infra-worker/pkg/infrastructure/repository"
	plugin2 "github.com/PaulBarrie/infra-worker/pkg/plugin/common"
)

type Dependency struct {
	Name      string            `bson:"name"`
	Version   string            `bson:"version"`
	Variables map[string]string `bson:"vars"`
	Path      string            `bson:"path"`
}

type DependencyList []Dependency

type DependencyCall struct {
	Name      string           `bson:"name"`
	Location  plugin2.Location `bson:"source"`
	Source    repository.IRepository
	Version   string  `bson:"version"`
	Variables []Input `bson:"vars"`
}

type DependencyCallList []DependencyCall
