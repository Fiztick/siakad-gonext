package model

import (
	"time"

	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	NIS               string
	Name              string
	Gender            string    `gorm:"type:char(1)"`
	PlaceOfBirth      string    `json:"place_of_birth"`
	DateOfBirth       time.Time `json:"date_of_birth"`
	Address           string
	Phone             string
	Religion          string
	Nationality       string
	ChildPosition     int    `json:"child_position"`
	SiblingCount      int    `json:"sibling_count"`
	SchoolOrigin      string `json:"school_origin"`
	Weight            int
	Height            int
	BloodType         string `gorm:"type:char(2)" json:"blood_type"`
	IQ                int
	Arts              string
	Sports            string
	Community         string
	QuranReading      string `json:"quran_reading"`
	IqraReadingVolume string `json:"iqra_reading_volume"`
	Status            string `gorm:"default:'Active'"`

	ClassID *uint `json:"class_id"`
	Class   Class `json:"class" gorm:"foreignKey:ClassID"`

	Guardian    Guardian
	Attendances []Attendance
	Grades      []CourseGrade
}
