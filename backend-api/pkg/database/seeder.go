package database

import (
	"log"
	"siakad-backend/internal/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	perms := []model.Permission{
		{Name: "edit_all_grades", Desc: "Can edit any student grade"},
		{Name: "edit_own_grades", Desc: "Can only edit grades for courses they teach"},
		{Name: "view_homeroom", Desc: "Can view all grades for their assigned class"},
	}
	for _, p := range perms {
		db.FirstOrCreate(&p, model.Permission{Name: p.Name})
	}

	principalPos := model.Position{Name: "Kepala Sekolah"}
	teacherPos := model.Position{Name: "Guru"}
	db.FirstOrCreate(&principalPos, model.Position{Name: "Kepala Sekolah"})
	db.FirstOrCreate(&teacherPos, model.Position{Name: "Guru"})

	teacher := model.Employee{
		NIP:        "123456789",
		Name:       "Test Employee",
		Gender:     "L",
		PositionID: teacherPos.ID,
	}
	db.FirstOrCreate(&teacher, model.Employee{NIP: teacher.NIP})

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := model.User{
		Username:   "admin",
		Email:      "admin@siakad.com",
		Password:   string(hashedPassword),
		EmployeeID: &teacher.ID,
	}

	if err := db.FirstOrCreate(&admin, model.User{Username: "admin"}).Error; err == nil {
		db.Model(&admin).Association("Permissions").Replace(perms)
	}

	year := model.Year{Name: "2025/2026 Ganjil", IsActive: true}
	db.FirstOrCreate(&year, model.Year{Name: "2025/2026 Ganjil"})

	class := model.Class{
		ClassCode:         "10A-2025",
		Name:              "10-A",
		YearID:            year.ID,
		HomeroomTeacherID: teacher.ID,
	}
	db.FirstOrCreate(&class, model.Class{ClassCode: "10A-2025"})

	log.Println("Seeding completed successfully!")
}
