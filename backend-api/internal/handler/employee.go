package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"siakad-backend/internal/model"
	"time"

	"github.com/labstack/echo/v5"
)

func (h *Handler) GetEmployees(c *echo.Context) error {
	var emp []model.Employee
	db := h.DB.Preload("Position")

	// filter by name
	if name := c.QueryParam("name"); name != "" {
		db = db.Where("name ILIKE ?", "%"+name+"%")
	}

	if err := db.Find(&emp).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"count": len(emp),
		"data":  emp,
	})
}

func (h *Handler) GetEmployee(c *echo.Context) error {
	id := c.Param("employeeId")
	var emp model.Employee

	if err := h.DB.Preload("Position").First(&emp, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, emp)
}

func (h *Handler) CreateEmployee(c *echo.Context) error {
	var emp model.Employee
	var pos model.Position

	if err := c.Bind(&emp); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request format"})
	}

	// check if position exist
	if err := h.DB.First(&pos, &emp.PositionID).Error; err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"message": "Position doesn't exist"})
	}

	// processing file image
	fileHeader, err := c.FormFile("image")
	if err == nil {
		// check file size
		if fileHeader.Size > 200*1024 {
			return c.JSON(http.StatusUnprocessableEntity, map[string]string{
				"message": "Image size must be less than 200kb",
			})
		}

		uploadBase := "./public/uploads/employees"
		os.MkdirAll(uploadBase, os.ModePerm)
		newFileName := fmt.Sprintf("%d_%s", time.Now().Unix(), fileHeader.Filename)
		dstPath := filepath.Join(uploadBase, newFileName)

		src, err := fileHeader.Open()
		if err == nil {
			defer src.Close()

			// check extension
			buffer := make([]byte, 512)
			src.Read(buffer)
			contentType := http.DetectContentType(buffer)

			allowedTypes := map[string]bool{
				"image/jpeg": true,
				"image/jpg":  true,
			}
			if !allowedTypes[contentType] {
				return c.JSON(http.StatusUnprocessableEntity, map[string]string{
					"message": "File must be jpeg or jpg",
				})
			}

			src.Seek(0, io.SeekStart)

			dst, err := os.Create(dstPath)
			if err == nil {
				defer dst.Close()

				if _, err = io.Copy(dst, src); err == nil {
					emp.Image = "/uploads/employees/" + newFileName
				}
			}
		}
	}

	if err := h.DB.Create(&emp).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	h.DB.Preload("Position").First(&emp, emp.ID)

	return c.JSON(http.StatusCreated, emp)
}

func (h *Handler) UpdateEmployee(c *echo.Context) error {
	id := c.Param("employeeId")
	var emp model.Employee
	var pos model.Position

	if err := h.DB.First(&emp, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Employee not found"})
	}

	if err := c.Bind(&emp); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request format"})
	}

	// check if position exist
	if err := h.DB.First(&pos, &emp.PositionID).Error; err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"message": "Position doesn't exist"})
	}

	// processing file image
	fileHeader, err := c.FormFile("image")
	if err == nil {
		// check file size
		if fileHeader.Size > 200*1024 {
			return c.JSON(http.StatusUnprocessableEntity, map[string]string{
				"message": "Image size must be less than 200kb",
			})
		}

		oldFilePath := "./public" + emp.Image
		uploadBase := "./public/uploads/employees"
		os.MkdirAll(uploadBase, os.ModePerm)
		newFileName := fmt.Sprintf("%d_%s", time.Now().Unix(), fileHeader.Filename)
		dstPath := filepath.Join(uploadBase, newFileName)

		src, err := fileHeader.Open()
		if err == nil {
			defer src.Close()

			// check extension
			buffer := make([]byte, 512)
			src.Read(buffer)
			contentType := http.DetectContentType(buffer)

			allowedTypes := map[string]bool{
				"image/jpeg": true,
				"image/jpg":  true,
			}
			if !allowedTypes[contentType] {
				return c.JSON(http.StatusUnprocessableEntity, map[string]string{
					"message": "File must be jpeg or jpg",
				})
			}

			src.Seek(0, io.SeekStart)

			dst, err := os.Create(dstPath)
			if err == nil {
				defer dst.Close()

				if _, err = io.Copy(dst, src); err == nil {
					emp.Image = "/uploads/employees/" + newFileName

					os.Remove(oldFilePath)
				}
			}
		}
	}

	if err := h.DB.Save(&emp).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	h.DB.Preload("Position").First(&emp, emp.ID)

	return c.JSON(http.StatusOK, emp)
}

func (h *Handler) DeleteEmployee(c *echo.Context) error {
	id := c.Param("employeeId")
	var emp model.Employee

	if err := h.DB.First(&emp, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Employee not found"})
	}

	if err := h.DB.Delete(&emp).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Employee deleted successfully"})
}
