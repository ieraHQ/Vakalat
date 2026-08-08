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

// ClientService defines the interface for client operations.
type ClientService interface {
	CreateClient(ctx context.Context, client *repositories.Client) error
	GetClientByID(ctx context.Context, id string) (*repositories.Client, error)
	UpdateClient(ctx context.Context, client *repositories.Client) error
	DeleteClient(ctx context.Context, id string) error
	ListClients(ctx context.Context, limit, offset int) ([]*repositories.Client, error)
}

// clientService implements ClientService.
type clientService struct {
	clientRepo    repositories.ClientRepository
	searchService SearchService
	embedder      ai.LLMClient
}

// NewClientService creates a new ClientService. searchService and embedder
// keep the search index up to date whenever a client is created or updated —
// indexing is best-effort and never fails the calling request.
func NewClientService(clientRepo repositories.ClientRepository, searchService SearchService, embedder ai.LLMClient) ClientService {
	return &clientService{clientRepo: clientRepo, searchService: searchService, embedder: embedder}
}

// CreateClient creates a new client.
func (s *clientService) CreateClient(ctx context.Context, client *repositories.Client) error {
	client.ID = uuid.New().String()
	if err := s.clientRepo.Create(ctx, client); err != nil {
		return err
	}
	s.indexClient(ctx, client)
	return nil
}

// GetClientByID retrieves a client by ID.
func (s *clientService) GetClientByID(ctx context.Context, id string) (*repositories.Client, error) {
	return s.clientRepo.FindByID(ctx, id)
}

// UpdateClient updates a client.
func (s *clientService) UpdateClient(ctx context.Context, client *repositories.Client) error {
	if err := s.clientRepo.Update(ctx, client); err != nil {
		return err
	}
	s.indexClient(ctx, client)
	return nil
}

// DeleteClient soft-deletes a client.
func (s *clientService) DeleteClient(ctx context.Context, id string) error {
	return s.clientRepo.Delete(ctx, id)
}

// ListClients retrieves a list of clients with pagination.
func (s *clientService) ListClients(ctx context.Context, limit, offset int) ([]*repositories.Client, error) {
	return s.clientRepo.List(ctx, limit, offset)
}

// indexClient embeds and upserts a client's searchable content. Failures are
// logged, not returned — search indexing must never block a client write.
func (s *clientService) indexClient(ctx context.Context, client *repositories.Client) {
	notes := ""
	if client.Notes != nil {
		notes = *client.Notes
	}
	content := fmt.Sprintf("%s %s %s %s", client.Name, client.Email, client.Phone, notes)

	embedding, err := s.embedder.Embed(ctx, content)
	if err != nil {
		logger.GetLogger().Warn("Failed to embed client for search index", zap.String("client_id", client.ID), zap.Error(err))
		embedding = nil
	}

	if err := s.searchService.IndexContent(ctx, "client", client.ID, content, embedding); err != nil {
		logger.GetLogger().Warn("Failed to index client", zap.String("client_id", client.ID), zap.Error(err))
	}
}
