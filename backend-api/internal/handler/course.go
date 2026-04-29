package handler

import (
	"net/http"
	"siakad-backend/internal/model"

	"github.com/labstack/echo/v5"
)

func (h *Handler) GetCourses(c *echo.Context) error {
	var courses []model.Course

	if err := h.DB.Find(&courses).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"count": len(courses),
		"data":  courses,
	})
}

func (h *Handler) GetCourse(c *echo.Context) error {
	id := c.Param("courseId")
	var course model.Course

	if err := h.DB.Find(&course, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Course not found"})
	}

	return c.JSON(http.StatusOK, course)
}

func (h *Handler) CreateCourse(c *echo.Context) error {
	var course model.Course

	if err := c.Bind(&course); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request format"})
	}

	if err := h.DB.Create(&course).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, course)
}

func (h *Handler) UpdateCourse(c *echo.Context) error {
	id := c.Param("courseId")
	var course model.Course

	if err := h.DB.Find(&course, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Course not found"})
	}

	if err := c.Bind(&course); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request format"})
	}

	if err := h.DB.Save(&course).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, course)
}

func (h *Handler) DeleteCourse(c *echo.Context) error {
	id := c.Param("courseId")
	var course model.Course

	if err := h.DB.Find(&course, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Course not found"})
	}

	if err := h.DB.Delete(&course).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Course deleted successfully"})
}
