package controllers

import (
	"go-blog/internal/models"
	services "go-blog/internal/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid" // Added UUID import for DTO
)

// PaymentRequest is the DTO used for binding incoming JSON payloads for payment creation.
type PaymentRequest struct {
	MemberID    uuid.UUID `json:"member_id" binding:"required"`
	AmountCents int       `json:"amount_cents" binding:"required"`
	Currency    string    `json:"currency" binding:"required"`
	Method      string    `json:"method" binding:"required"`
	Status      string    `json:"status" binding:"required"`
	Reference   string    `json:"reference"`
}

type PaymentController struct {
	service *services.PaymentService
}

func NewPaymentController(service *services.PaymentService) *PaymentController {
	return &PaymentController{service: service}
}

func (c *PaymentController) RecordPayment(ctx *gin.Context) {
	// Use DTO for binding the incoming JSON payload
	var req PaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format or missing field: " + err.Error()})
		return
	}

	// Validate for zero UUIDs just in case
	if req.MemberID == uuid.Nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Member ID cannot be an empty UUID."})
		return
	}

	// Map the DTO to the database model
	p := models.Payment{
		MemberID:    req.MemberID,
		AmountCents: req.AmountCents,
		Currency:    req.Currency,
		Method:      req.Method,
		Status:      req.Status,
		Reference:   req.Reference,
	}
	
	if err := c.service.RecordPayment(&p); err != nil {
		errMsg := err.Error()

		// Intercept Foreign Key violation errors for better client feedback (400 Bad Request).
		if strings.Contains(errMsg, "violates foreign key constraint \"fk_members_payments\"") {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Member ID: The member specified in the request does not exist in the system."})
			return
		}

		// Default: return 500 Internal Server Error for unhandled database/service errors
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record payment: " + errMsg})
		return
	}
	ctx.JSON(http.StatusCreated, p)
}

func (c *PaymentController) GetPayments(ctx *gin.Context) {
	memberID := ctx.Param("memberID")
	payments, err := c.service.GetPayments(memberID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, payments)
}
func (c *PaymentController) GetAllPayments(ctx *gin.Context) {
	payments, err := c.service.GetAllPayments()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, payments)
}
