package services

import (
	"context"
	"github.com/ieraHQ/Vakalat/backend/api/repositories"
	"github.com/google/uuid"
)

// MatterService defines the interface for matter operations.
type MatterService interface {
	CreateMatter(ctx context.Context, matter *repositories.Matter) error
	GetMatterByID(ctx context.Context, id string) (*repositories.Matter, error)
	UpdateMatter(ctx context.Context, matter *repositories.Matter) error
	DeleteMatter(ctx context.Context, id string) error
	ListMattersByClient(ctx context.Context, clientID string, limit, offset int) ([]*repositories.Matter, error)
	ListMattersByAdvocate(ctx context.Context, advocateID string, limit, offset int) ([]*repositories.Matter, error)
}

// matterService implements MatterService.
type matterService struct {
	matterRepo repositories.MatterRepository
}

// NewMatterService creates a new MatterService.
func NewMatterService(matterRepo repositories.MatterRepository) MatterService {
	return &matterService{matterRepo: matterRepo}
}

// CreateMatter creates a new matter.
func (s *matterService) CreateMatter(ctx context.Context, matter *repositories.Matter) error {
	matter.ID = uuid.New().String()
	return s.matterRepo.Create(ctx, matter)
}

// GetMatterByID retrieves a matter by ID.
func (s *matterService) GetMatterByID(ctx context.Context, id string) (*repositories.Matter, error) {
	return s.matterRepo.FindByID(ctx, id)
}

// UpdateMatter updates a matter.
func (s *matterService) UpdateMatter(ctx context.Context, matter *repositories.Matter) error {
	return s.matterRepo.Update(ctx, matter)
}

// DeleteMatter soft-deletes a matter.
func (s *matterService) DeleteMatter(ctx context.Context, id string) error {
	return s.matterRepo.Delete(ctx, id)
}

// ListMattersByClient retrieves a list of matters for a client with pagination.
func (s *matterService) ListMattersByClient(ctx context.Context, clientID string, limit, offset int) ([]*repositories.Matter, error) {
	return s.matterRepo.ListByClient(ctx, clientID, limit, offset)
}

// ListMattersByAdvocate retrieves a list of matters for an advocate with pagination.
func (s *matterService) ListMattersByAdvocate(ctx context.Context, advocateID string, limit, offset int) ([]*repositories.Matter, error) {
	return s.matterRepo.ListByAdvocate(ctx, advocateID, limit, offset)
}