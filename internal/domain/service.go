package domain

import (
	"log/slog"

	"github.com/robstave/gorag/internal/adapters/repositories"
	"github.com/robstave/gorag/internal/domain/types"
)

type Service struct {
	logger *slog.Logger
	repo   repositories.Repository
}

type Domain interface {
	GetdocumentByID(documentID string) (*types.Document, error)
	GetAlldocuments() ([]types.Document, error)
	Createdocument(document types.Document) (*types.Document, error)
	Updatedocument(document types.Document) (*types.Document, error)
	Deletedocument(documentID string) error
	Seeddocument() error
}

// internal/domain/service.go
// ...
func NewService(logger *slog.Logger, repo repositories.Repository) Domain {
	service := &Service{
		logger: logger,
		repo:   repo,
	}

	// Seed the initial user.   This is called on every startup, but will only create the user if it doesn't already exist
	// To reset the app, just delete the database file  ( assuming you're using the default sqlite3 database )
	if err := service.Seeddocument(); err != nil {
		logger.Error("Failed to seed initial user", "error", err)
	}

	return service
}
