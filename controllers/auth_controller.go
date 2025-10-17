package controllers

import (
	"go-blog/internal/models"
	services "go-blog/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthController struct {
	service *services.AuthService
}

func NewAuthController(service *services.AuthService) *AuthController {
	return &AuthController{service}
}
func (c *AuthController) Register(ctx *gin.Context) {
	var input struct {
		FirstName             string     `json:"first_name"`
		LastName              string     `json:"last_name"`
		Email                 string     `json:"email"`
		Password              string     `json:"password"`
		MembershipType        string     `json:"membership_type"`
		DateOfBirth           string     `json:"date_of_birth"`
		EmergencyContactName  string     `json:"emergency_contact_name"`
		EmergencyContactPhone string     `json:"emergency_contact_phone"`
		FitnessGoals          string     `json:"fitness_goals"`
		UserType              string     `json:"user_type"`
		Gender                string     `json:"gender"`

		// Membership info (optional, only for members)
		PlanID          string     `json:"plan_id"`
		MembershipStart *time.Time `json:"membership_start"`
		MembershipEnd   *time.Time `json:"membership_end"`
		AutoRenew       *bool      `json:"auto_renew"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse date_of_birth if provided
	var dob *time.Time
	if input.DateOfBirth != "" {
		t, err := time.Parse("2006-01-02", input.DateOfBirth)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid date_of_birth format, use YYYY-MM-DD"})
			return
		}
		dob = &t
	}

	user := &models.User{
		FirstName:             input.FirstName,
		LastName:              input.LastName,
		Email:                 input.Email,
		MembershipType:        input.MembershipType,
		DateOfBirth:           dob,
		EmergencyContactName:  input.EmergencyContactName,
		EmergencyContactPhone: input.EmergencyContactPhone,
		FitnessGoals:          input.FitnessGoals,
		UserType:              input.UserType,
		Gender:                input.Gender,
	}

	// Handle membership defaults
	var startDate, endDate time.Time
	if input.MembershipStart != nil {
		startDate = *input.MembershipStart
	} else {
		startDate = time.Now()
	}

	if input.MembershipEnd != nil {
		endDate = *input.MembershipEnd
	} else {
		endDate = startDate.AddDate(0, 1, 0) // default 1 month
	}

	autoRenew := true
	if input.AutoRenew != nil {
		autoRenew = *input.AutoRenew
	}

	// Handle plan_id for members
	var planID string
	if user.UserType == "member" || user.UserType == "Member" {
		if input.PlanID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "plan_id is required for members"})
			return
		}
		
		// Validate UUID format in controller
		if _, err := uuid.Parse(input.PlanID); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan_id: " + err.Error()})
			return
		}
		planID = input.PlanID
	} else {
		// For non-members, pass empty string
		planID = ""
	}

	// Call service
	if err := c.service.Register(user, input.Password, planID, startDate, endDate, autoRenew); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}
func (c *AuthController) Login(ctx *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, ok := c.service.Login(input.Email, input.Password)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	accessToken, refreshToken, err := c.service.GenerateTokens(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          user,
	})
}

// GetAllUsers returns all users
func (c *AuthController) GetAllUsers(ctx *gin.Context) {
	users, err := c.service.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"users": users,
		"count": len(users),
	})
}

// GetTrainers returns only users with user_type "Trainer"
func (c *AuthController) GetTrainers(ctx *gin.Context) {
	trainers, err := c.service.GetUsersByType("Trainer")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"trainers": trainers,
		"count":    len(trainers),
	})
}