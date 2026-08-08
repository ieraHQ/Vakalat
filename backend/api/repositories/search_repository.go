package repositories

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// SearchResult represents a single hit from the search index.
type SearchResult struct {
	ID         string  `json:"id"`
	EntityType string  `json:"entity_type"`
	EntityID   string  `json:"entity_id"`
	Content    string  `json:"content"`
	Score      float64 `json:"score"`
}

// SearchRepository defines the interface for full-text and semantic search
// over the search_index table.
type SearchRepository interface {
	// SearchByKeyword runs a PostgreSQL full-text search over indexed content.
	SearchByKeyword(ctx context.Context, query string, limit int) ([]*SearchResult, error)
	// SearchBySimilarity runs a pgvector cosine-similarity search using a
	// precomputed embedding.
	SearchBySimilarity(ctx context.Context, embedding []float32, limit int) ([]*SearchResult, error)
	// Upsert indexes (or re-indexes) a piece of content for an entity.
	Upsert(ctx context.Context, entityType, entityID, content string, embedding []float32) error
}

// searchRepository implements SearchRepository.
type searchRepository struct {
	db *pgxpool.Pool
}

// NewSearchRepository creates a new SearchRepository.
func NewSearchRepository(db *pgxpool.Pool) SearchRepository {
	return &searchRepository{db: db}
}

// vectorLiteral formats a float32 slice as a pgvector text literal, e.g. "[0.1,0.2,0.3]".
func vectorLiteral(embedding []float32) string {
	parts := make([]string, len(embedding))
	for i, v := range embedding {
		parts[i] = strconv.FormatFloat(float64(v), 'f', -1, 32)
	}
	return "[" + strings.Join(parts, ",") + "]"
}

// SearchByKeyword performs a PostgreSQL full-text search over indexed content,
// ranked by relevance.
func (r *searchRepository) SearchByKeyword(ctx context.Context, query string, limit int) ([]*SearchResult, error) {
	sql := `
		SELECT id, entity_type, entity_id, content,
		       ts_rank(to_tsvector('english', content), plainto_tsquery('english', $1)) AS score
		FROM search_index
		WHERE to_tsvector('english', content) @@ plainto_tsquery('english', $1)
		ORDER BY score DESC
		LIMIT $2
	`
	rows, err := r.db.Query(ctx, sql, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanSearchResults(rows)
}

// SearchBySimilarity performs a pgvector cosine-distance nearest-neighbor
// search against the given embedding.
func (r *searchRepository) SearchBySimilarity(ctx context.Context, embedding []float32, limit int) ([]*SearchResult, error) {
	sql := fmt.Sprintf(`
		SELECT id, entity_type, entity_id, content,
		       1 - (embedding <=> '%s'::vector) AS score
		FROM search_index
		WHERE embedding IS NOT NULL
		ORDER BY embedding <=> '%s'::vector
		LIMIT $1
	`, vectorLiteral(embedding), vectorLiteral(embedding))

	rows, err := r.db.Query(ctx, sql, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanSearchResults(rows)
}

// Upsert indexes (or re-indexes) a piece of content for an entity.
func (r *searchRepository) Upsert(ctx context.Context, entityType, entityID, content string, embedding []float32) error {
	var embeddingArg interface{}
	if embedding != nil {
		embeddingArg = vectorLiteral(embedding)
	}

	query := `
		INSERT INTO search_index (id, entity_type, entity_id, content, embedding)
		VALUES (gen_random_uuid(), $1, $2, $3, $4)
	`
	_, err := r.db.Exec(ctx, query, entityType, entityID, content, embeddingArg)
	return err
}

func scanSearchResults(rows interface {
	Next() bool
	Scan(dest ...interface{}) error
}) ([]*SearchResult, error) {
	var results []*SearchResult
	for rows.Next() {
		var result SearchResult
		if err := rows.Scan(&result.ID, &result.EntityType, &result.EntityID, &result.Content, &result.Score); err != nil {
			return nil, err
		}
		results = append(results, &result)
	}
	return results, nil
}
