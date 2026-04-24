package handler

import (
	"net/http"
	"os"
	"siakad-backend/internal/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

type JwtCustomClaims struct {
	UserID      uint     `json:"user_id"`
	EmployeeID  uint     `json:"employees_id"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

func (h *Handler) Login(c *echo.Context) error {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"message": "success"})
	}

	var user model.User
	if err := h.DB.Preload("Permissions").Where("username = ?", body.Username).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]any{"message": "Username not found"})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]any{"message": "Wrong Password"})
	}

	var perms []string
	for _, p := range user.Permissions {
		perms = append(perms, p.Name)
	}

	employeeID := uint(0)
	if user.EmployeeID != nil {
		employeeID = *user.EmployeeID
	}

	token, err := h.generateToken(user.ID, employeeID, perms)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{"message": "Token error"})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"token": token,
		"user": map[string]any{
			"username":    user.Username,
			"employee_id": employeeID,
			"permissions": perms,
		},
	})
}

func (h *Handler) generateToken(userID uint, employeeID uint, perms []string) (string, error) {
	claims := &JwtCustomClaims{
		UserID:      userID,
		EmployeeID:  employeeID,
		Permissions: perms,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
