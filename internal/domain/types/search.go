package types

// SearchQuery represents a search request
type SearchQuery struct {
	Query string `json:"query"`
	Limit int    `json:"limit,omitempty"`
}

// SearchResult represents a single document result with its similarity score
type SearchResult struct {
	ID       string   `json:"id"`
	Score    float64  `json:"score"`
	Document Document `json:"document"`
}

// SearchResponse represents the response to a search query
type SearchResponse struct {
	Results []SearchResult `json:"results"`
	Query   string         `json:"query"`
}
