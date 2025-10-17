package repositories

import (
	"go-blog/internal/models"

	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Create(p *models.Payment) error {
	return r.db.Create(p).Error
}

func (r *PaymentRepository) GetByMember(memberID string) ([]models.Payment, error) {
	var payments []models.Payment
	err := r.db.Where("member_id = ?", memberID).Find(&payments).Error
	return payments, err
}
func (r *PaymentRepository) GetAll() ([]models.Payment, error) {
	var payments []models.Payment
	err := r.db.Find(&payments).Error
	return payments, err
}
