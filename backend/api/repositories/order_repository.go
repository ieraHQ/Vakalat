package repositories

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Order represents a court order related to a legal matter.
type Order struct {
	ID          string    `json:"id"`
	MatterID    string    `json:"matter_id" validate:"required,uuid"`
	HearingID   *string   `json:"hearing_id"`
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description"`
	Date        time.Time `json:"date" validate:"required"`
	DocumentID  *string   `json:"document_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// OrderRepository defines the interface for order persistence.
type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	FindByID(ctx context.Context, id string) (*Order, error)
	Update(ctx context.Context, order *Order) error
	Delete(ctx context.Context, id string) error
	ListByMatter(ctx context.Context, matterID string) ([]*Order, error)
}

// orderRepository implements OrderRepository.
type orderRepository struct {
	db *pgxpool.Pool
}

// NewOrderRepository creates a new OrderRepository.
func NewOrderRepository(db *pgxpool.Pool) OrderRepository {
	return &orderRepository{db: db}
}

// Create inserts a new order into the database.
func (r *orderRepository) Create(ctx context.Context, order *Order) error {
	query := `
		INSERT INTO orders (id, matter_id, hearing_id, title, description, date, document_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(ctx, query, order.ID, order.MatterID, order.HearingID, order.Title, order.Description, order.Date, order.DocumentID)
	return err
}

// FindByID retrieves an order by ID.
func (r *orderRepository) FindByID(ctx context.Context, id string) (*Order, error) {
	query := `SELECT id, matter_id, hearing_id, title, description, date, document_id, created_at, updated_at FROM orders WHERE id = $1 AND deleted_at IS NULL`
	row := r.db.QueryRow(ctx, query, id)

	var order Order
	err := row.Scan(&order.ID, &order.MatterID, &order.HearingID, &order.Title, &order.Description, &order.Date, &order.DocumentID, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

// Update updates an order in the database.
func (r *orderRepository) Update(ctx context.Context, order *Order) error {
	query := `
		UPDATE orders
		SET hearing_id = $1, title = $2, description = $3, date = $4, document_id = $5, updated_at = NOW()
		WHERE id = $6
	`
	_, err := r.db.Exec(ctx, query, order.HearingID, order.Title, order.Description, order.Date, order.DocumentID, order.ID)
	return err
}

// Delete soft-deletes an order from the database.
func (r *orderRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE orders SET deleted_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// ListByMatter retrieves all orders for a matter, most recent first.
func (r *orderRepository) ListByMatter(ctx context.Context, matterID string) ([]*Order, error) {
	query := `SELECT id, matter_id, hearing_id, title, description, date, document_id, created_at, updated_at FROM orders WHERE matter_id = $1 AND deleted_at IS NULL ORDER BY date DESC`
	rows, err := r.db.Query(ctx, query, matterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*Order
	for rows.Next() {
		var order Order
		if err := rows.Scan(&order.ID, &order.MatterID, &order.HearingID, &order.Title, &order.Description, &order.Date, &order.DocumentID, &order.CreatedAt, &order.UpdatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}

	return orders, nil
}
