package services

import (
	"go-blog/repositories"
	"time"

	"go-blog/internal/models"

	"github.com/google/uuid"
)

type GymService struct {
	repo *repositories.GymRepository
}

func NewGymService(repo *repositories.GymRepository) *GymService {
	return &GymService{repo: repo}
}

// Create a new gym
func (s *GymService) CreateGym(name, address, phone, timezone string, openingHours, settings map[string]interface{}) (*models.Gym, error) {
	gym := &models.Gym{
		ID:           uuid.New(),
		Name:         name,
		Address:      address,
		Phone:        phone,
		Timezone:     timezone,
		OpeningHours: models.MapToJSON(openingHours),
		Settings:     models.MapToJSON(settings),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	err := s.repo.Create(gym)
	return gym, err
}

// List all gyms
func (s *GymService) ListGyms() ([]models.Gym, error) {
	var gyms []models.Gym
	err := s.repo.ListAll(&gyms)
	return gyms, err
}

// Get gym by ID
func (s *GymService) GetGym(id uuid.UUID) (*models.Gym, error) {
	return s.repo.GetByID(id.String())
}

// Update gym
func (s *GymService) UpdateGym(gym *models.Gym) error {
	gym.UpdatedAt = time.Now()
	return s.repo.Update(gym)
}
