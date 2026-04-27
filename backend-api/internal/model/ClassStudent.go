package model

import "gorm.io/gorm"

type ClassStudent struct {
	gorm.Model
	ClassID uint  `json:"class_id"`
	Class   Class `json:"class"`

	StudentID uint    `json:"student_id"`
	Student   Student `json:"student"`

	Status string `json:"status" gorm:"default:active"`
}
