package repositories

import (
	"go-blog/internal/models"

	"gorm.io/gorm"
)

type PlanRepository struct {
	db *gorm.DB
}

func NewPlanRepository(db *gorm.DB) *PlanRepository {
	return &PlanRepository{db: db}
}

func (r *PlanRepository) Create(plan *models.Plan) error {
	return r.db.Create(plan).Error
}

func (r *PlanRepository) GetAll() ([]models.Plan, error) {
	var plans []models.Plan
	err := r.db.Find(&plans).Error
	return plans, err
}

func (r *PlanRepository) GetByID(id string) (*models.Plan, error) {
	var plan models.Plan
	err := r.db.First(&plan, "id = ?", id).Error
	return &plan, err
}

func (r *PlanRepository) Update(plan *models.Plan) error {
	return r.db.Save(plan).Error
}

func (r *PlanRepository) Delete(id string) error {
	return r.db.Delete(&models.Plan{}, "id = ?", id).Error
}
