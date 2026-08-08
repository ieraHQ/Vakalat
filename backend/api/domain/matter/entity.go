package matter

import (
	"time"
)

// Matter represents a legal matter in the domain layer.
type Matter struct {
	ID              string
	Title           string
	Description     string
	ClientID        string
	CourtID         string
	JudgeID         string
	AdvocateID      string
	CaseNumber      string
	CaseType        string
	Stage           string
	Status          string
	Priority        PriorityLevel
	LimitationDate  time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}

// PriorityLevel represents the priority level of a matter.
type PriorityLevel string

const (
	PriorityLow     PriorityLevel = "low"
	PriorityMedium  PriorityLevel = "medium"
	PriorityHigh    PriorityLevel = "high"
	PriorityUrgent  PriorityLevel = "urgent"
)

// MatterParty represents a party involved in a legal matter.
type MatterParty struct {
	ID          string
	MatterID    string
	PartyType   PartyType
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

// PartyType represents the type of party in a legal matter.
type PartyType string

const (
	PartyPetitioner  PartyType = "petitioner"
	PartyRespondent  PartyType = "respondent"
	PartyIntervenor  PartyType = "intervenor"
)

// Hearing represents a court hearing for a legal matter.
type Hearing struct {
	ID               string
	MatterID         string
	Date             time.Time
	Notes            string
	Outcome          string
	NextHearingDate  *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time
}

// Order represents a court order related to a legal matter.
type Order struct {
	ID          string
	MatterID    string
	HearingID   *string
	Title       string
	Description string
	Date        time.Time
	DocumentID  *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}