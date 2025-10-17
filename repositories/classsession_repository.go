package repositories

import (
	"go-blog/internal/models"

	"gorm.io/gorm"
)

type ClassSessionRepository struct {
	db *gorm.DB
}

func NewClassSessionRepository(db *gorm.DB) *ClassSessionRepository {
	return &ClassSessionRepository{db: db}
}

// Create a new session
func (r *ClassSessionRepository) Create(session *models.ClassSession) error {
	return r.db.Create(session).Error
}

// List all sessions
func (r *ClassSessionRepository) ListAll(sessions *[]models.ClassSession) error {
	return r.db.Preload("Class").Find(sessions).Error
}

// Get a session by ID
func (r *ClassSessionRepository) GetByID(id string) (*models.ClassSession, error) {
	var session models.ClassSession
	err := r.db.Preload("Class").First(&session, "id = ?", id).Error
	return &session, err
}
