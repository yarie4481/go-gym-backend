package services

import (
	"go-blog/internal/models"
	"go-blog/repositories"
)

type MembershipService struct {
	repo *repositories.MembershipRepository
}

func NewMembershipService(repo *repositories.MembershipRepository) *MembershipService {
	return &MembershipService{repo: repo}
}

func (s *MembershipService) CreateMembership(m *models.Membership) error {
    // Ensure you create the UUID on the model if you are not relying on the database default
    // If your GORM model has `gorm:"default:gen_random_uuid()"` and you are using postgres, 
    // you can rely on the database. If not, consider setting it here:
    // m.ID = uuid.New()
	return s.repo.Create(m)
}

func (s *MembershipService) GetByMember(memberID string) ([]models.Membership, error) {
	return s.repo.GetByMember(memberID)
}