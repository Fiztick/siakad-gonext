package main

import (
	"flag"
	"net/http"
	"os"
	"siakad-backend/internal/handler"
	"siakad-backend/pkg/database"

	echoJWT "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
)

func main() {
	shouldSeed := flag.Bool("seed", false, "Set to true to seed the database")
	flag.Parse()

	db := database.InitDB()
	e := echo.New()
	h := &handler.Handler{DB: db}

	e.Static("/public", "public")

	if *shouldSeed {
		database.Seed(db)
		if flag.Lookup("seed").Value.String() == "true" {
			return
		}
	}

	e.GET("/health", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "SIAKAD API is Online",
			"db":     "Connected",
		})
	})
	e.POST("/login", h.Login)

	// Protected Routes
	protected := e.Group("/api")
	protected.Use(echoJWT.WithConfig(echoJWT.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}))

	protected.GET("/users", h.GetUsers)
	protected.POST("/user", h.CreateUser)
	protected.GET("/user/:userId", h.GetUser)
	protected.PUT("/user/:userId", h.UpdateUser)
	protected.DELETE("/user/:userId", h.DeleteUser)

	protected.GET("/permissions", h.GetPermissions)
	protected.PUT("/user/:userId/permissions", h.SyncUserPermissions)

	protected.GET("/students", h.GetStudents)
	protected.POST("/student", h.CreateStudent)
	protected.GET("/student/:studentId", h.GetStudent)
	protected.PUT("/student/:studentId", h.UpdateStudent)
	protected.DELETE("/student/:studentId", h.DeleteStudent)

	protected.POST("/student/:studentId/guardian", h.CreateGuardian)
	protected.PUT("/student/:studentId/guardian", h.UpdateGuardian)

	protected.GET("/positions", h.GetPositions)
	protected.POST("/position", h.CreatePosition)
	protected.GET("/position/:positionId", h.GetPosition)
	protected.PUT("/position/:positionId", h.UpdatePosition)
	protected.DELETE("/position/:positionId", h.DeletePosition)

	protected.GET("/employees", h.GetEmployees)
	protected.POST("/employee", h.CreateEmployee)
	protected.GET("/employees/:employeeId", h.GetEmployee)
	protected.PUT("/employee/:employeeId", h.UpdateEmployee)
	protected.DELETE("/employee/:employeeId", h.DeleteEmployee)

	protected.GET("/years", h.GetYears)
	protected.POST("/year", h.CreateYear)
	protected.GET("/year/:yearId", h.GetYear)
	protected.PUT("/year/:yearId", h.UpdateYear)
	protected.DELETE("/year/:yearId", h.DeleteYear)

	protected.GET("/classes", h.GetClasses)
	protected.POST("/class", h.CreateClass)
	protected.GET("/class/:classId", h.GetClass)
	protected.PUT("/class/:classId", h.UpdateClass)
	protected.DELETE("/class/:classId", h.DeleteClass)

	protected.GET("/class/:classId/students", h.GetClassStudents)
	protected.GET("/enrollment", h.EnrollStudent)
	protected.DELETE("/enrollment/:enrollmentId", h.DeleteEnrollment)

	protected.GET("/courses", h.GetCourses)
	protected.POST("/course", h.CreateCourse)
	protected.GET("/course/:courseId", h.GetCourse)
	protected.PUT("/course/:courseId", h.UpdateCourse)
	protected.DELETE("/course/:courseId", h.DeleteCourse)

	protected.POST("/grade", h.UpsertCourseGrade)
	protected.GET("/grades/class-view", h.GetGradeByClassCourse)
	protected.GET("/student/:studentId/report-card", h.GetStudentReportCard)

	protected.POST("/attendance", h.UpsertAttendance)
	protected.GET("/attendance/class", h.GetAttendanceByClass)
	protected.GET("/student/:studentId/attendance", h.GetStudentAttendance)

	if err := e.Start(":8080"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
