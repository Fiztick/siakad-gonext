package handler

import (
	"net/http"
	"siakad-backend/internal/model"

	"github.com/labstack/echo/v5"
)

func (h *Handler) GetClassStudents(c *echo.Context) error {
	classId := c.Param("classId")
	var class model.Class
	var members []model.ClassStudent

	if err := h.DB.First(&class, classId).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Class not found"})
	}

	if err := h.DB.Preload("Student").Preload("Student.Guardian").Preload("Class.Year").Preload("Class.HomeroomTeacher.Position").Where("class_id = ?", classId).Find(&members).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	h.DB.Preload("Year").Preload("HomeroomTeacher.Position").First(&class, class.ID)

	return c.JSON(http.StatusOK, map[string]any{
		"class":         class,
		"student_total": len(members),
		"members":       members,
	})
}

func (h *Handler) EnrollStudent(c *echo.Context) error {
	var enrollment model.ClassStudent
	var student model.Student
	var class model.Class

	if err := c.Bind(&enrollment); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request format"})
	}

	// check if student exist
	if err := h.DB.First(&student, enrollment.StudentID).Error; err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"message": "Student not found"})
	}

	// check if class exist
	if err := h.DB.First(&class, enrollment.ClassID).Error; err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"message": "Class not found"})
	}

	// prevent duplicate
	var count int64
	h.DB.Model(&model.ClassStudent{}).Where("class_id = ? AND student_id = ?", enrollment.ClassID, enrollment.StudentID).Count(&count)
	if count > 0 {
		return c.JSON(http.StatusConflict, map[string]string{"message": "Student already enrolled in this class"})
	}

	if err := h.DB.Create(&enrollment).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	h.DB.Preload("Class").Preload("Student").Preload("Student.Guardian").First(&enrollment, enrollment.ID)

	return c.JSON(http.StatusCreated, enrollment)
}

func (h *Handler) DeleteEnrollment(c *echo.Context) error {
	id := c.Param("enrollmentId")
	var enrollment model.ClassStudent

	if err := h.DB.First(&enrollment, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Enrollment not found"})
	}

	if err := h.DB.Delete(&enrollment).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Enrollment deleted successfully"})
}
