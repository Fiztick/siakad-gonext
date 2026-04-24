package database

import (
	"fmt"
	"log"
	"siakad-backend/internal/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	// Permission Seeder
	perms := []model.Permission{
		{Name: "edit_all_grades", Desc: "Can edit any student grade"},
		{Name: "edit_own_grades", Desc: "Can only edit grades for courses they teach"},
		{Name: "view_homeroom", Desc: "Can view all grades for their assigned class"},
	}
	for i := range perms {
		db.FirstOrCreate(&perms[i], model.Permission{Name: perms[i].Name})
	}

	// Position Seeder
	principalPos := model.Position{Name: "Kepala Sekolah"}
	teacherPos := model.Position{Name: "Guru"}
	db.FirstOrCreate(&principalPos, model.Position{Name: "Kepala Sekolah"})
	db.FirstOrCreate(&teacherPos, model.Position{Name: "Guru"})

	teacher := model.Employee{
		NIP:        "00000000",
		Name:       "Admin",
		Gender:     "L",
		PositionID: teacherPos.ID,
	}
	db.FirstOrCreate(&teacher, model.Employee{NIP: teacher.NIP})

	// Admin User Seeder
	var admin model.User
	err := db.Where("username = ?", "admin").First(&admin).Error

	if err != nil {
		fmt.Println("Creating admin user...")
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)

		admin = model.User{
			Username:   "admin",
			Email:      "admin@siakad.com",
			Password:   string(hashedPassword),
			EmployeeID: &teacher.ID,
		}
		db.Create(&admin)
	}
	fmt.Println("Syncing admin permissions...")
	db.Model(&admin).Association("Permissions").Replace(perms)

	// Year Seeder
	year := model.Year{Name: "2025/2026 Ganjil", IsActive: true}
	db.FirstOrCreate(&year, model.Year{Name: "2025/2026 Ganjil"})

	// Class Seeder
	class := model.Class{
		ClassCode:         "10A-2025",
		Name:              "10-A",
		YearID:            year.ID,
		HomeroomTeacherID: teacher.ID,
	}
	db.FirstOrCreate(&class, model.Class{ClassCode: "10A-2025"})

	log.Println("Seeding completed successfully!")
}
