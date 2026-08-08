package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/ieraHQ/Vakalat/backend/api/ai"
	"github.com/ieraHQ/Vakalat/backend/api/repositories"
)

// AIService defines the interface for AI-assisted operations: summarization,
// drafting, and research over a matter, backed by a local LLM.
type AIService interface {
	// Summarize generates and stores a summary of the given text for a matter.
	Summarize(ctx context.Context, matterID, text string) (*repositories.AISummary, error)
	// Ask runs a free-form prompt against the assistant, scoped to a matter, and
	// stores the exchange as an AI session.
	Ask(ctx context.Context, matterID, prompt string) (*repositories.AISession, error)
	// Draft generates a draft document (e.g. a legal notice or application) from
	// instructions, scoped to a matter.
	Draft(ctx context.Context, matterID, instructions string) (*repositories.AISession, error)
}

// aiService implements AIService.
type aiService struct {
	aiRepo repositories.AIRepository
	llm    ai.LLMClient
	model  string
}

// NewAIService creates a new AIService.
func NewAIService(aiRepo repositories.AIRepository, llm ai.LLMClient, model string) AIService {
	return &aiService{aiRepo: aiRepo, llm: llm, model: model}
}

const summarizeSystemPrompt = "You are a legal assistant. Summarize the following case material concisely and accurately, preserving names, dates, and case-critical facts. Do not invent information that is not present."

// Summarize generates and stores a concise summary of the given text for a matter.
func (s *aiService) Summarize(ctx context.Context, matterID, text string) (*repositories.AISummary, error) {
	response, err := s.llm.Chat(ctx, summarizeSystemPrompt, text)
	if err != nil {
		return nil, err
	}

	summary := &repositories.AISummary{
		ID:       uuid.New().String(),
		MatterID: matterID,
		Summary:  response,
	}

	if err := s.aiRepo.CreateSummary(ctx, summary); err != nil {
		return nil, err
	}

	return summary, nil
}

const askSystemPrompt = "You are a legal research and case assistant. Answer questions about the matter clearly and cite relevant facts. If you are unsure, say so instead of guessing."

// Ask runs a free-form prompt against the assistant and stores the exchange.
func (s *aiService) Ask(ctx context.Context, matterID, prompt string) (*repositories.AISession, error) {
	response, err := s.llm.Chat(ctx, askSystemPrompt, prompt)
	if err != nil {
		return nil, err
	}

	session := &repositories.AISession{
		ID:       uuid.New().String(),
		MatterID: matterID,
		Prompt:   prompt,
		Response: response,
		Model:    s.model,
	}

	if err := s.aiRepo.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

const draftSystemPrompt = "You are a legal drafting assistant. Draft the requested document in a formal legal style appropriate for Indian courts. Leave clearly marked placeholders (e.g. [DATE], [PARTY NAME]) for any facts not provided in the instructions."

// Draft generates a draft legal document from instructions and stores the exchange.
func (s *aiService) Draft(ctx context.Context, matterID, instructions string) (*repositories.AISession, error) {
	response, err := s.llm.Chat(ctx, draftSystemPrompt, instructions)
	if err != nil {
		return nil, err
	}

	session := &repositories.AISession{
		ID:       uuid.New().String(),
		MatterID: matterID,
		Prompt:   instructions,
		Response: response,
		Model:    s.model,
	}

	if err := s.aiRepo.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}
