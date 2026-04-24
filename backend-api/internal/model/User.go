package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email      string `gorm:"uniqueIndex;not null"`
	Username   string `gorm:"uniqueIndex;not null"`
	Password   string `gorm:"not null"`
	EmployeeID *uint

	Employee    Employee
	Permissions []Permission `gorm:"many2many:user_permissions;"`
}
