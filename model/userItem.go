package model

import (
	"github.com/jinzhu/gorm"
)

type UserItem struct {
	gorm.Model
	UserName      string	`gorm:"not null"`
	Password      string	`gorm:"not null"`
	WarehouseAddr string	`gorm:"not null"`
}
