package services

import (
	"go-blog/internal/models"
	"go-blog/repositories"

	"github.com/google/uuid"
)

type MemberService interface {
	CreateMember(member *models.Member) error
	GetAllMembers() ([]models.Member, error)
	GetMemberByID(id uuid.UUID) (*models.Member, error)
	UpdateMember(member *models.Member) error
	DeleteMember(id uuid.UUID) error
}

type memberService struct {
	repo repositories.MemberRepository
}

func NewMemberService(repo repositories.MemberRepository) MemberService {
	return &memberService{repo: repo}
}

func (s *memberService) CreateMember(member *models.Member) error {
	return s.repo.Create(member)
}

func (s *memberService) GetAllMembers() ([]models.Member, error) {
	return s.repo.GetAll()
}

func (s *memberService) GetMemberByID(id uuid.UUID) (*models.Member, error) {
	return s.repo.GetByID(id)
}

func (s *memberService) UpdateMember(member *models.Member) error {
	return s.repo.Update(member)
}

func (s *memberService) DeleteMember(id uuid.UUID) error {
	return s.repo.Delete(id)
}
