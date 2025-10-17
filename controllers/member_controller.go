package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"go-blog/internal/models"
	services "go-blog/internal/service"
)

type MemberController struct {
	service services.MemberService
}

func NewMemberController(service services.MemberService) *MemberController {
	return &MemberController{service: service}
}

// POST /members
func (c *MemberController) CreateMember(ctx *gin.Context) {
	var member models.Member
	if err := ctx.ShouldBindJSON(&member); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.CreateMember(&member); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, member)
}

// GET /members
func (c *MemberController) GetAllMembers(ctx *gin.Context) {
	members, err := c.service.GetAllMembers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, members)
}

// GET /members/:id
func (c *MemberController) GetMemberByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	member, err := c.service.GetMemberByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "member not found"})
		return
	}
	ctx.JSON(http.StatusOK, member)
}

// PUT /members/:id
func (c *MemberController) UpdateMember(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	var member models.Member
	if err := ctx.ShouldBindJSON(&member); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	member.ID = id

	if err := c.service.UpdateMember(&member); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, member)
}

// DELETE /members/:id
func (c *MemberController) DeleteMember(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	if err := c.service.DeleteMember(id); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "member deleted"})
}
