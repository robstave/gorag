package vectorstore

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/robstave/gorag/internal/domain/types"
)

// ChromaClient implements the VectorStore interface for Chroma DB
type ChromaClient struct {
	baseURL      string
	collectionID string
	client       *http.Client
	logger       *slog.Logger
}

// NewChromaClient creates a new Chroma client
func NewChromaClient(baseURL string, collectionName string, logger *slog.Logger) (*ChromaClient, error) {
	client := &http.Client{
		Timeout: time.Second * 30,
	}

	chromaClient := &ChromaClient{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		client:  client,
		logger:  logger,
	}

	// Ensure the collection exists
	collID, err := chromaClient.getOrCreateCollection(collectionName)
	if err != nil {
		return nil, err
	}
	chromaClient.collectionID = collID

	return chromaClient, nil
}

// getOrCreateCollection gets or creates a collection in Chroma
func (c *ChromaClient) getOrCreateCollection(name string) (string, error) {
	// First check if collection exists
	collections, err := c.listCollections()
	if err != nil {
		return "", err
	}

	for _, col := range collections {
		if col.Name == name {
			return col.ID, nil
		}
	}

	// If not, create it
	return c.createCollection(name)
}

type collection struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type listCollectionsResponse struct {
	Collections []collection `json:"collections"`
}

// listCollections lists all collections in Chroma
func (c *ChromaClient) listCollections() ([]collection, error) {
	url := fmt.Sprintf("%s/api/v1/collections", c.baseURL)

	resp, err := c.client.Get(url)
	if err != nil {
		c.logger.Error("Failed to list collections", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		c.logger.Error("Chroma API error", "status", resp.Status, "body", string(body))
		return nil, fmt.Errorf("failed to list collections: %s", resp.Status)
	}

	var result listCollectionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		c.logger.Error("Failed to decode collections response", "error", err)
		return nil, err
	}

	return result.Collections, nil
}

type createCollectionRequest struct {
	Name        string                 `json:"name"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	GetOrCreate bool                   `json:"get_or_create"`
}

type createCollectionResponse struct {
	ID string `json:"id"`
}

// createCollection creates a new collection in Chroma
func (c *ChromaClient) createCollection(name string) (string, error) {
	url := fmt.Sprintf("%s/api/v1/collections", c.baseURL)

	reqBody := createCollectionRequest{
		Name:        name,
		GetOrCreate: true,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		c.logger.Error("Failed to marshal create collection request", "error", err)
		return "", err
	}

	resp, err := c.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		c.logger.Error("Failed to create collection", "error", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		c.logger.Error("Chroma API error", "status", resp.Status, "body", string(body))
		return "", fmt.Errorf("failed to create collection: %s", resp.Status)
	}

	var result createCollectionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		c.logger.Error("Failed to decode create collection response", "error", err)
		return "", err
	}

	return result.ID, nil
}

// AddDocument adds a document to Chroma
func (c *ChromaClient) AddDocument(doc types.Document, embedding []float32) error {
	if c.collectionID == "" {
		return errors.New("collection ID not set")
	}

	url := fmt.Sprintf("%s/api/v1/collections/%s/add", c.baseURL, c.collectionID)

	// Add document to Chroma
	reqBody := map[string]interface{}{
		"ids":        []string{doc.ID},
		"embeddings": [][]float32{embedding},
		"metadatas": []map[string]interface{}{
			{
				"name":       doc.Name,
				"created_at": doc.CreatedAt.Format(time.RFC3339),
			},
		},
		"documents": []string{doc.Value},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		c.logger.Error("Failed to marshal add document request", "error", err)
		return err
	}

	resp, err := c.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		c.logger.Error("Failed to add document to Chroma", "error", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		c.logger.Error("Chroma API error", "status", resp.Status, "body", string(body))
		return fmt.Errorf("failed to add document: %s", resp.Status)
	}

	return nil
}

// QueryDocuments queries documents from Chroma
func (c *ChromaClient) QueryDocuments(query string, embedding []float32, limit int) ([]types.SearchResult, error) {
	if c.collectionID == "" {
		return nil, errors.New("collection ID not set")
	}

	url := fmt.Sprintf("%s/api/v1/collections/%s/query", c.baseURL, c.collectionID)

	reqBody := map[string]interface{}{
		"query_embeddings": [][]float32{embedding},
		"n_results":        limit,
		"include":          []string{"metadatas", "documents", "distances"},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		c.logger.Error("Failed to marshal query request", "error", err)
		return nil, err
	}

	resp, err := c.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		c.logger.Error("Failed to query Chroma", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		c.logger.Error("Chroma API error", "status", resp.Status, "body", string(body))
		return nil, fmt.Errorf("failed to query documents: %s", resp.Status)
	}

	// Parse response
	var queryResp struct {
		IDs       [][]string                 `json:"ids"`
		Documents [][]string                 `json:"documents"`
		Metadatas [][]map[string]interface{} `json:"metadatas"`
		Distances [][]float64                `json:"distances"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&queryResp); err != nil {
		c.logger.Error("Failed to decode query response", "error", err)
		return nil, err
	}

	// Ensure we have results
	if len(queryResp.IDs) == 0 || len(queryResp.IDs[0]) == 0 {
		return []types.SearchResult{}, nil
	}

	// Map to search results
	results := make([]types.SearchResult, len(queryResp.IDs[0]))
	for i, id := range queryResp.IDs[0] {
		result := types.SearchResult{
			ID:    id,
			Score: queryResp.Distances[0][i],
			Document: types.Document{
				ID:    id,
				Value: queryResp.Documents[0][i],
			},
		}

		// Add metadata if available
		if i < len(queryResp.Metadatas[0]) {
			metadata := queryResp.Metadatas[0][i]
			if name, ok := metadata["name"].(string); ok {
				result.Document.Name = name
			}
		}

		results[i] = result
	}

	return results, nil
}

// DeleteDocument deletes a document from Chroma
func (c *ChromaClient) DeleteDocument(id string) error {
	if c.collectionID == "" {
		return errors.New("collection ID not set")
	}

	url := fmt.Sprintf("%s/api/v1/collections/%s/delete", c.baseURL, c.collectionID)

	reqBody := map[string]interface{}{
		"ids": []string{id},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		c.logger.Error("Failed to marshal delete request", "error", err)
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		c.logger.Error("Failed to create delete request", "error", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error("Failed to delete from Chroma", "error", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		c.logger.Error("Chroma API error", "status", resp.Status, "body", string(body))
		return fmt.Errorf("failed to delete document: %s", resp.Status)
	}

	return nil
}
