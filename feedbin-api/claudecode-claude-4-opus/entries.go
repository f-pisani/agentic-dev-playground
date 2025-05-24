package feedbin

import (
	"fmt"
	"strconv"
	"strings"
)

// GetEntries retrieves entries with optional filters
func (c *Client) GetEntries(options ...QueryOption) ([]Entry, error) {
	var entries []Entry
	query := buildQuery(options...)
	err := c.get("/entries.json", query, &entries)
	return entries, err
}

// GetEntry retrieves a single entry by ID
func (c *Client) GetEntry(id int) (*Entry, error) {
	var entry Entry
	path := fmt.Sprintf("/entries/%d.json", id)
	err := c.get(path, nil, &entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

// GetEntriesByIDs retrieves multiple entries by their IDs (max 100)
func (c *Client) GetEntriesByIDs(ids []int) ([]Entry, error) {
	if len(ids) == 0 {
		return []Entry{}, nil
	}
	if len(ids) > 100 {
		return nil, fmt.Errorf("maximum 100 entry IDs allowed per request, got %d", len(ids))
	}

	var entries []Entry
	query := buildQuery(WithIDs(ids))
	err := c.get("/entries.json", query, &entries)
	return entries, err
}

// GetUnreadEntries retrieves all unread entry IDs
func (c *Client) GetUnreadEntries() ([]int, error) {
	var ids []int
	err := c.get("/unread_entries.json", nil, &ids)
	return ids, err
}

// MarkEntriesRead marks the specified entries as read
func (c *Client) MarkEntriesRead(entryIDs []int) error {
	if len(entryIDs) == 0 {
		return nil
	}
	if len(entryIDs) > 1000 {
		return fmt.Errorf("maximum 1000 entry IDs allowed per request, got %d", len(entryIDs))
	}

	body := struct {
		UnreadEntries []int `json:"unread_entries"`
	}{
		UnreadEntries: entryIDs,
	}
	return c.delete("/unread_entries.json", body)
}

// MarkEntriesUnread marks the specified entries as unread
func (c *Client) MarkEntriesUnread(entryIDs []int) error {
	if len(entryIDs) == 0 {
		return nil
	}
	if len(entryIDs) > 1000 {
		return fmt.Errorf("maximum 1000 entry IDs allowed per request, got %d", len(entryIDs))
	}

	body := struct {
		UnreadEntries []int `json:"unread_entries"`
	}{
		UnreadEntries: entryIDs,
	}
	return c.post("/unread_entries.json", body, nil)
}

// GetStarredEntries retrieves all starred entry IDs
func (c *Client) GetStarredEntries() ([]int, error) {
	var ids []int
	err := c.get("/starred_entries.json", nil, &ids)
	return ids, err
}

// StarEntries marks the specified entries as starred
func (c *Client) StarEntries(entryIDs []int) error {
	if len(entryIDs) == 0 {
		return nil
	}
	if len(entryIDs) > 1000 {
		return fmt.Errorf("maximum 1000 entry IDs allowed per request, got %d", len(entryIDs))
	}

	body := struct {
		StarredEntries []int `json:"starred_entries"`
	}{
		StarredEntries: entryIDs,
	}
	return c.post("/starred_entries.json", body, nil)
}

// UnstarEntries removes the star from the specified entries
func (c *Client) UnstarEntries(entryIDs []int) error {
	if len(entryIDs) == 0 {
		return nil
	}
	if len(entryIDs) > 1000 {
		return fmt.Errorf("maximum 1000 entry IDs allowed per request, got %d", len(entryIDs))
	}

	body := struct {
		StarredEntries []int `json:"starred_entries"`
	}{
		StarredEntries: entryIDs,
	}
	return c.delete("/starred_entries.json", body)
}

// CreateEntry creates a new entry from a URL
func (c *Client) CreateEntry(url string, title string) (*Entry, error) {
	body := struct {
		URL   string `json:"url"`
		Title string `json:"title,omitempty"`
	}{
		URL:   url,
		Title: title,
	}

	var entry Entry
	err := c.post("/entries.json", body, &entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

// GetRecentlyReadEntries retrieves recently read entry IDs
func (c *Client) GetRecentlyReadEntries() ([]int, error) {
	var ids []int
	err := c.get("/recently_read_entries.json", nil, &ids)
	return ids, err
}

// AddToRecentlyRead adds entries to the recently read list
func (c *Client) AddToRecentlyRead(entryIDs []int) error {
	if len(entryIDs) == 0 {
		return nil
	}
	if len(entryIDs) > 1000 {
		return fmt.Errorf("maximum 1000 entry IDs allowed per request, got %d", len(entryIDs))
	}

	body := struct {
		RecentlyReadEntries []int `json:"recently_read_entries"`
	}{
		RecentlyReadEntries: entryIDs,
	}
	return c.post("/recently_read_entries.json", body, nil)
}

// GetUpdatedEntries retrieves entries that have been updated
func (c *Client) GetUpdatedEntries(options ...QueryOption) ([]UpdatedEntry, error) {
	var entries []UpdatedEntry
	query := buildQuery(options...)
	err := c.get("/updated_entries.json", query, &entries)
	return entries, err
}

// MarkAsRead marks all entries as read up to the specified timestamp
func (c *Client) MarkAsRead(timestamp string) error {
	body := struct {
		CreatedEntriesTimestamp string `json:"created_entries_timestamp"`
	}{
		CreatedEntriesTimestamp: timestamp,
	}
	return c.post("/unread_entries/delete_all.json", body, nil)
}

// MarkEntriesReadAlt is an alternative POST endpoint for marking entries as read
// Use this for clients that don't support DELETE with body
func (c *Client) MarkEntriesReadAlt(entryIDs []int) error {
	if len(entryIDs) == 0 {
		return nil
	}
	if len(entryIDs) > 1000 {
		return fmt.Errorf("maximum 1000 entry IDs allowed per request, got %d", len(entryIDs))
	}

	body := struct {
		UnreadEntries []int `json:"unread_entries"`
	}{
		UnreadEntries: entryIDs,
	}
	return c.post("/unread_entries/delete.json", body, nil)
}

// UnstarEntriesAlt is an alternative POST endpoint for unstarring entries
// Use this for clients that don't support DELETE with body
func (c *Client) UnstarEntriesAlt(entryIDs []int) error {
	if len(entryIDs) == 0 {
		return nil
	}
	if len(entryIDs) > 1000 {
		return fmt.Errorf("maximum 1000 entry IDs allowed per request, got %d", len(entryIDs))
	}

	body := struct {
		StarredEntries []int `json:"starred_entries"`
	}{
		StarredEntries: entryIDs,
	}
	return c.post("/starred_entries/delete.json", body, nil)
}

// MarkFeedAsRead marks all entries in a feed as read
func (c *Client) MarkFeedAsRead(feedID int) error {
	path := fmt.Sprintf("/feeds/%d/entries.json", feedID)
	body := struct {
		Read bool `json:"read"`
	}{
		Read: true,
	}
	return c.patch(path, body, nil)
}

// GetEntryIDs is a helper to convert string IDs to int IDs
func parseEntryIDs(idStr string) ([]int, error) {
	if idStr == "" {
		return []int{}, nil
	}

	parts := strings.Split(idStr, ",")
	ids := make([]int, 0, len(parts))

	for _, part := range parts {
		id, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			return nil, fmt.Errorf("invalid entry ID: %s", part)
		}
		ids = append(ids, id)
	}

	return ids, nil
}
