package repositories

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Hearing represents a court hearing for a legal matter.
type Hearing struct {
	ID              string     `json:"id"`
	MatterID        string     `json:"matter_id" validate:"required,uuid"`
	Date            time.Time  `json:"date" validate:"required"`
	Notes           string     `json:"notes"`
	Outcome         string     `json:"outcome"`
	NextHearingDate *time.Time `json:"next_hearing_date"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// HearingRepository defines the interface for hearing persistence.
type HearingRepository interface {
	Create(ctx context.Context, hearing *Hearing) error
	FindByID(ctx context.Context, id string) (*Hearing, error)
	Update(ctx context.Context, hearing *Hearing) error
	Delete(ctx context.Context, id string) error
	ListByMatter(ctx context.Context, matterID string) ([]*Hearing, error)
}

// hearingRepository implements HearingRepository.
type hearingRepository struct {
	db *pgxpool.Pool
}

// NewHearingRepository creates a new HearingRepository.
func NewHearingRepository(db *pgxpool.Pool) HearingRepository {
	return &hearingRepository{db: db}
}

// Create inserts a new hearing into the database.
func (r *hearingRepository) Create(ctx context.Context, hearing *Hearing) error {
	query := `
		INSERT INTO hearings (id, matter_id, date, notes, outcome, next_hearing_date)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(ctx, query, hearing.ID, hearing.MatterID, hearing.Date, hearing.Notes, hearing.Outcome, hearing.NextHearingDate)
	return err
}

// FindByID retrieves a hearing by ID.
func (r *hearingRepository) FindByID(ctx context.Context, id string) (*Hearing, error) {
	query := `SELECT id, matter_id, date, notes, outcome, next_hearing_date, created_at, updated_at FROM hearings WHERE id = $1 AND deleted_at IS NULL`
	row := r.db.QueryRow(ctx, query, id)

	var hearing Hearing
	err := row.Scan(&hearing.ID, &hearing.MatterID, &hearing.Date, &hearing.Notes, &hearing.Outcome, &hearing.NextHearingDate, &hearing.CreatedAt, &hearing.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &hearing, nil
}

// Update updates a hearing in the database.
func (r *hearingRepository) Update(ctx context.Context, hearing *Hearing) error {
	query := `
		UPDATE hearings
		SET date = $1, notes = $2, outcome = $3, next_hearing_date = $4, updated_at = NOW()
		WHERE id = $5
	`
	_, err := r.db.Exec(ctx, query, hearing.Date, hearing.Notes, hearing.Outcome, hearing.NextHearingDate, hearing.ID)
	return err
}

// Delete soft-deletes a hearing from the database.
func (r *hearingRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE hearings SET deleted_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// ListByMatter retrieves all hearings for a matter, most recent first.
func (r *hearingRepository) ListByMatter(ctx context.Context, matterID string) ([]*Hearing, error) {
	query := `SELECT id, matter_id, date, notes, outcome, next_hearing_date, created_at, updated_at FROM hearings WHERE matter_id = $1 AND deleted_at IS NULL ORDER BY date DESC`
	rows, err := r.db.Query(ctx, query, matterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hearings []*Hearing
	for rows.Next() {
		var hearing Hearing
		if err := rows.Scan(&hearing.ID, &hearing.MatterID, &hearing.Date, &hearing.Notes, &hearing.Outcome, &hearing.NextHearingDate, &hearing.CreatedAt, &hearing.UpdatedAt); err != nil {
			return nil, err
		}
		hearings = append(hearings, &hearing)
	}

	return hearings, nil
}
