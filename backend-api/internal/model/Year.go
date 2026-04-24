package model

import "gorm.io/gorm"

type Year struct {
	gorm.Model
	Name     string `gorm:"not null"`
	IsActive bool   `gorm:"default:true"`
}
