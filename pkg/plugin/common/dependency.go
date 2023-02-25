package common

import "github.com/jinzhu/gorm"

type Dependency struct {
	gorm.Model
	Name    string `bson:"name"`
	Version string `bson:"version"`
	Path    string `bson:"path"`
}
