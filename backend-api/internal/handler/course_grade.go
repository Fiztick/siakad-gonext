package handler

import (
	"net/http"
	"siakad-backend/internal/model"

	"github.com/labstack/echo/v5"
)

func (h *Handler) GetGradeByClassCourse(c *echo.Context) error {
	classId := c.QueryParam("classId")
	courseId := c.QueryParam("courseId")
	var grades []model.CourseGrade

	err := h.DB.Preload("Student").
		Where("class_id = ? AND course_id = ?", classId, courseId).
		Find(&grades).Error

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, grades)
}

func (h *Handler) GetStudentReportCard(c *echo.Context) error {
	studentId := c.Param("studentId")
	classId := c.QueryParam("classId")
	var report []model.CourseGrade
	err := h.DB.Preload("Course").Preload("Employee").
		Where("student_id = ? AND class_id = ?", studentId, classId).
		Find(&report).Error

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, report)
}

func (h *Handler) UpsertCourseGrade(c *echo.Context) error {
	var grade model.CourseGrade
	var existingGrade model.CourseGrade

	if err := c.Bind(&grade); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request format"})
	}

	// check student
	var student model.Student
	if err := h.DB.First(&student, grade.StudentID).Error; err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"message": "Student not found"})
	}

	// check class
	var class model.Class
	if err := h.DB.First(&class, grade.ClassID).Error; err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"message": "Class not found"})
	}

	// check course
	var course model.Course
	if err := h.DB.First(&course, grade.CourseID).Error; err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"message": "Course not found"})
	}

	// check employee
	var employee model.Employee
	if err := h.DB.First(&employee, grade.EmployeeID).Error; err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"message": "Employee not found"})
	}

	// update if grade already exist
	err := h.DB.Where("student_id = ? AND class_id = ? AND course_id = ?",
		grade.StudentID, grade.ClassID, grade.CourseID).First(&existingGrade).Error
	if err == nil {
		existingGrade.AverageDailyGrade = grade.AverageDailyGrade
		existingGrade.FirstMidtermGrade = grade.FirstMidtermGrade
		existingGrade.SecondMidtermGrade = grade.SecondMidtermGrade
		existingGrade.FirstFinalGrade = grade.FirstFinalGrade
		existingGrade.SecondFinalGrade = grade.SecondFinalGrade
		if err := h.DB.Save(&existingGrade).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}

		h.DB.Preload("Student").Preload("Class").Preload("Course").First(&existingGrade, existingGrade.ID)

		return c.JSON(http.StatusOK, existingGrade)
	}

	// create if grade dont exist
	if err := h.DB.Create(&grade).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	h.DB.Preload("Student").Preload("Class").Preload("Course").First(&grade, grade.ID)

	return c.JSON(http.StatusCreated, grade)
}
