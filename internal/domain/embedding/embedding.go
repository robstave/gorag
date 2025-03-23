package embedding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

// EmbeddingService handles the generation of embeddings for documents
type EmbeddingService struct {
	client    *http.Client
	apiKey    string
	model     string
	endpoint  string
	dimension int
	logger    *slog.Logger
}

// NewOpenAIEmbeddingService creates a new embedding service using OpenAI
func NewOpenAIEmbeddingService(logger *slog.Logger) *EmbeddingService {
	apiKey := os.Getenv("OPENAI_API_KEY")
	model := "text-embedding-3-small"
	if m := os.Getenv("OPENAI_EMBEDDING_MODEL"); m != "" {
		model = m
	}

	dimension := 1536 // default for text-embedding-3-small

	return &EmbeddingService{
		client: &http.Client{
			Timeout: time.Second * 30,
		},
		apiKey:    apiKey,
		model:     model,
		endpoint:  "https://api.openai.com/v1/embeddings",
		dimension: dimension,
		logger:    logger,
	}
}

// CreateEmbedding generates an embedding for the given text
func (s *EmbeddingService) CreateEmbedding(text string) ([]float32, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("OpenAI API key not set")
	}

	// Clean up text - remove excessive whitespace
	text = strings.TrimSpace(text)

	// Create request to OpenAI
	reqBody := map[string]interface{}{
		"input": text,
		"model": s.model,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		s.logger.Error("Failed to marshal embedding request", "error", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", s.endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		s.logger.Error("Failed to create embedding request", "error", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))

	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Error("Failed to call OpenAI API", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		s.logger.Error("OpenAI API error", "status", resp.Status, "body", string(body))
		return nil, fmt.Errorf("OpenAI API error: %s", resp.Status)
	}

	// Parse response
	var result struct {
		Object string `json:"object"`
		Data   []struct {
			Embedding []float32 `json:"embedding"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		s.logger.Error("Failed to decode embedding response", "error", err)
		return nil, err
	}

	if len(result.Data) == 0 || len(result.Data[0].Embedding) == 0 {
		return nil, fmt.Errorf("empty embedding result")
	}

	return result.Data[0].Embedding, nil
}

// GetEmbeddingDimension returns the dimension of embeddings from this service
func (s *EmbeddingService) GetEmbeddingDimension() int {
	return s.dimension
}

// MockEmbeddingService is a simple mock implementation for testing
type MockEmbeddingService struct {
	dimension int
}

// NewMockEmbeddingService creates a new mock embedding service
func NewMockEmbeddingService() *MockEmbeddingService {
	return &MockEmbeddingService{
		dimension: 384, // A smaller dimension for testing
	}
}

// CreateEmbedding generates a mock embedding with random values
func (s *MockEmbeddingService) CreateEmbedding(text string) ([]float32, error) {
	// Create a deterministic mock embedding based on the text length
	embedding := make([]float32, s.dimension)
	for i := range embedding {
		// Use text length to seed a simple value
		embedding[i] = float32(len(text)%10) / 10.0
	}
	return embedding, nil
}

// GetEmbeddingDimension returns the dimension of mock embeddings
func (s *MockEmbeddingService) GetEmbeddingDimension() int {
	return s.dimension
}
