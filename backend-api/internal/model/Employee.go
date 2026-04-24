package model

import (
	"time"

	"gorm.io/gorm"
)

type Employee struct {
	gorm.Model
	NIP                  string `gorm:"uniqueIndex;not null"`
	Name                 string `gorm:"not null"`
	Gender               string `gorm:"type:char(1);"`
	PlaceOfBirth         string
	DateOfBirth          time.Time
	Phone                string
	Address              string
	EmployeeStatus       string
	LastEducation        string
	LastMajor            string
	LastPlaceOfEducation string
	Image                string

	PositionID uint
	Position   Position
}
