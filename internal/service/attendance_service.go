package services

import (
	"errors"
	"go-blog/repositories"
	"time"

	"go-blog/internal/models"

	"github.com/google/uuid"
)

type AttendanceService struct {
	repo *repositories.AttendanceRepository
}

func NewAttendanceService(repo *repositories.AttendanceRepository) *AttendanceService {
	return &AttendanceService{repo: repo}
}

// ✅ Check-in logic
func (s *AttendanceService) CheckIn(memberID, sessionID uuid.UUID, method string) error {
	// prevent double check-in
	existing, _ := s.repo.FindByMemberAndSession(memberID, sessionID)
	if existing != nil {
		return errors.New("member already checked in for this session")
	}

	record := &models.Attendance{
		ID:            uuid.New(),
		MemberID:      memberID,
		SessionID:     sessionID,
		CheckinMethod: method,
		CheckedInAt:   time.Now(),
	}

	return s.repo.Create(record)
}

// ✅ Get attendance by member
func (s *AttendanceService) GetMemberAttendance(memberID uuid.UUID) ([]models.Attendance, error) {
	return s.repo.FindAllByMember(memberID)
}

// ✅ Get all attendance (for admin)
func (s *AttendanceService) GetAllAttendance() ([]models.Attendance, error) {
	return s.repo.FindAll()
}
