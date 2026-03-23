package ops

import "github.com/jrogala/jellyfin-cli/client"

// UpdateItem updates metadata fields on an item.
func UpdateItem(c *client.Client, itemID string, updates map[string]any) error {
	return c.UpdateItem(itemID, updates)
}

// IdentifyItem matches an item with metadata providers.
func IdentifyItem(c *client.Client, itemID, name string, year int, providerIDs map[string]string) error {
	return c.IdentifyItem(itemID, name, year, providerIDs)
}

// RefreshMetadata triggers a metadata refresh for an item.
func RefreshMetadata(c *client.Client, itemID string) error {
	return c.RefreshMetadata(itemID)
}
