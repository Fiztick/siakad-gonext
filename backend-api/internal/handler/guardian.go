package handler

import (
	"net/http"
	"siakad-backend/internal/model"

	"github.com/labstack/echo/v5"
)

func (h *Handler) CreateGuardian(c *echo.Context) error {
	id := c.Param("studentId")

	var student model.Student
	if err := h.DB.First(&student, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Student not found"})
	}

	var guardian model.Guardian
	if err := c.Bind(&guardian); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request format"})
	}

	guardian.StudentID = student.ID

	if err := h.DB.Create(&guardian).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to created guardian: " + err.Error()})
	}

	return c.JSON(http.StatusCreated, guardian)
}

func (h *Handler) UpdateGuardian(c *echo.Context) error {
	id := c.Param("studentId")
	var guardian model.Guardian

	if err := h.DB.Where("student_id = ?", id).First(&guardian).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Guardian not found"})
	}

	if err := c.Bind(&guardian); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request format"})
	}

	if err := h.DB.Save(&guardian).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, guardian)
}
