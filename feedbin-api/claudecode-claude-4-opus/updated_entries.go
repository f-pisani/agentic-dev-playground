package feedbin

import (
	"fmt"
)

// MarkUpdatedEntriesRead marks the specified updated entries as read
func (c *Client) MarkUpdatedEntriesRead(entryIDs []int) error {
	if len(entryIDs) == 0 {
		return nil
	}
	if len(entryIDs) > 1000 {
		return fmt.Errorf("maximum 1000 entry IDs allowed per request, got %d", len(entryIDs))
	}

	body := struct {
		UpdatedEntries []int `json:"updated_entries"`
	}{
		UpdatedEntries: entryIDs,
	}
	return c.delete("/updated_entries.json", body)
}

// MarkUpdatedEntriesReadAlt marks updated entries as read using POST (alternative to DELETE)
func (c *Client) MarkUpdatedEntriesReadAlt(entryIDs []int) error {
	if len(entryIDs) == 0 {
		return nil
	}
	if len(entryIDs) > 1000 {
		return fmt.Errorf("maximum 1000 entry IDs allowed per request, got %d", len(entryIDs))
	}

	body := struct {
		UpdatedEntries []int `json:"updated_entries"`
	}{
		UpdatedEntries: entryIDs,
	}
	return c.post("/updated_entries/delete.json", body, nil)
}
