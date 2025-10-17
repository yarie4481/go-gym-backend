package repositories

import (
	"go-blog/internal/models"

	"gorm.io/gorm"
)

type MembershipRepository struct {
	db *gorm.DB
}

func NewMembershipRepository(db *gorm.DB) *MembershipRepository {
	return &MembershipRepository{db: db}
}

func (r *MembershipRepository) Create(m *models.Membership) error {
	return r.db.Create(m).Error
}

func (r *MembershipRepository) GetByMember(memberID string) ([]models.Membership, error) {
	var memberships []models.Membership
	err := r.db.Where("member_id = ?", memberID).Preload("Plan").Find(&memberships).Error
	return memberships, err
}

func (r *MembershipRepository) Update(m *models.Membership) error {
	return r.db.Save(m).Error
}

func (r *MembershipRepository) Delete(id string) error {
	return r.db.Delete(&models.Membership{}, "id = ?", id).Error
}
