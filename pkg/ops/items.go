package ops

import (
	"fmt"
	"strconv"

	"github.com/jrogala/jellyfin-cli/client"
)

// ItemEntry represents a media item in listing/search results.
type ItemEntry struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Type            string            `json:"type"`
	Path            string            `json:"path,omitempty"`
	ProductionYear  int               `json:"production_year,omitempty"`
	CommunityRating float64           `json:"community_rating,omitempty"`
	OfficialRating  string            `json:"official_rating,omitempty"`
	Overview        string            `json:"overview,omitempty"`
	Runtime         string            `json:"runtime,omitempty"`
	ProviderIds     map[string]string `json:"provider_ids,omitempty"`
}

// ItemsResult holds the result of an items query.
type ItemsResult struct {
	Items            []ItemEntry `json:"items"`
	TotalRecordCount int         `json:"total_record_count"`
}

// ListItemsOptions configures item listing.
type ListItemsOptions struct {
	LibraryID string
	ItemType  string
	Limit     int
}

// ListItems returns items matching the given options.
func ListItems(c *client.Client, opts ListItemsOptions) (*ItemsResult, error) {
	result, err := c.GetItems(opts.LibraryID, opts.ItemType, "SortName", opts.Limit)
	if err != nil {
		return nil, err
	}
	entries := make([]ItemEntry, 0, len(result.Items))
	for _, item := range result.Items {
		entries = append(entries, itemToEntry(item))
	}
	return &ItemsResult{
		Items:            entries,
		TotalRecordCount: result.TotalRecordCount,
	}, nil
}

// ListMovies returns all movies sorted by name.
func ListMovies(c *client.Client) (*ItemsResult, error) {
	return ListItems(c, ListItemsOptions{ItemType: "Movie"})
}

// GetItemInfo returns detailed information about a single item.
func GetItemInfo(c *client.Client, itemID string) (*ItemEntry, error) {
	item, err := c.GetItem(itemID)
	if err != nil {
		return nil, err
	}
	entry := itemToEntry(*item)
	return &entry, nil
}

func itemToEntry(item client.Item) ItemEntry {
	return ItemEntry{
		ID:              item.ID,
		Name:            item.Name,
		Type:            item.Type,
		Path:            item.Path,
		ProductionYear:  item.ProductionYear,
		CommunityRating: item.CommunityRating,
		OfficialRating:  item.OfficialRating,
		Overview:        item.Overview,
		Runtime:         client.FormatRuntime(item.RunTimeTicks),
		ProviderIds:     item.ProviderIds,
	}
}

// FormatYear formats a year for display, returning empty string for zero.
func FormatYear(year int) string {
	if year > 0 {
		return strconv.Itoa(year)
	}
	return ""
}

// FormatRating formats a rating for display, returning empty string for zero.
func FormatRating(rating float64) string {
	if rating > 0 {
		return fmt.Sprintf("%.1f", rating)
	}
	return ""
}
