// controllers/classsession_controller.go
package controllers

import (
	services "go-blog/internal/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ClassSessionController struct {
	service *services.ClassSessionService
}

func NewClassSessionController(service *services.ClassSessionService) *ClassSessionController {
	return &ClassSessionController{service: service}
}

// POST /classsession
func (c *ClassSessionController) CreateSession(ctx *gin.Context) {
	var body struct {
		ClassID  string `json:"class_id" binding:"required"` // Add class_id to request
		StartsAt string `json:"starts_at" binding:"required"`
		EndsAt   string `json:"ends_at" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Parse class ID from request
	classUUID, err := uuid.Parse(body.ClassID)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid class ID"})
		return
	}

	startTime, err := time.Parse(time.RFC3339, body.StartsAt)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid start time format"})
		return
	}

	endTime, err := time.Parse(time.RFC3339, body.EndsAt)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid end time format"})
		return
	}

	session, err := c.service.CreateSession(classUUID, startTime, endTime)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, session)
}

// GET /classsession
func (c *ClassSessionController) ListSessions(ctx *gin.Context) {
	sessions, err := c.service.ListSessions()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, sessions)
}

// GET /classsession/:id
func (c *ClassSessionController) GetSession(ctx *gin.Context) {
	id := ctx.Param("id")
	sessionUUID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid session ID"})
		return
	}
	
	session, err := c.service.GetSession(sessionUUID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, session)
}