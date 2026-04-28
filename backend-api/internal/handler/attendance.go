package handler

import (
	"net/http"
	"siakad-backend/internal/model"

	"github.com/labstack/echo/v5"
)

func (h *Handler) GetAttendanceByClass(c *echo.Context) error {
	classId := c.QueryParam("classId")
	var attendances []model.Attendance

	if err := h.DB.Preload("Student").Where("class_id = ?", classId).Find(&attendances).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, attendances)
}

func (h *Handler) GetStudentAttendance(c *echo.Context) error {
	studentId := c.Param("studentId")
	classId := c.QueryParam("classId")
	var attendance model.Attendance

	if err := h.DB.Preload("Class").Where("student_id = ? AND class_id = ?", studentId, classId).First(&attendance).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Attendance record not found"})
	}

	return c.JSON(http.StatusOK, attendance)
}

func (h *Handler) UpsertAttendance(c *echo.Context) error {
	var attendance model.Attendance
	var input model.Attendance
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request format"})
	}

	var studentCount, classCount int64
	h.DB.Model(&model.Student{}).Where("id = ?", input.StudentID).Count(&studentCount)
	h.DB.Model(&model.Class{}).Where("id = ?", input.ClassID).Count(&classCount)

	if studentCount == 0 || classCount == 0 {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"message": "Student or Class not found"})
	}

	// update if record already exist
	err := h.DB.Where("student_id = ? AND class_id = ?", input.StudentID, input.ClassID).First(&attendance).Error

	if err == nil {
		attendance.TotalHadir = input.TotalHadir
		attendance.TotalAlpa = input.TotalAlpa
		attendance.TotalIzin = input.TotalIzin
		attendance.TotalSakit = input.TotalSakit

		if err := h.DB.Save(&attendance).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}

		return c.JSON(http.StatusOK, attendance)
	}

	// create if record doesnt exist
	if err := h.DB.Create(&input).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, input)
}
