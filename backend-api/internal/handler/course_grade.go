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
	var input model.CourseGrade
	var existingGrade model.CourseGrade

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request format"})
	}

	// check student
	var studentCount int64
	h.DB.Model(&model.Student{}).Where("id = ?", input.StudentID).Count(&studentCount)
	if studentCount == 0 {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"message": "Student not found"})
	}

	// check class
	var classCount int64
	h.DB.Model(&model.Class{}).Where("id = ?", input.ClassID).Count(&classCount)
	if classCount == 0 {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"message": "Class not found"})
	}

	// check course
	var courseCount int64
	h.DB.Model(&model.Course{}).Where("id = ?", input.CourseID).Count(&courseCount)
	if courseCount == 0 {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"message": "Course not found"})
	}

	// check employee
	var employeeCount int64
	h.DB.Model(&model.Employee{}).Where("id = ?", input.EmployeeID).Count(&employeeCount)
	if employeeCount == 0 {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"message": "Employee not found"})
	}

	// update if grade already exist
	err := h.DB.Where("student_id = ? AND class_id = ? AND course_id = ?",
		input.StudentID, input.ClassID, input.CourseID).First(&existingGrade).Error
	if err == nil {
		existingGrade.AverageDailyGrade = input.AverageDailyGrade
		existingGrade.FirstMidtermGrade = input.FirstMidtermGrade
		existingGrade.SecondMidtermGrade = input.SecondMidtermGrade
		existingGrade.FirstFinalGrade = input.FirstFinalGrade
		existingGrade.SecondFinalGrade = input.SecondFinalGrade
		if err := h.DB.Save(&existingGrade).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}

		h.DB.Preload("Student").Preload("Class").Preload("Course").First(&existingGrade, existingGrade.ID)

		return c.JSON(http.StatusOK, existingGrade)
	}

	// create if input dont exist
	if err := h.DB.Create(&input).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	h.DB.Preload("Student").Preload("Class").Preload("Course").First(&input, input.ID)

	return c.JSON(http.StatusCreated, input)
}
