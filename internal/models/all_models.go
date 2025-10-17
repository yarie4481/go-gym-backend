package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// The User model MUST be defined first due to the dependency chain.
type User struct {
	UserID                  uuid.UUID  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"user_id"`
	FirstName               string     `gorm:"size:50;not null" json:"first_name"`
	LastName                string     `gorm:"size:50;not null" json:"last_name"`
	Email                   string     `gorm:"size:100;unique;not null" json:"email"`
	PasswordHash            string     `gorm:"size:255;not null" json:"-"`
	PhoneNumber             string     `gorm:"size:15" json:"phone_number"`
	DateOfBirth             *time.Time `json:"date_of_birth"`
	MembershipType          string     `gorm:"type:varchar(20);not null" json:"membership_type"`
	JoinDate                time.Time  `gorm:"default:CURRENT_DATE" json:"join_date"`
	MembershipStart         *time.Time `json:"membership_start"`
	MembershipEnd           *time.Time `json:"membership_end"`
	FitnessGoals            string     `json:"fitness_goals"`
	EmergencyContactName    string     `gorm:"size:100" json:"emergency_contact_name"`
	EmergencyContactPhone   string     `gorm:"size:15" json:"emergency_contact_phone"`
	Status                  string     `gorm:"type:varchar(20);default:'Active'" json:"status"`
	UserType                string     `gorm:"type:varchar(20);default:'Member'" json:"user_type"`
	ProfilePictureURL       string     `gorm:"size:255" json:"profile_picture_url"`
	CreatedAt               time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt               time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	Gender           string

	// Relationships
	// Member (One-to-One): Added a pointer to prevent circular reference/migration issues.
	// GORM will use Member.UserID as the foreign key.
	Member *Member `gorm:"foreignKey:UserID"` 
}

// Member model - profile extension for role = member
type Member struct {
	ID               uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	FirstName               string     `gorm:"size:50;not null" json:"first_name"`
	LastName               string     `gorm:"size:50;not null" json:"last_name"`
	UserID           uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex"`
	Dob              *time.Time
	Gender           string
	EmergencyContact datatypes.JSON `gorm:"type:jsonb"`
	Notes            string         `gorm:"type:text"`
	CreatedAt        time.Time
	UpdatedAt        time.Time

	// Relationships
	// User (One-to-One): Corrected the constraint tag. The foreign key is the UserID on this table.
	User        User         `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"` // <-- FIX IS HERE
	Memberships []Membership `gorm:"foreignKey:MemberID"`
	Bookings    []Booking    `gorm:"foreignKey:MemberID"`
	Attendance  []Attendance `gorm:"foreignKey:MemberID"`
	Payments    []Payment    `gorm:"foreignKey:MemberID"`
}

// Gym model
type Gym struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name           string         `gorm:"not null"`
	Address        string
	Phone          string
	Timezone       string
	OpeningHours   datatypes.JSON `gorm:"type:jsonb"`
	Settings       datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt      time.Time
	UpdatedAt      time.Time

	// Relationships
	Classes        []Class         `gorm:"foreignKey:GymID"`
	InventoryItems []InventoryItem `gorm:"foreignKey:GymID"`
}

// Plan model
type Plan struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Title        string    `gorm:"not null"`
	Description  string    `gorm:"type:text"`
	PriceCents   int       `gorm:"not null"`
	BillingCycle string    `gorm:"not null"`
	NumSessions  *int      // null = unlimited
	Access       string    `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time

	// Relationships
	Memberships []Membership `gorm:"foreignKey:PlanID"`
}

// Membership model
type Membership struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	    MemberID uuid.UUID `gorm:"type:uuid;not null"`
    PlanID   uuid.UUID `gorm:"type:uuid;not null"`
	StartDate       time.Time `gorm:"not null"`
	EndDate         time.Time `gorm:"not null"`
	Status          string    `gorm:"not null;default:'active'"`
	AutoRenew       bool      `gorm:"default:true"`
	PaymentMethodID *string
	CreatedAt       time.Time
	UpdatedAt       time.Time



    // ...
    Member Member `gorm:"foreignKey:MemberID"`
    Plan   Plan   `gorm:"foreignKey:PlanID"`



}

// Class model (activities)
type Class struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	GymID           uuid.UUID      `gorm:"type:uuid;not null"`
	Title           string         `gorm:"not null"`
	Description     string         `gorm:"type:text"`
	TrainerID       uuid.UUID      `gorm:"type:uuid;not null"`
	Capacity        int            `gorm:"not null"`
	RecurringRule   datatypes.JSON `gorm:"type:jsonb"`
	DurationMinutes int            `gorm:"not null"`
	CreatedAt       time.Time
	UpdatedAt       time.Time

	// Relationships
	Gym      Gym            `gorm:"foreignKey:GymID"`
	Trainer  User           `gorm:"foreignKey:TrainerID;references:UserID"` // Trainer is a User
	Sessions []ClassSession `gorm:"foreignKey:ClassID"`
}

// ClassSession model (concrete occurrence)
type ClassSession struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ClassID   uuid.UUID `gorm:"type:uuid;not null"`
	StartsAt  time.Time `gorm:"not null"`
	EndsAt    time.Time `gorm:"not null"`
	Capacity  int       `gorm:"not null"`
	Status    string    `gorm:"not null;default:'scheduled'"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// Relationships
	Class      Class        `gorm:"foreignKey:ClassID"`
	Bookings   []Booking    `gorm:"foreignKey:SessionID"`
	Attendance []Attendance `gorm:"foreignKey:SessionID"`
}

// Booking model
type Booking struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SessionID uuid.UUID `gorm:"type:uuid;not null"`
	MemberID  uuid.UUID `gorm:"type:uuid;not null"`
	Status    string    `gorm:"not null;default:'booked'"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// Relationships
	Session ClassSession `gorm:"foreignKey:SessionID"`
	Member  Member       `gorm:"foreignKey:MemberID"`
}

// Attendance model
type Attendance struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SessionID       uuid.UUID `gorm:"type:uuid;not null"`
	MemberID        uuid.UUID `gorm:"type:uuid;not null"`
	CheckinMethod   string    `gorm:"not null"`
	CheckedInAt     time.Time `gorm:"not null"`
	CreatedAt       time.Time

	// Relationships
	Session ClassSession `gorm:"foreignKey:SessionID"`
	Member  Member       `gorm:"foreignKey:MemberID"`
}

// Payment model
type Payment struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	MemberID    uuid.UUID `gorm:"type:uuid;not null"`
	AmountCents int       `gorm:"not null"`
	Currency    string    `gorm:"not null;default:'Birr'"`
	Method      string    `gorm:"not null"`
	Status      string    `gorm:"not null;default:'pending'"`
	Reference   string
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// Relationships
	Member Member `gorm:"foreignKey:MemberID"`
}

// InventoryItem model
type InventoryItem struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	GymID            uuid.UUID `gorm:"type:uuid;not null"`
	Name             string    `gorm:"not null"`
	SKU              string    `gorm:"uniqueIndex;not null"`
	Quantity         int       `gorm:"not null;default:0"`
	ReorderThreshold int       `gorm:"not null;default:5"`
	MaintenanceNotes string    `gorm:"type:text"`
	CreatedAt        time.Time
	UpdatedAt        time.Time

	// Relationships
	Gym Gym `gorm:"foreignKey:GymID"`
}

// Notification model
type Notification struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null"`
	Type        string         `gorm:"not null"`
	Payload     datatypes.JSON `gorm:"type:jsonb"`
	SentAt      *time.Time
	DeliveredAt *time.Time
	CreatedAt   time.Time

	// User is a one-to-many relationship (one user, many notifications)
	// You may or may not explicitly define the User relationship here if it's not needed.
}

// AuditLog model
type AuditLog struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ActorUserID uuid.UUID      `gorm:"type:uuid;not null"`
	ActionType  string         `gorm:"not null"`
	TargetType  string         `gorm:"not null"`
	TargetID    uuid.UUID      `gorm:"type:uuid;not null"`
	Metadata    datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt   time.Time

	Actor User `gorm:"foreignKey:ActorUserID;references:UserID"`
}

func MigrateModels(db *gorm.DB) {
	// Make sure pgcrypto extension exists before anything else
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto";`).Error; err != nil {
		panic("❌ Failed to enable pgcrypto extension: " + err.Error())
	}

	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)

	// Explicit order migration is crucial to avoid foreign key issues
	models := []interface{}{
		&User{},            // 1. MUST come first (no foreign keys)
		&Member{},          // 2. Depends on User
		&Gym{},             // 3. Independent
		&Plan{},            // 4. Independent
		&Membership{},      // 5. Depends on Member, Plan
		&Class{},           // 6. Depends on Gym, User (as Trainer)
		&ClassSession{},    // 7. Depends on Class
		&Booking{},         // 8. Depends on ClassSession, Member
		&Attendance{},      // 9. Depends on ClassSession, Member
		&Payment{},         // 10. Depends on Member
		&InventoryItem{},   // 11. Depends on Gym
		&Notification{},    // 12. Depends on User
		&AuditLog{},        // 13. Depends on User
	}

	for _, m := range models {
		if err := db.AutoMigrate(m); err != nil {
			panic(fmt.Sprintf("❌ Failed to migrate model %T: %v", m, err))
		}
	}

	fmt.Println("✅ All database migrations completed successfully!")
}