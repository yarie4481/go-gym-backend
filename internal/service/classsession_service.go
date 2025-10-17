package services

import (
	"fmt"
	"go-blog/repositories"
	"time"

	"go-blog/internal/models"

	"github.com/google/uuid"
)

type ClassSessionService struct {
	repo *repositories.ClassSessionRepository
}

func NewClassSessionService(repo *repositories.ClassSessionRepository) *ClassSessionService {
	return &ClassSessionService{repo: repo}
}

// Create a session
func (s *ClassSessionService) CreateSession(classID uuid.UUID, startsAt, endsAt time.Time) (*models.ClassSession, error) {
    
    // 1. Validate the incoming classID to ensure it's not a zero value
    if classID == uuid.Nil {
        return nil, fmt.Errorf("class ID cannot be empty") // Added fmt.Errorf import might be needed
    }

    session := &models.ClassSession{
        // 🔑 THE CRITICAL FIX IS HERE: Assign the passed-in classID
        ClassID:  classID, 
        
        // ID is often set by GORM/DB, but using uuid.New() is fine if you prefer client-side generation.
        // If you remove this line, GORM will use the default UUID function defined in your model/DB.
        ID:       uuid.New(), 
        
        StartsAt: startsAt,
        EndsAt:   endsAt,
        Capacity: 20,
        Status:   "scheduled",
    }
    
    err := s.repo.Create(session)
    return session, err
}


// List all sessions
func (s *ClassSessionService) ListSessions() ([]models.ClassSession, error) {
	var sessions []models.ClassSession
	err := s.repo.ListAll(&sessions)
	return sessions, err
}

// Get a session by ID
func (s *ClassSessionService) GetSession(id uuid.UUID) (*models.ClassSession, error) {
	return s.repo.GetByID(id.String())
}
