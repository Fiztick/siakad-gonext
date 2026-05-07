package model

import "gorm.io/gorm"

type CourseGrade struct {
	gorm.Model
	AverageDailyGrade  int `json:"average_daily_grade"`
	FirstMidtermGrade  int `json:"first_midterm_grade"`
	SecondMidtermGrade int `json:"second_midterm_grade"`
	FirstFinalGrade    int `json:"first_final_grade"`
	SecondFinalGrade   int `json:"second_final_grade"`

	StudentID  uint `gorm:"not null;index" json:"student_id"`
	ClassID    uint `gorm:"not null;index" json:"class_id"`
	CourseID   uint `gorm:"not null;index" json:"course_id"`
	EmployeeID uint `gorm:"not null;index" json:"employee_id"`

	Student  Student
	Class    Class
	Course   Course
	Employee Employee
}
