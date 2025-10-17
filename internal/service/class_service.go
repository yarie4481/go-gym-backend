package services

import (
	"go-blog/internal/models"
	"go-blog/repositories"
	"time"

	"github.com/google/uuid"
)

type ClassService struct {
	repo *repositories.ClassRepository
}

func NewClassService(repo *repositories.ClassRepository) *ClassService {
	return &ClassService{repo: repo}
}

// Create a new class
func (s *ClassService) CreateClass(gymID, trainerID uuid.UUID, title, description string, capacity, durationMinutes int) (*models.Class, error) {
	class := &models.Class{
		ID:              uuid.New(),
		GymID:           gymID,
		TrainerID:       trainerID,
		Title:           title,
		Description:     description,
		Capacity:        capacity,
		DurationMinutes: durationMinutes,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	err := s.repo.Create(class)
	return class, err
}

// List all classes
func (s *ClassService) ListClasses() ([]models.Class, error) {
	var classes []models.Class
	err := s.repo.ListAll(&classes)
	return classes, err
}

// Get class by ID
func (s *ClassService) GetClass(id uuid.UUID) (*models.Class, error) {
	return s.repo.GetByID(id.String())
}
