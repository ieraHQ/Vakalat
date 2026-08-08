// Command backfill-embeddings indexes and embeds every existing client,
// matter, and OCR-completed document into the search index. The live
// pipeline (services/client_service.go, services/matter_service.go,
// worker/ocr_worker.go) indexes new/updated content automatically — this
// command is only needed once, for rows that existed before that pipeline
// was wired up.
//
// Usage: go run ./cmd/backfill-embeddings
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ieraHQ/Vakalat/backend/api/ai"
	"github.com/ieraHQ/Vakalat/backend/api/config"
	"github.com/ieraHQ/Vakalat/backend/api/database"
	"github.com/ieraHQ/Vakalat/backend/api/repositories"
	"github.com/ieraHQ/Vakalat/backend/api/services"
)

const pageSize = 100

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := database.InitDB(cfg); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer database.CloseDB()

	ctx := context.Background()

	clientRepo := repositories.NewClientRepository(database.DB)
	matterRepo := repositories.NewMatterRepository(database.DB)
	documentRepo := repositories.NewDocumentRepository(database.DB)
	searchRepo := repositories.NewSearchRepository(database.DB)

	searchService := services.NewSearchService(searchRepo)
	embedder := ai.NewLLMClient(cfg.AI.BaseURL, cfg.AI.Model, cfg.AI.EmbeddingModel)

	indexed, failed := 0, 0

	for offset := 0; ; offset += pageSize {
		clients, err := clientRepo.List(ctx, pageSize, offset)
		if err != nil {
			log.Fatalf("failed to list clients: %v", err)
		}
		if len(clients) == 0 {
			break
		}
		for _, c := range clients {
			notes := ""
			if c.Notes != nil {
				notes = *c.Notes
			}
			content := fmt.Sprintf("%s %s %s %s", c.Name, c.Email, c.Phone, notes)
			if indexContent(ctx, searchService, embedder, "client", c.ID, content) {
				indexed++
			} else {
				failed++
			}
		}
	}

	for offset := 0; ; offset += pageSize {
		matters, err := matterRepo.ListAll(ctx, pageSize, offset)
		if err != nil {
			log.Fatalf("failed to list matters: %v", err)
		}
		if len(matters) == 0 {
			break
		}
		for _, m := range matters {
			content := fmt.Sprintf("%s %s %s", m.Title, m.Description, m.CaseNumber)
			if indexContent(ctx, searchService, embedder, "matter", m.ID, content) {
				indexed++
			} else {
				failed++
			}
		}
	}

	// Only OCR-completed documents have text worth indexing.
	documents, err := documentRepo.ListByOCRStatus(ctx, "completed")
	if err != nil {
		log.Fatalf("failed to list completed documents: %v", err)
	}
	for _, d := range documents {
		content := d.Name + "\n" + d.OCRText
		if indexContent(ctx, searchService, embedder, "document", d.ID, content) {
			indexed++
		} else {
			failed++
		}
	}

	log.Printf("Backfill complete: %d indexed, %d failed", indexed, failed)
}

func indexContent(ctx context.Context, searchService services.SearchService, embedder ai.LLMClient, entityType, entityID, content string) bool {
	embedding, err := embedder.Embed(ctx, content)
	if err != nil {
		log.Printf("warning: failed to embed %s %s: %v (indexing without embedding)", entityType, entityID, err)
		embedding = nil
	}

	if err := searchService.IndexContent(ctx, entityType, entityID, content, embedding); err != nil {
		log.Printf("error: failed to index %s %s: %v", entityType, entityID, err)
		return false
	}
	return true
}
