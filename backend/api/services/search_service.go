package services

import (
	"context"
	"sort"

	"github.com/ieraHQ/Vakalat/backend/api/repositories"
)

// SearchService defines the interface for hybrid (keyword + semantic) search.
type SearchService interface {
	// Search runs a keyword search, and additionally blends in vector
	// similarity results when an embedding is supplied.
	Search(ctx context.Context, query string, embedding []float32, limit int) ([]*repositories.SearchResult, error)
	// IndexContent upserts searchable content for an entity.
	IndexContent(ctx context.Context, entityType, entityID, content string, embedding []float32) error
}

// searchService implements SearchService.
type searchService struct {
	searchRepo repositories.SearchRepository
}

// NewSearchService creates a new SearchService.
func NewSearchService(searchRepo repositories.SearchRepository) SearchService {
	return &searchService{searchRepo: searchRepo}
}

// Search performs a hybrid search: PostgreSQL full-text search over keywords,
// merged with pgvector semantic similarity when an embedding is provided.
// Results are deduplicated by ID and ranked by combined score.
func (s *searchService) Search(ctx context.Context, query string, embedding []float32, limit int) ([]*repositories.SearchResult, error) {
	keywordResults, err := s.searchRepo.SearchByKeyword(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	if len(embedding) == 0 {
		return keywordResults, nil
	}

	similarityResults, err := s.searchRepo.SearchBySimilarity(ctx, embedding, limit)
	if err != nil {
		return nil, err
	}

	merged := make(map[string]*repositories.SearchResult, len(keywordResults)+len(similarityResults))
	for _, r := range keywordResults {
		merged[r.ID] = r
	}
	for _, r := range similarityResults {
		if existing, ok := merged[r.ID]; ok {
			// Blend: average the keyword rank and vector similarity scores.
			existing.Score = (existing.Score + r.Score) / 2
			continue
		}
		merged[r.ID] = r
	}

	results := make([]*repositories.SearchResult, 0, len(merged))
	for _, r := range merged {
		results = append(results, r)
	}
	sort.Slice(results, func(i, j int) bool { return results[i].Score > results[j].Score })

	if len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

// IndexContent upserts searchable content for an entity.
func (s *searchService) IndexContent(ctx context.Context, entityType, entityID, content string, embedding []float32) error {
	return s.searchRepo.Upsert(ctx, entityType, entityID, content, embedding)
}
