package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/ieraHQ/Vakalat/backend/api/repositories"
)

// HearingService defines the interface for hearing operations.
type HearingService interface {
	CreateHearing(ctx context.Context, hearing *repositories.Hearing) error
	GetHearingByID(ctx context.Context, id string) (*repositories.Hearing, error)
	UpdateHearing(ctx context.Context, hearing *repositories.Hearing) error
	DeleteHearing(ctx context.Context, id string) error
	ListHearingsByMatter(ctx context.Context, matterID string) ([]*repositories.Hearing, error)
}

// hearingService implements HearingService.
type hearingService struct {
	hearingRepo repositories.HearingRepository
}

// NewHearingService creates a new HearingService.
func NewHearingService(hearingRepo repositories.HearingRepository) HearingService {
	return &hearingService{hearingRepo: hearingRepo}
}

// CreateHearing creates a new hearing.
func (s *hearingService) CreateHearing(ctx context.Context, hearing *repositories.Hearing) error {
	hearing.ID = uuid.New().String()
	return s.hearingRepo.Create(ctx, hearing)
}

// GetHearingByID retrieves a hearing by ID.
func (s *hearingService) GetHearingByID(ctx context.Context, id string) (*repositories.Hearing, error) {
	return s.hearingRepo.FindByID(ctx, id)
}

// UpdateHearing updates a hearing.
func (s *hearingService) UpdateHearing(ctx context.Context, hearing *repositories.Hearing) error {
	return s.hearingRepo.Update(ctx, hearing)
}

// DeleteHearing soft-deletes a hearing.
func (s *hearingService) DeleteHearing(ctx context.Context, id string) error {
	return s.hearingRepo.Delete(ctx, id)
}

// ListHearingsByMatter retrieves all hearings for a matter.
func (s *hearingService) ListHearingsByMatter(ctx context.Context, matterID string) ([]*repositories.Hearing, error) {
	return s.hearingRepo.ListByMatter(ctx, matterID)
}
