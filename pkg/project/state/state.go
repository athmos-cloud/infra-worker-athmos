package state

import (
	"github.com/PaulBarrie/infra-worker/pkg/auth"
	"github.com/jinzhu/gorm"
)

type State struct {
	gorm.Model
	Id             int        `json:"id" gorm:"primaryKey"`
	Owner          auth.Owner `json:"owner"`
	Versions       []Version
	CurrentVersion Version
	BackendType    BackendType
	UpdatedAt      string
	CreatedAt      string
}

type Version struct {
	gorm.Model
	Id        int `json:"id"`
	CreatedAt string
}
