package repositories

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

// MatterRepository defines the interface for matter operations.
type MatterRepository interface {
	Create(ctx context.Context, matter *Matter) error
	FindByID(ctx context.Context, id string) (*Matter, error)
	Update(ctx context.Context, matter *Matter) error
	Delete(ctx context.Context, id string) error
	ListByClient(ctx context.Context, clientID string, limit, offset int) ([]*Matter, error)
	ListByAdvocate(ctx context.Context, advocateID string, limit, offset int) ([]*Matter, error)
}

// Matter represents a legal matter in the system.
type Matter struct {
	ID              string `json:"id"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	ClientID        string `json:"client_id"`
	CourtID         string `json:"court_id"`
	JudgeID         string `json:"judge_id"`
	AdvocateID      string `json:"advocate_id"`
	CaseNumber      string `json:"case_number"`
	CaseType        string `json:"case_type"`
	Stage           string `json:"stage"`
	Status          string `json:"status"`
	Priority        string `json:"priority"`
	LimitationDate  string `json:"limitation_date"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

// matterRepository implements MatterRepository.
type matterRepository struct {
	db *pgxpool.Pool
}

// NewMatterRepository creates a new MatterRepository.
func NewMatterRepository(db *pgxpool.Pool) MatterRepository {
	return &matterRepository{db: db}
}

// Create inserts a new matter into the database.
func (r *matterRepository) Create(ctx context.Context, matter *Matter) error {
	query := `
		INSERT INTO matters (id, title, description, client_id, court_id, judge_id, advocate_id, case_number, case_type, stage, status, priority, limitation_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	_, err := r.db.Exec(ctx, query, matter.ID, matter.Title, matter.Description, matter.ClientID, matter.CourtID, matter.JudgeID, matter.AdvocateID, matter.CaseNumber, matter.CaseType, matter.Stage, matter.Status, matter.Priority, matter.LimitationDate)
	return err
}

// FindByID retrieves a matter by ID.
func (r *matterRepository) FindByID(ctx context.Context, id string) (*Matter, error) {
	query := `SELECT id, title, description, client_id, court_id, judge_id, advocate_id, case_number, case_type, stage, status, priority, limitation_date, created_at, updated_at FROM matters WHERE id = $1 AND deleted_at IS NULL`
	row := r.db.QueryRow(ctx, query, id)

	var matter Matter
	err := row.Scan(&matter.ID, &matter.Title, &matter.Description, &matter.ClientID, &matter.CourtID, &matter.JudgeID, &matter.AdvocateID, &matter.CaseNumber, &matter.CaseType, &matter.Stage, &matter.Status, &matter.Priority, &matter.LimitationDate, &matter.CreatedAt, &matter.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &matter, nil
}

// Update updates a matter in the database.
func (r *matterRepository) Update(ctx context.Context, matter *Matter) error {
	query := `
		UPDATE matters
		SET title = $1, description = $2, client_id = $3, court_id = $4, judge_id = $5, advocate_id = $6, case_number = $7, case_type = $8, stage = $9, status = $10, priority = $11, limitation_date = $12, updated_at = NOW()
		WHERE id = $13
	`
	_, err := r.db.Exec(ctx, query, matter.Title, matter.Description, matter.ClientID, matter.CourtID, matter.JudgeID, matter.AdvocateID, matter.CaseNumber, matter.CaseType, matter.Stage, matter.Status, matter.Priority, matter.LimitationDate, matter.ID)
	return err
}

// Delete soft-deletes a matter from the database.
func (r *matterRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE matters SET deleted_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// ListByClient retrieves a list of matters for a client with pagination.
func (r *matterRepository) ListByClient(ctx context.Context, clientID string, limit, offset int) ([]*Matter, error) {
	query := `SELECT id, title, description, client_id, court_id, judge_id, advocate_id, case_number, case_type, stage, status, priority, limitation_date, created_at, updated_at FROM matters WHERE client_id = $1 AND deleted_at IS NULL LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(ctx, query, clientID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matters []*Matter
	for rows.Next() {
		var matter Matter
		if err := rows.Scan(&matter.ID, &matter.Title, &matter.Description, &matter.ClientID, &matter.CourtID, &matter.JudgeID, &matter.AdvocateID, &matter.CaseNumber, &matter.CaseType, &matter.Stage, &matter.Status, &matter.Priority, &matter.LimitationDate, &matter.CreatedAt, &matter.UpdatedAt); err != nil {
			return nil, err
		}
		matters = append(matters, &matter)
	}

	return matters, nil
}

// ListByAdvocate retrieves a list of matters for an advocate with pagination.
func (r *matterRepository) ListByAdvocate(ctx context.Context, advocateID string, limit, offset int) ([]*Matter, error) {
	query := `SELECT id, title, description, client_id, court_id, judge_id, advocate_id, case_number, case_type, stage, status, priority, limitation_date, created_at, updated_at FROM matters WHERE advocate_id = $1 AND deleted_at IS NULL LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(ctx, query, advocateID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matters []*Matter
	for rows.Next() {
		var matter Matter
		if err := rows.Scan(&matter.ID, &matter.Title, &matter.Description, &matter.ClientID, &matter.CourtID, &matter.JudgeID, &matter.AdvocateID, &matter.CaseNumber, &matter.CaseType, &matter.Stage, &matter.Status, &matter.Priority, &matter.LimitationDate, &matter.CreatedAt, &matter.UpdatedAt); err != nil {
			return nil, err
		}
		matters = append(matters, &matter)
	}

	return matters, nil
}