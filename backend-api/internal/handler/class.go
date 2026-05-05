package handler

import (
	"net/http"
	"siakad-backend/internal/model"

	"github.com/labstack/echo/v5"
)

func (h *Handler) GetClasses(c *echo.Context) error {
	var classes []model.Class

	if err := h.DB.Preload("Year").Preload("HomeroomTeacher").Find(&classes).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"count": len(classes),
		"data":  classes,
	})
}

func (h *Handler) GetClass(c *echo.Context) error {
	id := c.Param("classId")
	var class model.Class

	if err := h.DB.Preload("Year").Preload("HomeroomTeacher").First(&class, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Class not found"})
	}

	return c.JSON(http.StatusOK, class)
}

func (h *Handler) CreateClass(c *echo.Context) error {
	var class model.Class
	var year model.Year
	var ht model.Employee

	if err := c.Bind(&class); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request format"})
	}

	// check if year exist
	if err := h.DB.First(&year, class.YearID).Error; err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"message": "Year not found"})
	}

	//check if homeroom teacher exist
	if err := h.DB.First(&ht, class.HomeroomTeacherID).Error; err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"message": "Teacher not found in employee table"})
	}

	if err := h.DB.Create(&class).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	h.DB.Preload("Year").Preload("HomeroomTeacher").First(&class, class.ID)

	return c.JSON(http.StatusCreated, class)
}

func (h *Handler) UpdateClass(c *echo.Context) error {
	id := c.Param("classId")
	var class model.Class
	var year model.Year
	var ht model.Employee

	if err := h.DB.First(&class, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Class not found"})
	}

	if err := c.Bind(&class); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request format"})
	}

	// check if year exist
	if err := h.DB.First(&year, &class.YearID).Error; err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"message": "Year not found"})
	}

	//check if homeroom teacher exist
	if err := h.DB.First(&ht, &class.HomeroomTeacherID).Error; err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"message": "Teacher not found in employee table"})
	}

	if err := h.DB.Save(&class).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	h.DB.Preload("Year").Preload("HomeroomTeacher").Preload("HomeroomTeacher.Position").First(&class, class.ID)

	return c.JSON(http.StatusOK, class)
}

func (h *Handler) DeleteClass(c *echo.Context) error {
	id := c.Param("classId")
	var class model.Class

	if err := h.DB.First(&class, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Class not found"})
	}

	if err := h.DB.Delete(&class).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Class deleted successfully"})
}
