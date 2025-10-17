package services

import (
	"go-blog/internal/models"
	"go-blog/repositories"
)

type PaymentService struct {
	repo *repositories.PaymentRepository
}

func NewPaymentService(repo *repositories.PaymentRepository) *PaymentService {
	return &PaymentService{repo: repo}
}

func (s *PaymentService) RecordPayment(p *models.Payment) error {
	return s.repo.Create(p)
}

func (s *PaymentService) GetPayments(memberID string) ([]models.Payment, error) {
	return s.repo.GetByMember(memberID)
}
func (s *PaymentService) GetAllPayments() ([]models.Payment, error) {
	return s.repo.GetAll()
}
