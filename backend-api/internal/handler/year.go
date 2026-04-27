package handler

import (
	"net/http"
	"siakad-backend/internal/model"

	"github.com/labstack/echo/v5"
)

func (h *Handler) GetYears(c *echo.Context) error {
	var years []model.Year

	if err := h.DB.Find(&years).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"message": len(years),
		"data":    years,
	})
}

func (h *Handler) GetYear(c *echo.Context) error {
	id := c.Param("yearId")
	var year model.Year

	if err := h.DB.First(&year, id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, year)
}

func (h *Handler) CreateYear(c *echo.Context) error {
	var year model.Year

	if err := c.Bind(&year); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request format"})
	}

	if err := h.DB.Create(&year).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, year)
}
func (h *Handler) UpdateYear(c *echo.Context) error {
	id := c.Param("yearId")
	var year model.Year

	if err := h.DB.First(&year, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Year not found"})
	}

	if err := c.Bind(&year); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request format"})
	}

	if err := h.DB.Save(&year).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, year)
}

func (h *Handler) DeleteYear(c *echo.Context) error {
	id := c.Param("yearId")
	var year model.Year

	if err := h.DB.First(&year, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Year not found"})
	}

	if err := h.DB.Delete(&year).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Year deleted successfully"})
}
