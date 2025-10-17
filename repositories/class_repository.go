package repositories

import (
	"go-blog/internal/models"

	"gorm.io/gorm"
)

type ClassRepository struct {
	db *gorm.DB
}

func NewClassRepository(db *gorm.DB) *ClassRepository {
	return &ClassRepository{db: db}
}

// Create a new class
func (r *ClassRepository) Create(class *models.Class) error {
	return r.db.Create(class).Error
}

// List all classes
func (r *ClassRepository) ListAll(classes *[]models.Class) error {
	return r.db.Preload("Gym").Preload("Trainer").Find(classes).Error
}

// Get class by ID
func (r *ClassRepository) GetByID(id string) (*models.Class, error) {
	var class models.Class
	err := r.db.Preload("Gym").Preload("Trainer").Preload("Sessions").First(&class, "id = ?", id).Error
	return &class, err
}
