package model

import (
	"github.com/jinzhu/gorm"
)

type UserItem struct {
	gorm.Model
	UserName      string `gorm:"primary_key"`
	Password      string `gorm:"not null"`
	WarehouseAddr string `gorm:"not null"`
}
