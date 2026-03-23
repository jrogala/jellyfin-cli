// Package ops contains business logic for Jellyfin operations.
// Functions return Go structs and errors -- zero I/O, zero formatting.
package ops

import "github.com/jrogala/jellyfin-cli/client"

// LibraryEntry represents a library in listing results.
type LibraryEntry struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	CollectionType string `json:"collection_type"`
}

// ListLibraries returns all libraries/views visible to the user.
func ListLibraries(c *client.Client) ([]LibraryEntry, error) {
	libs, err := c.GetLibraries()
	if err != nil {
		return nil, err
	}
	var entries []LibraryEntry
	for _, l := range libs {
		entries = append(entries, LibraryEntry{
			ID:             l.ID,
			Name:           l.Name,
			CollectionType: l.CollectionType,
		})
	}
	return entries, nil
}

// ScanLibrary triggers a full library scan.
func ScanLibrary(c *client.Client) error {
	return c.ScanLibrary()
}
