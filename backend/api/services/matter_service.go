package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ieraHQ/Vakalat/backend/api/ai"
	"github.com/ieraHQ/Vakalat/backend/api/logger"
	"github.com/ieraHQ/Vakalat/backend/api/repositories"
	"go.uber.org/zap"
)

// MatterService defines the interface for matter operations.
type MatterService interface {
	CreateMatter(ctx context.Context, matter *repositories.Matter) error
	GetMatterByID(ctx context.Context, id string) (*repositories.Matter, error)
	UpdateMatter(ctx context.Context, matter *repositories.Matter) error
	DeleteMatter(ctx context.Context, id string) error
	ListMattersByClient(ctx context.Context, clientID string, limit, offset int) ([]*repositories.Matter, error)
	ListMattersByAdvocate(ctx context.Context, advocateID string, limit, offset int) ([]*repositories.Matter, error)
	ListMatters(ctx context.Context, limit, offset int) ([]*repositories.Matter, error)
}

// matterService implements MatterService.
type matterService struct {
	matterRepo    repositories.MatterRepository
	searchService SearchService
	embedder      ai.LLMClient
}

// NewMatterService creates a new MatterService. searchService and embedder
// keep the search index up to date whenever a matter is created or updated —
// indexing is best-effort and never fails the calling request.
func NewMatterService(matterRepo repositories.MatterRepository, searchService SearchService, embedder ai.LLMClient) MatterService {
	return &matterService{matterRepo: matterRepo, searchService: searchService, embedder: embedder}
}

// CreateMatter creates a new matter.
func (s *matterService) CreateMatter(ctx context.Context, matter *repositories.Matter) error {
	matter.ID = uuid.New().String()
	if err := s.matterRepo.Create(ctx, matter); err != nil {
		return err
	}
	s.indexMatter(ctx, matter)
	return nil
}

// GetMatterByID retrieves a matter by ID.
func (s *matterService) GetMatterByID(ctx context.Context, id string) (*repositories.Matter, error) {
	return s.matterRepo.FindByID(ctx, id)
}

// UpdateMatter updates a matter.
func (s *matterService) UpdateMatter(ctx context.Context, matter *repositories.Matter) error {
	if err := s.matterRepo.Update(ctx, matter); err != nil {
		return err
	}
	s.indexMatter(ctx, matter)
	return nil
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

// ListMatters retrieves matters with pagination, regardless of client/advocate.
func (s *matterService) ListMatters(ctx context.Context, limit, offset int) ([]*repositories.Matter, error) {
	return s.matterRepo.ListAll(ctx, limit, offset)
}

// indexMatter embeds and upserts a matter's searchable content. Failures are
// logged, not returned — search indexing must never block a matter write.
func (s *matterService) indexMatter(ctx context.Context, matter *repositories.Matter) {
	content := fmt.Sprintf("%s %s %s", matter.Title, matter.Description, matter.CaseNumber)

	embedding, err := s.embedder.Embed(ctx, content)
	if err != nil {
		logger.GetLogger().Warn("Failed to embed matter for search index", zap.String("matter_id", matter.ID), zap.Error(err))
		embedding = nil
	}

	if err := s.searchService.IndexContent(ctx, "matter", matter.ID, content, embedding); err != nil {
		logger.GetLogger().Warn("Failed to index matter", zap.String("matter_id", matter.ID), zap.Error(err))
	}
}
