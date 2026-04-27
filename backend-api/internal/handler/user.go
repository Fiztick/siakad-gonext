package handler

import (
	"net/http"
	"siakad-backend/internal/model"

	"github.com/labstack/echo/v5"
)

func (h *Handler) GetUsers(c *echo.Context) error {
	var users []model.User

	if err := h.DB.Preload("Employee").Preload("Permissions").Find(&users).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"count": len(users),
		"data":  users,
	})
}

func (h *Handler) GetUser(c *echo.Context) error {
	id := c.Param("userId")
	var user model.User

	if err := h.DB.First(&user, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
	}

	return c.JSON(http.StatusOK, user)
}

func (h *Handler) CreateUser(c *echo.Context) error {
	var user model.User
	var employee model.Employee

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request format"})
	}

	// check if employee exist
	if err := h.DB.First(&employee, user.EmployeeID).Error; err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"messsage": "Employee not found"})
	}

	if err := h.DB.Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, user)
}

func (h *Handler) UpdateUser(c *echo.Context) error {
	id := c.Param("userId")
	var user model.User
	var employee model.Employee

	if err := h.DB.Find(&user, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
	}

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request format"})
	}

	// check if employee exist if changing
	if user.EmployeeID != nil {
		var count int64
		h.DB.Model(&employee).Where("id = ?", *user.EmployeeID).Count(&count)
		if count == 0 {
			return c.JSON(http.StatusUnprocessableEntity, map[string]string{"message": "Employee not found"})
		}
	}

	if err := h.DB.Save(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	h.DB.Preload("Employee").Preload("Permissions").First(&user, user.ID)

	return c.JSON(http.StatusOK, user)
}

func (h *Handler) DeleteUser(c *echo.Context) error {
	id := c.Param("userId")
	var user model.User

	if err := h.DB.First(&user, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
	}

	if err := h.DB.Delete(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "User deleted successfully"})
}
