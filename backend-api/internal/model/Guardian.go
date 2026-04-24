package model

import "gorm.io/gorm"

type Guardian struct {
	gorm.Model
	Name          string `gorm:"not null"`
	Address       string
	Phone         string
	Relation      string
	Occupation    string
	Income        int
	LastEducation string

	StudentID uint
}
