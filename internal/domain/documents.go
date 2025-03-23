// internal/domain/documents.go
package domain

import (
	"errors"

	"github.com/google/uuid"
	"github.com/robstave/gorag/internal/domain/types"
)

func (s *Service) GetdocumentByID(documentID string) (*types.Document, error) {
	s.logger.Info("Retrieving document by ID", "documentID", documentID)

	document, err := s.repo.GetdocumentById(documentID)
	if err != nil {
		s.logger.Error("Error retrieving document", "error", err)
		return nil, err
	}

	if document == nil {
		s.logger.Warn("document not found", "documentID", documentID)
		return nil, errors.New("document not found")
	}

	return document, nil
}

func (s *Service) GetAlldocuments() ([]types.Document, error) {
	s.logger.Info("Retrieving all documents")

	documents, err := s.repo.GetAlldocuments()
	if err != nil {
		s.logger.Error("Error retrieving all documents", "error", err)
		return nil, err
	}

	return documents, nil
}

func (s *Service) Createdocument(document types.Document) (*types.Document, error) {
	s.logger.Info("Creating new document", "name", document.Name)

	// Generate UUID if not provided
	if document.ID == "" {
		document.ID = uuid.New().String()
	}

	if err := s.repo.Createdocument(document); err != nil {
		s.logger.Error("Failed to create document", "error", err)
		return nil, err
	}

	return &document, nil
}

func (s *Service) Updatedocument(document types.Document) (*types.Document, error) {
	s.logger.Info("Updating document", "id", document.ID)

	// Check if document exists
	existingdocument, err := s.repo.GetdocumentById(document.ID)
	if err != nil {
		s.logger.Error("Error checking document existence", "error", err)
		return nil, err
	}

	if existingdocument == nil {
		s.logger.Warn("document not found for update", "id", document.ID)
		return nil, errors.New("document not found")
	}

	if err := s.repo.Updatedocument(document); err != nil {
		s.logger.Error("Failed to update document", "error", err)
		return nil, err
	}

	return &document, nil
}

func (s *Service) Deletedocument(documentID string) error {
	s.logger.Info("Deleting document", "id", documentID)

	// Check if document exists
	existingdocument, err := s.repo.GetdocumentById(documentID)
	if err != nil {
		s.logger.Error("Error checking document existence", "error", err)
		return err
	}

	if existingdocument == nil {
		s.logger.Warn("document not found for deletion", "id", documentID)
		return errors.New("document not found")
	}

	if err := s.repo.Deletedocument(documentID); err != nil {
		s.logger.Error("Failed to delete document", "error", err)
		return err
	}

	return nil
}
