package handler

import (
	"net/http"
	"siakad-backend/internal/model"

	"github.com/labstack/echo/v5"
)

func (h *Handler) GetPositions(c *echo.Context) error {
	var positions []model.Position

	if err := h.DB.Find(&positions).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"message": len(positions),
		"data":    positions,
	})
}

func (h *Handler) GetPosition(c *echo.Context) error {
	id := c.Param("positionId")
	var pos model.Position

	if err := h.DB.First(&pos, id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, pos)
}

func (h *Handler) CreatePosition(c *echo.Context) error {
	var pos model.Position

	if err := c.Bind(&pos); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request format"})
	}

	if err := h.DB.Create(&pos).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, pos)
}

func (h *Handler) UpdatePosition(c *echo.Context) error {
	id := c.Param("positionId")
	var pos model.Position

	if err := h.DB.First(&pos, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Position not found"})
	}

	if err := c.Bind(&pos); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request format"})
	}

	if err := h.DB.Save(&pos).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, pos)
}

func (h *Handler) DeletePosition(c *echo.Context) error {
	id := c.Param("positionId")
	var pos model.Position

	if err := h.DB.First(&pos, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Position not found"})
	}

	if err := h.DB.Delete(&pos).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Position deleted successfully"})
}
