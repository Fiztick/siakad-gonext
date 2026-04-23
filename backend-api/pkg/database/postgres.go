package database

import (
	"fmt"
	"log"
	"os"
	"siakad-backend/internal/model"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port)

	var db *gorm.DB
	var err error
	attempt := 10

	for i := 1; i <= attempt; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			sqlDB, _ := db.DB()
			if err = sqlDB.Ping(); err == nil {
				break
			}
		}
		log.Printf("Postgres not ready yet (attempt %d/%d), waiting 5s... ", i, attempt)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Fatalf("Could not connect to database after %d retries: %v", attempt, err)
	}

	fmt.Println("Running Database Migrations...")
	err = db.AutoMigrate(
		&model.Permission{},
		&model.Position{},
		&model.Year{},
		&model.Employee{},
		&model.User{},
		&model.Student{},
		&model.Guardian{},
		&model.Class{},
		&model.Course{},
		&model.CourseGrade{},
		&model.Attendance{},
	)

	if err != nil {
		log.Fatal("Migration Failed: ", err)
	}

	fmt.Println("Database Connected & Migrated Successfully")
	return db
}
