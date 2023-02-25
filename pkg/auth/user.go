package auth

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Id       int    `json:"id" gorm:"primaryKey"`
	Username string `json:"username"`
}
