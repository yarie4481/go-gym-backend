package repositories

import (
	"go-blog/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetAllUsers() ([]*models.User, error)
	GetUsersByUserType(userType string) ([]*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := r.db.Where("email = ?", email).First(&user)
	return &user, result.Error
}

// GetAllUsers retrieves all users from the database
func (r *userRepository) GetAllUsers() ([]*models.User, error) {
	var users []*models.User
	result := r.db.Find(&users)
	return users, result.Error
}

// GetUsersByUserType retrieves users filtered by user_type
func (r *userRepository) GetUsersByUserType(userType string) ([]*models.User, error) {
	var users []*models.User
	result := r.db.Where("user_type = ?", userType).Find(&users)
	return users, result.Error
}

// Helper to expose DB if needed (optional)
func (r *userRepository) GetDB() *gorm.DB {
	return r.db
}