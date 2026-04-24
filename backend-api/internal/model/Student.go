package model

import (
	"time"

	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	NIS               string
	Name              string
	Gender            string `gorm:"type:char(1)"`
	PlaceOfBirth      string
	DateOfBirth       time.Time
	Address           string
	Phone             string
	Religion          string
	Nationality       string
	ChildPosition     int
	SiblingCount      int
	SchoolOrigin      string
	Weight            int
	Height            int
	BloodType         string `gorm:"type:char(2)"`
	IQ                int
	Arts              string
	Sports            string
	Community         string
	QuranReading      string
	IqraReadingVolume string
	Status            string `gorm:"default:'Active'"`

	ClassID uint  `json:"class_id"`
	Class   Class `json:"class" gorm:"foreignKey:ClassID"`

	Guardian    Guardian
	Attendances []Attendance
	Grades      []CourseGrade
}
