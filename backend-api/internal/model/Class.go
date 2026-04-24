package model

import "gorm.io/gorm"

type Class struct {
	gorm.Model
	Name   string
	YearID uint
	Year   Year

	ClassCode         string `gorm:"uniqueIndex"`
	HomeroomTeacherID uint
	HomeroomTeacher   Employee `gorm:"foreignKey:HomeroomTeacherID"`
}
