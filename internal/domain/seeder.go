// internal/domain/seeder.go
package domain

import (
	"github.com/google/uuid"
	"github.com/robstave/gorag/internal/domain/types"
)

func (s *Service) Seeddocument() error {
	// Default document data
	defaultdocuments := []types.Document{
		{
			ID:    uuid.New().String(),
			Name:  "Sample document 1",
			Value: "This is a sample document to demonstrate functionality",
		},
		{
			ID:    uuid.New().String(),
			Name:  "Sample document 2",
			Value: "Another document example for testing",
		},
	}

	// Check if we already have documents
	existingdocuments, err := s.repo.GetAlldocuments()
	if err != nil {
		s.logger.Error("Failed to check existing documents", "error", err)
		return err
	}

	// If we already have documents, don't seed
	if len(existingdocuments) > 0 {
		s.logger.Info("Database already has documents, skipping seed")
		return nil
	}

	// Create the sample documents
	for _, document := range defaultdocuments {
		if err := s.repo.Createdocument(document); err != nil {
			s.logger.Error("Failed to seed document", "name", document.Name, "error", err)
			return err
		}
		s.logger.Info("Seeded document successfully", "id", document.ID, "name", document.Name)
	}

	s.logger.Info("Successfully seeded initial documents")
	return nil
}
