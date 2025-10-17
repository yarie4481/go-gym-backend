package controllers

import (
	"go-blog/internal/models"
	services "go-blog/internal/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GymController struct {
	service *services.GymService
}

func NewGymController(service *services.GymService) *GymController {
	return &GymController{service: service}
}

// POST /gym
func (c *GymController) CreateGym(ctx *gin.Context) {
	var body struct {
		Name         string                 `json:"name" binding:"required"`
		Address      string                 `json:"address"`
		Phone        string                 `json:"phone"`
		Timezone     string                 `json:"timezone"`
		OpeningHours map[string]interface{} `json:"opening_hours"`
		Settings     map[string]interface{} `json:"settings"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	gym, err := c.service.CreateGym(body.Name, body.Address, body.Phone, body.Timezone, body.OpeningHours, body.Settings)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gym)
}

// GET /gym
func (c *GymController) ListGyms(ctx *gin.Context) {
	gyms, err := c.service.ListGyms()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gyms)
}

// GET /gym/:id
func (c *GymController) GetGym(ctx *gin.Context) {
	id := ctx.Param("id")
	gymUUID, _ := uuid.Parse(id)

	gym, err := c.service.GetGym(gymUUID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gym)
}

// PUT /gym/:id
func (c *GymController) UpdateGym(ctx *gin.Context) {
	id := ctx.Param("id")
	gymUUID, _ := uuid.Parse(id)

	var body struct {
		Name         string                 `json:"name"`
		Address      string                 `json:"address"`
		Phone        string                 `json:"phone"`
		Timezone     string                 `json:"timezone"`
		OpeningHours map[string]interface{} `json:"opening_hours"`
		Settings     map[string]interface{} `json:"settings"`
	}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	gym, err := c.service.GetGym(gymUUID)
	if err != nil {
		ctx.JSON(404, gin.H{"error": "Gym not found"})
		return
	}

	if body.Name != "" {
		gym.Name = body.Name
	}
	if body.Address != "" {
		gym.Address = body.Address
	}
	if body.Phone != "" {
		gym.Phone = body.Phone
	}
	if body.Timezone != "" {
		gym.Timezone = body.Timezone
	}
	if body.OpeningHours != nil {
		gym.OpeningHours = models.MapToJSON(body.OpeningHours)
	}
	if body.Settings != nil {
		gym.Settings = models.MapToJSON(body.Settings)
	}

	gym.UpdatedAt = time.Now()

	if err := c.service.UpdateGym(gym); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gym)
}
