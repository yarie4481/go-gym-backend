package repositories

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"go-blog/internal/models"
)

type MemberRepository interface {
	Create(member *models.Member) error
	GetAll() ([]models.Member, error)
	GetByID(id uuid.UUID) (*models.Member, error)
	Update(member *models.Member) error
	Delete(id uuid.UUID) error
}

type memberRepository struct {
	db *gorm.DB
}

func NewMemberRepository(db *gorm.DB) MemberRepository {
	return &memberRepository{db: db}
}

func (r *memberRepository) Create(member *models.Member) error {
	return r.db.Create(member).Error
}

func (r *memberRepository) GetAll() ([]models.Member, error) {
	var members []models.Member
	err := r.db.
		Preload("Memberships").
		Preload("Bookings").
		Preload("Attendance").
		Preload("Payments").
		Find(&members).Error
	return members, err
}

func (r *memberRepository) GetByID(id uuid.UUID) (*models.Member, error) {
	var member models.Member
	err := r.db.Preload("User").First(&member, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *memberRepository) Update(member *models.Member) error {
	return r.db.Save(member).Error
}

func (r *memberRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&models.Member{}, "id = ?", id)
	if result.RowsAffected == 0 {
		return errors.New("member not found")
	}
	return result.Error
}
