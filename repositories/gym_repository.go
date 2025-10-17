package repositories

import (
	"go-blog/internal/models"

	"gorm.io/gorm"
)

type GymRepository struct {
	db *gorm.DB
}

func NewGymRepository(db *gorm.DB) *GymRepository {
	return &GymRepository{db: db}
}

// Create a new gym
func (r *GymRepository) Create(gym *models.Gym) error {
	return r.db.Create(gym).Error
}

// List all gyms
func (r *GymRepository) ListAll(gyms *[]models.Gym) error {
	return r.db.Find(gyms).Error
}

// Get gym by ID
func (r *GymRepository) GetByID(id string) (*models.Gym, error) {
	var gym models.Gym
	err := r.db.First(&gym, "id = ?", id).Error
	return &gym, err
}

// Update gym
func (r *GymRepository) Update(gym *models.Gym) error {
	return r.db.Save(gym).Error
}
