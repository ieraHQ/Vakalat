package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// LLMClient defines the interface for talking to a local, OpenAI-compatible
// chat completion endpoint (e.g. Ollama, LM Studio). Keeping this behind an
// interface means the assistant is never tied to a single model or provider.
type LLMClient interface {
	Chat(ctx context.Context, systemPrompt, userPrompt string) (string, error)
	// Embed returns a vector embedding for text, using the OpenAI-compatible
	// /embeddings endpoint. Used to populate search_index.embedding and
	// embeddings.embedding so hybrid search has a semantic component instead
	// of falling back to keyword-only.
	Embed(ctx context.Context, text string) ([]float32, error)
}

// openAICompatibleClient implements LLMClient against any server that speaks
// the OpenAI /chat/completions and /embeddings wire formats.
type openAICompatibleClient struct {
	baseURL        string
	model          string
	embeddingModel string
	http           *http.Client
}

// NewLLMClient creates a new LLMClient pointed at baseURL (e.g.
// "http://localhost:11434/v1" for Ollama), using the given chat model and
// embedding model. embeddingModel should produce 384-dimension vectors to
// match the schema (e.g. Ollama's "all-minilm").
func NewLLMClient(baseURL, model, embeddingModel string) LLMClient {
	return &openAICompatibleClient{
		baseURL:        baseURL,
		model:          model,
		embeddingModel: embeddingModel,
		http:           &http.Client{Timeout: 60 * time.Second},
	}
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatCompletionRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

type chatCompletionResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
}

// Chat sends a system+user prompt pair to the configured local LLM and
// returns the assistant's reply text.
func (c *openAICompatibleClient) Chat(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	reqBody := chatCompletionRequest{
		Model: c.model,
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("llm request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var completion chatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&completion); err != nil {
		return "", err
	}

	if len(completion.Choices) == 0 {
		return "", fmt.Errorf("llm returned no choices")
	}

	return completion.Choices[0].Message.Content, nil
}

type embeddingRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type embeddingResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
}

// Embed requests a vector embedding for text from the configured embedding model.
func (c *openAICompatibleClient) Embed(ctx context.Context, text string) ([]float32, error) {
	reqBody := embeddingRequest{
		Model: c.embeddingModel,
		Input: text,
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/embeddings", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("embedding request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var embedding embeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&embedding); err != nil {
		return nil, err
	}

	if len(embedding.Data) == 0 {
		return nil, fmt.Errorf("embedding request returned no data")
	}

	return embedding.Data[0].Embedding, nil
}
