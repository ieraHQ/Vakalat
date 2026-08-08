package repositories

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// AISession represents a single prompt/response exchange with the AI assistant.
type AISession struct {
	ID        string    `json:"id"`
	MatterID  string    `json:"matter_id"`
	Prompt    string    `json:"prompt"`
	Response  string    `json:"response"`
	Model     string    `json:"model"`
	CreatedAt time.Time `json:"created_at"`
}

// AISummary represents a generated summary for a matter.
type AISummary struct {
	ID        string    `json:"id"`
	MatterID  string    `json:"matter_id"`
	Summary   string    `json:"summary"`
	CreatedAt time.Time `json:"created_at"`
}

// AIRepository defines the interface for persisting AI sessions and summaries.
type AIRepository interface {
	CreateSession(ctx context.Context, session *AISession) error
	ListSessionsByMatter(ctx context.Context, matterID string) ([]*AISession, error)
	CreateSummary(ctx context.Context, summary *AISummary) error
	ListSummariesByMatter(ctx context.Context, matterID string) ([]*AISummary, error)
}

// aiRepository implements AIRepository.
type aiRepository struct {
	db *pgxpool.Pool
}

// NewAIRepository creates a new AIRepository.
func NewAIRepository(db *pgxpool.Pool) AIRepository {
	return &aiRepository{db: db}
}

// CreateSession inserts a new AI session record.
func (r *aiRepository) CreateSession(ctx context.Context, session *AISession) error {
	query := `INSERT INTO ai_sessions (id, matter_id, prompt, response, model) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(ctx, query, session.ID, session.MatterID, session.Prompt, session.Response, session.Model)
	return err
}

// ListSessionsByMatter retrieves all AI sessions for a matter, most recent first.
func (r *aiRepository) ListSessionsByMatter(ctx context.Context, matterID string) ([]*AISession, error) {
	query := `SELECT id, matter_id, prompt, response, model, created_at FROM ai_sessions WHERE matter_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query, matterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*AISession
	for rows.Next() {
		var session AISession
		if err := rows.Scan(&session.ID, &session.MatterID, &session.Prompt, &session.Response, &session.Model, &session.CreatedAt); err != nil {
			return nil, err
		}
		sessions = append(sessions, &session)
	}
	return sessions, nil
}

// CreateSummary inserts a new AI-generated summary record.
func (r *aiRepository) CreateSummary(ctx context.Context, summary *AISummary) error {
	query := `INSERT INTO ai_summaries (id, matter_id, summary) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, query, summary.ID, summary.MatterID, summary.Summary)
	return err
}

// ListSummariesByMatter retrieves all AI summaries for a matter, most recent first.
func (r *aiRepository) ListSummariesByMatter(ctx context.Context, matterID string) ([]*AISummary, error) {
	query := `SELECT id, matter_id, summary, created_at FROM ai_summaries WHERE matter_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query, matterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []*AISummary
	for rows.Next() {
		var summary AISummary
		if err := rows.Scan(&summary.ID, &summary.MatterID, &summary.Summary, &summary.CreatedAt); err != nil {
			return nil, err
		}
		summaries = append(summaries, &summary)
	}
	return summaries, nil
}
