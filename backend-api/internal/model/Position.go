package model

import "gorm.io/gorm"

type Position struct {
	gorm.Model
	Name string `gorm:"not null"`
}
