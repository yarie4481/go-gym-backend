package controllers

import (
	services "go-blog/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ClassController struct {
	service *services.ClassService
}

func NewClassController(service *services.ClassService) *ClassController {
	return &ClassController{service: service}
}

// POST /class
func (c *ClassController) CreateClass(ctx *gin.Context) {
	var body struct {
		GymID           string `json:"gym_id" binding:"required"`
		TrainerID       string `json:"trainer_id" binding:"required"`
		Title           string `json:"title" binding:"required"`
		Description     string `json:"description"`
		Capacity        int    `json:"capacity" binding:"required"`
		DurationMinutes int    `json:"duration_minutes" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	gymUUID, _ := uuid.Parse(body.GymID)
	trainerUUID, _ := uuid.Parse(body.TrainerID)

	class, err := c.service.CreateClass(gymUUID, trainerUUID, body.Title, body.Description, body.Capacity, body.DurationMinutes)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, class)
}

// GET /class
func (c *ClassController) ListClasses(ctx *gin.Context) {
	classes, err := c.service.ListClasses()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, classes)
}

// GET /class/:id
func (c *ClassController) GetClass(ctx *gin.Context) {
	id := ctx.Param("id")
	classUUID, _ := uuid.Parse(id)

	class, err := c.service.GetClass(classUUID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, class)
}
