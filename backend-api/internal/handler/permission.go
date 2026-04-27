package handler

import (
	"net/http"
	"siakad-backend/internal/model"

	"github.com/labstack/echo/v5"
)

func (h *Handler) GetPermissions(c *echo.Context) error {
	var perms []model.Permission

	if err := h.DB.Find(&perms).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"count": len(perms),
		"data":  perms,
	})
}

func (h *Handler) SyncUserPermissions(c *echo.Context) error {
	userId := c.Param("userId")
	var user model.User
	var perms []model.Permission

	var input struct {
		PermissionIDs []uint `json:"permission_ids"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request format"})
	}

	// check if user exist
	if err := h.DB.First(&user, userId).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
	}

	// get all perms
	if err := h.DB.Find(&perms, input.PermissionIDs).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error fetching permissions"})
	}

	// sync perms to user
	if err := h.DB.Model(&user).Association("Permissions").Replace(perms); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to sync permissions"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Permissions synced successfully"})
}
