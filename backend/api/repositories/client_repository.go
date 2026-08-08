package repositories

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ClientRepository defines the interface for client operations.
type ClientRepository interface {
	Create(ctx context.Context, client *Client) error
	FindByID(ctx context.Context, id string) (*Client, error)
	Update(ctx context.Context, client *Client) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*Client, error)
}

// Client represents a client in the system.
// Address, PAN, GSTIN, and Notes are nullable columns (the seed data and any
// client created without them leaves them NULL) — *string, not string, or
// pgx fails every read with "cannot scan NULL into *string".
type Client struct {
	ID      string  `json:"id"`
	Name    string  `json:"name" validate:"required"`
	Type    string  `json:"type" validate:"required,oneof=individual organization"`
	Email   string  `json:"email" validate:"omitempty,email"`
	Phone   string  `json:"phone"`
	Address *string `json:"address"`
	PAN     *string `json:"pan"`
	GSTIN   *string `json:"gstin"`
	Notes   *string `json:"notes"`
}

// clientRepository implements ClientRepository.
type clientRepository struct {
	db *pgxpool.Pool
}

// NewClientRepository creates a new ClientRepository.
func NewClientRepository(db *pgxpool.Pool) ClientRepository {
	return &clientRepository{db: db}
}

// Create inserts a new client into the database.
func (r *clientRepository) Create(ctx context.Context, client *Client) error {
	query := `
		INSERT INTO clients (id, name, type, email, phone, address, pan, gstin, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(ctx, query, client.ID, client.Name, client.Type, client.Email, client.Phone, client.Address, client.PAN, client.GSTIN, client.Notes)
	return err
}

// FindByID retrieves a client by ID.
func (r *clientRepository) FindByID(ctx context.Context, id string) (*Client, error) {
	query := `SELECT id, name, type, email, phone, address, pan, gstin, notes FROM clients WHERE id = $1 AND deleted_at IS NULL`
	row := r.db.QueryRow(ctx, query, id)

	var client Client
	err := row.Scan(&client.ID, &client.Name, &client.Type, &client.Email, &client.Phone, &client.Address, &client.PAN, &client.GSTIN, &client.Notes)
	if err != nil {
		return nil, err
	}

	return &client, nil
}

// Update updates a client in the database.
func (r *clientRepository) Update(ctx context.Context, client *Client) error {
	query := `
		UPDATE clients
		SET name = $1, type = $2, email = $3, phone = $4, address = $5, pan = $6, gstin = $7, notes = $8, updated_at = NOW()
		WHERE id = $9
	`
	_, err := r.db.Exec(ctx, query, client.Name, client.Type, client.Email, client.Phone, client.Address, client.PAN, client.GSTIN, client.Notes, client.ID)
	return err
}

// Delete soft-deletes a client from the database.
func (r *clientRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE clients SET deleted_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// List retrieves a list of clients with pagination.
func (r *clientRepository) List(ctx context.Context, limit, offset int) ([]*Client, error) {
	query := `SELECT id, name, type, email, phone, address, pan, gstin, notes FROM clients WHERE deleted_at IS NULL LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []*Client
	for rows.Next() {
		var client Client
		if err := rows.Scan(&client.ID, &client.Name, &client.Type, &client.Email, &client.Phone, &client.Address, &client.PAN, &client.GSTIN, &client.Notes); err != nil {
			return nil, err
		}
		clients = append(clients, &client)
	}

	return clients, nil
}