package model

import "gorm.io/gorm"

type Attendance struct {
	gorm.Model
	TotalHadir int
	TotalAlpa  int
	TotalIzin  int
	TotalSakit int

	StudentID uint `gorm:"not null"`
	Student   Student

	ClassID uint `gorm:"not null"`
	Class   Class
}
