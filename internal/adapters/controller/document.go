// internal/adapters/controller/document.go
package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/robstave/gorag/internal/domain/types"
)

// Createdocument handles document creation
// @Summary Create a new document
// @Description Create a new document
// @Tags documents
// @Accept json
// @Produce json
// @Success 201 {object} types.Document
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /documents [post]
func (hc *Controller) Createdocument(c echo.Context) error {
	var document types.Document
	if err := c.Bind(&document); err != nil {
		hc.logger.Error("Failed to bind document data", "error", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid document data"})
	}

	// Call the service to create the document
	createddocument, err := hc.service.Createdocument(document)
	if err != nil {
		hc.logger.Error("Failed to create document", "error", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to create document"})
	}

	return c.JSON(http.StatusCreated, createddocument)
}

// Getdocument retrieves a document by ID
// @Summary Get a document by ID
// @Description Get a document by its ID
// @Tags documents
// @Produce json
// @Param id path string true "document ID"
// @Success 200 {object} types.Document
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /documents/{id} [get]
func (hc *Controller) Getdocument(c echo.Context) error {
	id := c.Param("id")

	document, err := hc.service.GetdocumentByID(id)
	if err != nil {
		hc.logger.Error("Failed to retrieve document", "id", id, "error", err)
		return c.JSON(http.StatusNotFound, echo.Map{"message": "document not found"})
	}

	return c.JSON(http.StatusOK, document)
}

// GetAlldocuments retrieves all documents
// @Summary Get all documents
// @Description Get all available documents
// @Tags documents
// @Produce json
// @Success 200 {array} types.Document
// @Failure 500 {object} map[string]string
// @Router /documents [get]
func (hc *Controller) GetAlldocuments(c echo.Context) error {
	documents, err := hc.service.GetAlldocuments()
	if err != nil {
		hc.logger.Error("Failed to retrieve documents", "error", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve documents"})
	}

	return c.JSON(http.StatusOK, documents)
}

// Updatedocument updates an existing document
// @Summary Update a document
// @Description Update an existing document
// @Tags documents
// @Accept json
// @Produce json
// @Param id path string true "document ID"
// @Success 200 {object} types.Document
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /documents/{id} [put]
func (hc *Controller) Updatedocument(c echo.Context) error {
	id := c.Param("id")

	var document types.Document
	if err := c.Bind(&document); err != nil {
		hc.logger.Error("Failed to bind document data", "error", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid document data"})
	}

	// Ensure ID in path matches body
	document.ID = id

	updateddocument, err := hc.service.Updatedocument(document)
	if err != nil {
		hc.logger.Error("Failed to update document", "id", id, "error", err)
		return c.JSON(http.StatusNotFound, echo.Map{"message": "document not found or update failed"})
	}

	return c.JSON(http.StatusOK, updateddocument)
}

// Deletedocument deletes a document by ID
// @Summary Delete a document
// @Description Delete a document by its ID
// @Tags documents
// @Param id path string true "document ID"
// @Success 204 {object} nil
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /documents/{id} [delete]
func (hc *Controller) Deletedocument(c echo.Context) error {
	id := c.Param("id")

	if err := hc.service.Deletedocument(id); err != nil {
		hc.logger.Error("Failed to delete document", "id", id, "error", err)
		return c.JSON(http.StatusNotFound, echo.Map{"message": "document not found or delete failed"})
	}

	return c.NoContent(http.StatusNoContent)
}
