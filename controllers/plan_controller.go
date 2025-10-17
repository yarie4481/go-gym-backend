package controllers

import (
	"go-blog/internal/models"
	services "go-blog/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PlanController struct {
	service *services.PlanService
}

func NewPlanController(service *services.PlanService) *PlanController {
	return &PlanController{service: service}
}

func (c *PlanController) CreatePlan(ctx *gin.Context) {
	var plan models.Plan
	if err := ctx.ShouldBindJSON(&plan); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.service.CreatePlan(&plan); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, plan)
}

func (c *PlanController) GetPlans(ctx *gin.Context) {
	plans, err := c.service.GetPlans()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, plans)
}
