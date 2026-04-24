package handler

import (
	"net/http"
	"siakad-backend/internal/model"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func (h *Handler) GetStudents(c *echo.Context) error {
	var students []model.Student
	db := h.DB.Preload("Class")

	// filter by name
	if name := c.QueryParam("name"); name != "" {
		db = db.Where("name ILIKE ?", "%"+name+"%")
	}

	//filter by query
	if classID := c.QueryParam("class_Id"); classID != "" {
		db = db.Where("class_id = ?", classID)
	}

	// execute
	if err := db.Find(&students).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"message": len(students),
		"data":    students,
	})
}

func (h *Handler) GetStudent(c *echo.Context) error {
	id := c.Param("studentId")
	var student model.Student

	if err := h.DB.Preload("Guardian").Preload("Class").First(&student, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"message": "Student retrieved successfully",
		"data":    student,
	})
}

func (h *Handler) CreateStudent(c *echo.Context) error {
	var student model.Student

	if err := c.Bind(&student); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request format"})
	}

	if err := h.DB.Omit("Guardian", "Attendances", "Grades").Create(&student).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]any{
		"message": "Student created successfully",
		"data":    student,
	})
}

func (h *Handler) UpdateStudent(c *echo.Context) error {
	id := c.Param("studentId")
	var student model.Student

	if err := h.DB.First(&student, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Student not found"})
	}

	if err := c.Bind(&student); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request format"})
	}

	if err := h.DB.Save(&student).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"message": "Student updated successfully",
		"data":    student,
	})
}

func (h *Handler) DeleteStudent(c *echo.Context) error {
	id := c.Param("studentId")
	var student model.Student

	if err := h.DB.First(&student, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Student not found"})
	}

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&student).Error; err != nil {
			return err
		}

		if err := tx.Where("student_id = ?", student.ID).Delete(&model.Guardian{}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Student and guardian deleted successfully",
	})
}
