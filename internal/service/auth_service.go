package services

import (
	"errors"
	"fmt"
	"go-blog/internal/config"
	"go-blog/internal/models"
	"go-blog/repositories"
	"go-blog/utils"
	"log"

	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type AuthService struct {
	repo repositories.UserRepository
	db   *gorm.DB
}

func NewAuthService(repo repositories.UserRepository, db *gorm.DB) *AuthService {
	return &AuthService{repo: repo, db: db}
}

func (s *AuthService) Register(
	user *models.User,
	password string,
	planID string,
	startDate, endDate time.Time,
	autoRenew bool,
) error {
	// 1. Hash password
	hashed, err := utils.HashPassword(password)
	if err != nil {
		log.Printf("ERROR: AuthService.Register failed to hash password for user %s: %v", user.Email, err)
		return fmt.Errorf("failed to process password")
	}
	user.PasswordHash = hashed

	return s.db.Transaction(func(tx *gorm.DB) error {
		// 2. Create User
		if err := tx.Create(user).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				log.Printf("WARN: Registration failed, user email %s already exists.", user.Email)
				return fmt.Errorf("user with this email already exists")
			}
			log.Printf("ERROR: AuthService.Register failed to create user %s: %v", user.Email, err)
			return fmt.Errorf("failed to register user due to database error")
		}

		// 3. Only create Member and Membership if user.UserType == "member"
		if user.UserType == "member" || user.UserType == "Member" {
			// Validate planID for members
			if planID == "" {
				return fmt.Errorf("plan_id is required for members")
			}

			// Convert planID string to uuid.UUID
			planUUID, err := uuid.Parse(planID)
			if err != nil {
				return fmt.Errorf("invalid plan_id: %v", err)
			}

			// Optional: Check if plan exists in database (uncomment if you have plan validation)
			/*
			var plan models.Plan
			if err := tx.First(&plan, "id = ?", planUUID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return fmt.Errorf("plan with id %s does not exist", planID)
				}
				return fmt.Errorf("failed to verify plan: %v", err)
			}
			*/

			// Create Member
			member := models.Member{
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Gender:    user.Gender,
				UserID:    user.UserID,
				Dob:       user.DateOfBirth,
				EmergencyContact: datatypes.JSON([]byte(
					fmt.Sprintf(`{"name":"%s","phone":"%s"}`, user.EmergencyContactName, user.EmergencyContactPhone),
				)),
			}

			if err := tx.Create(&member).Error; err != nil {
				log.Printf("ERROR: AuthService.Register failed to create member for user ID %s: %v", user.UserID, err)
				return fmt.Errorf("failed to finalize member profile; registration rolled back")
			}

			// Create Membership
			membership := models.Membership{
				MemberID:  member.ID, // <-- use Member.ID
				PlanID:    planUUID,
				StartDate: startDate,
				EndDate:   endDate,
				Status:    "active",
				AutoRenew: autoRenew,
			}

			if err := tx.Create(&membership).Error; err != nil {
				log.Printf("ERROR: AuthService.Register failed to create membership for member ID %s: %v", member.ID, err)
				return fmt.Errorf("failed to create membership; registration rolled back")
			}
		}

		return nil
	})
}

// In your services/auth_service.go

func (s *AuthService) GetAllUsers() ([]*models.User, error) {
	return s.repo.GetAllUsers()
}

// GetUsersByType retrieves users by their user type
func (s *AuthService) GetUsersByType(userType string) ([]*models.User, error) {
	return s.repo.GetUsersByUserType(userType)
}


func (s *AuthService) Login(email, password string) (*models.User, bool) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, false
	}
	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return nil, false
	}
	return user, true
}
func (s *AuthService) GenerateTokens(user *models.User) (accessToken, refreshToken string, err error) {
	accessToken, err = config.GenerateJWT(user.UserID.String(), user.UserType, time.Minute*15)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = config.GenerateJWT(user.UserID.String(), user.UserType, time.Hour*24*7)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
