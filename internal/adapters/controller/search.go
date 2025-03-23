package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/robstave/gorag/internal/domain/types"
)

// Search handles document search requests
// @Summary Search for documents
// @Description Search for documents using semantic similarity
// @Tags search
// @Accept json
// @Produce json
// @Param query query string true "Search query"
// @Param limit query int false "Maximum number of results to return (default 5)"
// @Success 200 {object} types.SearchResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /search [get]
func (c *Controller) Search(ctx echo.Context) error {
	query := ctx.QueryParam("query")
	if query == "" {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"message": "Query parameter is required"})
	}

	// Default limit is 5
	limit := 5
	if limitStr := ctx.QueryParam("limit"); limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid limit parameter"})
		}
		if parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Create search query
	searchQuery := types.SearchQuery{
		Query: query,
		Limit: limit,
	}

	c.logger.Info("Searching documents", "query", query, "limit", limit)

	// Call the service to search documents
	results, err := c.service.SearchDocuments(searchQuery)
	if err != nil {
		c.logger.Error("Failed to search documents", "error", err)
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to search documents"})
	}

	response := types.SearchResponse{
		Results: results,
		Query:   query,
	}

	return ctx.JSON(http.StatusOK, response)
}
