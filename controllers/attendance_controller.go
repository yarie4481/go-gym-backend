package controllers

import (
	"net/http"

	services "go-blog/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AttendanceController struct {
	service *services.AttendanceService
}

func NewAttendanceController(service *services.AttendanceService) *AttendanceController {
	return &AttendanceController{service: service}
}

// ✅ POST /attendance/checkin
func (c *AttendanceController) CheckIn(ctx *gin.Context) {
	var payload struct {
		MemberID string `json:"member_id" binding:"required,uuid"`
		SessionID string `json:"session_id" binding:"required,uuid"`
		Method string `json:"method" binding:"required"` // qr, staff
	}

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	memberID, _ := uuid.Parse(payload.MemberID)
	sessionID, _ := uuid.Parse(payload.SessionID)

	if err := c.service.CheckIn(memberID, sessionID, payload.Method); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "✅ Member checked in successfully"})
}

// ✅ GET /attendance/member/:member_id
func (c *AttendanceController) GetMemberAttendance(ctx *gin.Context) {
	memberID, err := uuid.Parse(ctx.Param("member_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid member_id"})
		return
	}

	records, err := c.service.GetMemberAttendance(memberID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, records)
}

// ✅ GET /attendance/all
func (c *AttendanceController) GetAllAttendance(ctx *gin.Context) {
	records, err := c.service.GetAllAttendance()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, records)
}
