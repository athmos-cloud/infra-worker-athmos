package project

import (
	"github.com/PaulBarrie/infra-worker/pkg/auth"
	"github.com/PaulBarrie/infra-worker/pkg/project/state"
	"github.com/jinzhu/gorm"
)

type Stack struct {
	gorm.Model
	Plugins []PluginInstance `json:"plugins"`
	State   state.State      `json:"state"`
}

type Project struct {
	gorm.Model
	Id               int       `json:"id" gorm:"primaryKey"`
	Name             string    `json:"name" gorm:"varchar(255);"`
	Owner            auth.User `json:"owner"`
	TerraformPlugins Stack     `json:"terraform_plugins"`
	CreatedAt        string    `json:"created_at"`
	UpdatedAt        string    `json:"updated_at"`
}
