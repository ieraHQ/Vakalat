package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/ieraHQ/Vakalat/backend/api/repositories"
)

// OrderService defines the interface for order operations.
type OrderService interface {
	CreateOrder(ctx context.Context, order *repositories.Order) error
	GetOrderByID(ctx context.Context, id string) (*repositories.Order, error)
	UpdateOrder(ctx context.Context, order *repositories.Order) error
	DeleteOrder(ctx context.Context, id string) error
	ListOrdersByMatter(ctx context.Context, matterID string) ([]*repositories.Order, error)
}

// orderService implements OrderService.
type orderService struct {
	orderRepo repositories.OrderRepository
}

// NewOrderService creates a new OrderService.
func NewOrderService(orderRepo repositories.OrderRepository) OrderService {
	return &orderService{orderRepo: orderRepo}
}

// CreateOrder creates a new order.
func (s *orderService) CreateOrder(ctx context.Context, order *repositories.Order) error {
	order.ID = uuid.New().String()
	return s.orderRepo.Create(ctx, order)
}

// GetOrderByID retrieves an order by ID.
func (s *orderService) GetOrderByID(ctx context.Context, id string) (*repositories.Order, error) {
	return s.orderRepo.FindByID(ctx, id)
}

// UpdateOrder updates an order.
func (s *orderService) UpdateOrder(ctx context.Context, order *repositories.Order) error {
	return s.orderRepo.Update(ctx, order)
}

// DeleteOrder soft-deletes an order.
func (s *orderService) DeleteOrder(ctx context.Context, id string) error {
	return s.orderRepo.Delete(ctx, id)
}

// ListOrdersByMatter retrieves all orders for a matter.
func (s *orderService) ListOrdersByMatter(ctx context.Context, matterID string) ([]*repositories.Order, error) {
	return s.orderRepo.ListByMatter(ctx, matterID)
}
