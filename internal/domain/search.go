package domain

import (
	"github.com/robstave/gorag/internal/domain/types"
)

// SearchDocuments searches for documents using vector similarity
func (s *Service) SearchDocuments(query types.SearchQuery) ([]types.SearchResult, error) {
	s.logger.Info("Searching documents", "query", query.Query)

	// Generate embedding for the query
	embedding, err := s.embedService.CreateEmbedding(query.Query)
	if err != nil {
		s.logger.Error("Failed to create embedding", "error", err)
		return nil, err
	}

	// Use a default limit if not specified
	limit := query.Limit
	if limit <= 0 {
		limit = 5
	}

	// Query the vector store
	results, err := s.vectorStore.QueryDocuments(query.Query, embedding, limit)
	if err != nil {
		s.logger.Error("Failed to query vector store", "error", err)
		return nil, err
	}

	// Ensure we have complete document information
	for i, result := range results {
		// If the document content is empty, fetch it from the SQL database
		if result.Document.Value == "" {
			doc, err := s.repo.GetdocumentById(result.ID)
			if err != nil {
				s.logger.Error("Failed to get document from DB", "id", result.ID, "error", err)
				continue
			}
			if doc != nil {
				results[i].Document = *doc
			}
		}
	}

	return results, nil
}
