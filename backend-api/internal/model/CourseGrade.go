package model

import "gorm.io/gorm"

type CourseGrade struct {
	gorm.Model
	AverageDailyGrade  int
	FirstMidtermGrade  int
	SecondMidtermGrade int
	FirstFinalGrade    int
	SecondFinalGrade   int

	StudentID  uint `gorm:"not null;index"`
	ClassID    uint `gorm:"not null;index"`
	CourseID   uint `gorm:"not null;index"`
	EmployeeID uint `gorm:"not null;index"`

	Student  Student
	Class    Class
	Course   Course
	Employee Employee
}
