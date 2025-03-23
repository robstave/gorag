// internal/adapters/repositories/user.go
package repositories

import (
	"github.com/robstave/gorag/internal/domain/types"
	"gorm.io/gorm"
)

type Repository interface {
	GetdocumentById(id string) (*types.Document, error)
	GetAlldocuments() ([]types.Document, error)
	Createdocument(document types.Document) error
	Updatedocument(document types.Document) error
	Deletedocument(id string) error
}

type RepositorySQLite struct {
	db *gorm.DB
}

func NewRepositorySQLite(db *gorm.DB) Repository {
	return &RepositorySQLite{db: db}
}

func (r *RepositorySQLite) GetdocumentById(id string) (*types.Document, error) {
	var document types.Document
	result := r.db.First(&document, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &document, nil
}

func (r *RepositorySQLite) GetAlldocuments() ([]types.Document, error) {
	var documents []types.Document
	result := r.db.Find(&documents)
	if result.Error != nil {
		return nil, result.Error
	}
	return documents, nil
}

func (r *RepositorySQLite) Createdocument(document types.Document) error {
	return r.db.Create(&document).Error
}

func (r *RepositorySQLite) Updatedocument(document types.Document) error {
	return r.db.Save(&document).Error
}

func (r *RepositorySQLite) Deletedocument(id string) error {
	return r.db.Delete(&types.Document{}, "id = ?", id).Error
}
