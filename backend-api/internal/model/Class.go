package model

import "gorm.io/gorm"

type Class struct {
	gorm.Model
	ClassCode string `gorm:"uniqueIndex" json:"class_code"`
	Name      string
	YearID    uint `json:"year_id"`
	Year      Year `json:"year"`

	HomeroomTeacherID uint     `json:"homeroom_teacher_id"`
	HomeroomTeacher   Employee `gorm:"foreignKey:HomeroomTeacherID" json:"homeroom_teacher"`
}
