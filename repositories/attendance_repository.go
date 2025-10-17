package repositories

import (
	"go-blog/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AttendanceRepository struct {
	db *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) *AttendanceRepository {
	return &AttendanceRepository{db: db}
}

func (r *AttendanceRepository) Create(attendance *models.Attendance) error {
	return r.db.Create(attendance).Error
}

func (r *AttendanceRepository) FindByMemberAndSession(memberID, sessionID uuid.UUID) (*models.Attendance, error) {
	var record models.Attendance
	err := r.db.Where("member_id = ? AND session_id = ?", memberID, sessionID).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *AttendanceRepository) FindAllByMember(memberID uuid.UUID) ([]models.Attendance, error) {
	var records []models.Attendance
	err := r.db.Preload("Session").Where("member_id = ?", memberID).Order("checked_in_at desc").Find(&records).Error
	return records, err
}

func (r *AttendanceRepository) FindAll() ([]models.Attendance, error) {
	var records []models.Attendance
	err := r.db.Preload("Member").Preload("Session").Find(&records).Error
	return records, err
}
