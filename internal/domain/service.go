package domain

import (
	"log/slog"

	"github.com/robstave/gorag/internal/adapters/repositories"
	"github.com/robstave/gorag/internal/adapters/repositories/vectorstore"
	"github.com/robstave/gorag/internal/domain/embedding"
	"github.com/robstave/gorag/internal/domain/types"
)

type Service struct {
	logger       *slog.Logger
	repo         repositories.Repository
	vectorStore  vectorstore.VectorStore
	embedService embedding.EmbeddingService
}

type Domain interface {
	GetdocumentByID(documentID string) (*types.Document, error)
	GetAlldocuments() ([]types.Document, error)
	Createdocument(document types.Document) (*types.Document, error)
	Updatedocument(document types.Document) (*types.Document, error)
	Deletedocument(documentID string) error
	Seeddocument() error
	SearchDocuments(query types.SearchQuery) ([]types.SearchResult, error)
}

// NewService creates a new instance of the domain service
func NewService(logger *slog.Logger, repo repositories.Repository, vectorStore vectorstore.VectorStore, embedService embedding.EmbeddingService) Domain {
	service := &Service{
		logger:       logger,
		repo:         repo,
		vectorStore:  vectorStore,
		embedService: embedService,
	}

	// Seed the initial documents. This is called on every startup but will only create documents if they don't already exist
	// To reset the app, just delete the database file (assuming you're using the default sqlite3 database)
	if err := service.Seeddocument(); err != nil {
		logger.Error("Failed to seed initial documents", "error", err)
	}

	return service
}
