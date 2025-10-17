package services

import (
	"go-blog/internal/models"
	"go-blog/repositories"
)

type PlanService struct {
	repo *repositories.PlanRepository
}

func NewPlanService(repo *repositories.PlanRepository) *PlanService {
	return &PlanService{repo: repo}
}

func (s *PlanService) CreatePlan(plan *models.Plan) error {
	return s.repo.Create(plan)
}

func (s *PlanService) GetPlans() ([]models.Plan, error) {
	return s.repo.GetAll()
}
