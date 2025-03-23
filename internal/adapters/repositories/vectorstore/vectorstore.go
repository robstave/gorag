package vectorstore

import (
	"github.com/robstave/gorag/internal/domain/types"
)

// VectorStore represents a repository for storing and querying document embeddings
type VectorStore interface {
	// AddDocument adds a document and its embedding to the vector store
	AddDocument(doc types.Document, embedding []float32) error

	// QueryDocuments finds similar documents based on the query embedding
	QueryDocuments(query string, embedding []float32, limit int) ([]types.SearchResult, error)

	// DeleteDocument removes a document from the vector store
	DeleteDocument(id string) error
}
