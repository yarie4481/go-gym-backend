package controllers

import (
	"go-blog/internal/models"
	services "go-blog/internal/service"
	"net/http"
	"strings" // Added for error string checking
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// MembershipRequest is the Data Transfer Object (DTO) used for binding
// incoming JSON payloads for creating or updating a membership.
type MembershipRequest struct {
	// The UUID strings from JSON will be parsed here
	MemberID  uuid.UUID `json:"member_id" binding:"required"`
	PlanID    uuid.UUID `json:"plan_id" binding:"required"`
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
	Status    string    `json:"status" binding:"required"`
	AutoRenew bool      `json:"auto_renew"`
}

// MembershipController holds the business logic dependencies.
type MembershipController struct {
	service *services.MembershipService
}

// NewMembershipController creates a new instance of MembershipController.
func NewMembershipController(service *services.MembershipService) *MembershipController {
	return &MembershipController{service: service}
}

// CreateMembership handles the POST request to create a new membership record.
func (c *MembershipController) CreateMembership(ctx *gin.Context) {
	// Use the request DTO for binding.
	// The type is available directly since it's defined in the current package.
	var req MembershipRequest 
	
	// ShouldBindJSON with the `binding:"required"` tag will handle validation 
	// and correctly parse the UUIDs, returning an error on malformed UUIDs.
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format or missing field: " + err.Error()})
		return
	}
	
	// Validate for zero UUIDs (nil UUID) to ensure valid IDs were provided.
	if req.MemberID == uuid.Nil || req.PlanID == uuid.Nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Member ID and Plan ID cannot be zero/empty UUIDs."}) 
		return
	}

	// Map DTO to Model
	m := &models.Membership{
		MemberID:  req.MemberID,
		PlanID:    req.PlanID,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Status:    req.Status,
		AutoRenew: req.AutoRenew,
		// ID, CreatedAt, UpdatedAt are managed by the service/repository
	}

	// Call service layer with the populated model
	if err := c.service.CreateMembership(m); err != nil {
		
		errMsg := err.Error()

		// Intercept Foreign Key violation errors (SQLSTATE 23503) for better client feedback (400 Bad Request).
		// Note: A more robust solution involves returning domain-specific errors from the service layer.
		if strings.Contains(errMsg, "violates foreign key constraint \"fk_members_memberships\"") {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Member ID: The member specified in the request does not exist."})
			return
		}

		if strings.Contains(errMsg, "violates foreign key constraint \"fk_plans_memberships\"") {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Plan ID: The plan specified in the request does not exist."})
			return
		}

		// Default: return 500 Internal Server Error for unhandled database/service errors
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create membership: " + errMsg})
		return
	}
	
	ctx.JSON(http.StatusCreated, m)
}

// GetByMember retrieves all memberships associated with a given member ID.
func (c *MembershipController) GetByMember(ctx *gin.Context) {
	memberID := ctx.Param("memberID")
	// The service layer should handle validation and parsing of memberID string into uuid.UUID
	memberships, err := c.service.GetByMember(memberID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, memberships)
}
