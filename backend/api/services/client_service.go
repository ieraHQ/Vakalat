package services

import (
	"context"
	"github.com/ieraHQ/Vakalat/backend/api/repositories"
	"github.com/google/uuid"
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
	clientRepo repositories.ClientRepository
}

// NewClientService creates a new ClientService.
func NewClientService(clientRepo repositories.ClientRepository) ClientService {
	return &clientService{clientRepo: clientRepo}
}

// CreateClient creates a new client.
func (s *clientService) CreateClient(ctx context.Context, client *repositories.Client) error {
	client.ID = uuid.New().String()
	return s.clientRepo.Create(ctx, client)
}

// GetClientByID retrieves a client by ID.
func (s *clientService) GetClientByID(ctx context.Context, id string) (*repositories.Client, error) {
	return s.clientRepo.FindByID(ctx, id)
}

// UpdateClient updates a client.
func (s *clientService) UpdateClient(ctx context.Context, client *repositories.Client) error {
	return s.clientRepo.Update(ctx, client)
}

// DeleteClient soft-deletes a client.
func (s *clientService) DeleteClient(ctx context.Context, id string) error {
	return s.clientRepo.Delete(ctx, id)
}

// ListClients retrieves a list of clients with pagination.
func (s *clientService) ListClients(ctx context.Context, limit, offset int) ([]*repositories.Client, error) {
	return s.clientRepo.List(ctx, limit, offset)
}