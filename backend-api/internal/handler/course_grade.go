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

	// check if these exist
	checks := map[string]uint{
		"Student":  grade.StudentID,
		"Class":    grade.ClassID,
		"Course":   grade.CourseID,
		"Employee": grade.EmployeeID,
	}

	for table, id := range checks {
		var count int64
		h.DB.Table(table+"s").Where("id = ?", id).Count((&count))
		if count == 0 {
			return c.JSON(http.StatusUnprocessableEntity, map[string]string{"message": table + " not found"})
		}
	}

	// update if grade already exist
	err := h.DB.Where("student_id = ? AND class_id = ? AND course_id = ?",
		grade.StudentID, grade.ClassID, grade.CourseID).First(&existingGrade).Error
	if err == nil {
		grade.ID = existingGrade.ID
		if err := h.DB.Save(&grade).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}
		return c.JSON(http.StatusOK, grade)
	}

	// create if grade havent existed
	if err := h.DB.Create(&grade).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, grade)
}
