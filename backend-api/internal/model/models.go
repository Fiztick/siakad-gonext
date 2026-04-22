package model

import (
	"time"

	"gorm.io/gorm"
)

type Permission struct {
	gorm.Model
	Name string `gorm:"uniqueIndex;not null"`
	Desc string
}

type User struct {
	gorm.Model
	Email       string `gorm:"uniqueIndex;not null"`
	Username    string `gorm:"uniqueIndex;not null"`
	Password    string `gorm:"not null"`
	EmployeeID  *uint
	Employee    Employee
	Permissions []Permission `gorm:"many2many:user_permissions;"`
}

type Position struct {
	gorm.Model
	Name string `gorm:"not null"`
}

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

type Guardian struct {
	gorm.Model
	StudentID     uint
	Name          string `gorm:"not null"`
	Address       string
	Phone         string
	Relation      string
	Occupation    string
	Income        int
	LastEducation string
}

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

	Guardian    Guardian
	Attendances []Attendance
	Grades      []CourseGrade
}

type Year struct {
	gorm.Model
	Name     string `gorm:"not null"`
	IsActive bool   `gorm:"default:true"`
}

type Class struct {
	gorm.Model
	ClassCode         string `gorm:"uniqueIndex"`
	Name              string
	YearID            uint
	Year              Year
	HomeroomTeacherID uint
	HomeroomTeacher   Employee `gorm:"foreignKey:HomeroomTeacherID"`
}

type Course struct {
	gorm.Model
	Name string `gorm:"not null"`
}

type CourseGrade struct {
	gorm.Model
	StudentID  uint `gorm:"not null;index"`
	ClassID    uint `gorm:"not null;index"`
	CourseID   uint `gorm:"not null;index"`
	EmployeeID uint `gorm:"not null;index"`

	AverageDailyGrade  int
	FirstMidtermGrade  int
	SecondMidtermGrade int
	FirstFinalGrade    int
	SecondFinalGrade   int

	Student  Student
	Class    Class
	Course   Course
	Employee Employee
}

type Attendance struct {
	gorm.Model
	StudentID  uint `gorm:"not null"`
	ClassID    uint `gorm:"not null"`
	TotalHadir int
	TotalAlpa  int
	TotalIzin  int
	TotalSakit int
}
