package ops

import "github.com/jrogala/jellyfin-cli/client"

// SearchHintEntry represents a search result.
type SearchHintEntry struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Type           string `json:"type"`
	ProductionYear int    `json:"production_year,omitempty"`
}

// SearchResult holds the result of a search query.
type SearchResult struct {
	Hints            []SearchHintEntry `json:"hints"`
	TotalRecordCount int               `json:"total_record_count"`
}

// Search searches for media matching the given query.
func Search(c *client.Client, query string, limit int) (*SearchResult, error) {
	result, err := c.Search(query, limit)
	if err != nil {
		return nil, err
	}
	hints := make([]SearchHintEntry, 0, len(result.SearchHints))
	for _, h := range result.SearchHints {
		hints = append(hints, SearchHintEntry{
			ID:             h.ID,
			Name:           h.Name,
			Type:           h.Type,
			ProductionYear: h.ProductionYear,
		})
	}
	return &SearchResult{
		Hints:            hints,
		TotalRecordCount: result.TotalRecordCount,
	}, nil
}
