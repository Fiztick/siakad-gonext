package model

import "gorm.io/gorm"

type Attendance struct {
	gorm.Model
	TotalHadir int `json:"total_hadir"`
	TotalAlpa  int `json:"total_alpa"`
	TotalIzin  int `json:"total_izin"`
	TotalSakit int `json:"total_sakit"`

	StudentID uint `gorm:"not null" json:"student_id"`
	Student   Student

	ClassID uint `gorm:"not null" json:"class_id"`
	Class   Class
}
