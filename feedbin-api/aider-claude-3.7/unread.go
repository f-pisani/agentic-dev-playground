package feedbin

import (
	"fmt"
	"net/http"
	"strconv"
)

// GetUnreadEntries retrieves all unread entry IDs
func (c *Client) GetUnreadEntries() ([]int64, error) {
	req, err := c.NewRequest(http.MethodGet, "/unread_entries.json", nil)
	if err != nil {
		return nil, err
	}
	
	var unreadIDs []int64
	_, err = c.Do(req, &unreadIDs)
	if err != nil {
		return nil, err
	}
	
	return unreadIDs, nil
}

// MarkEntriesAsUnread marks entries as unread
func (c *Client) MarkEntriesAsUnread(entryIDs []int64) error {
	if len(entryIDs) == 0 {
		return nil
	}
	
	unreadReq := &UnreadEntryRequest{
		UnreadEntries: entryIDs,
	}
	
	req, err := c.NewRequest(http.MethodPost, "/unread_entries.json", unreadReq)
	if err != nil {
		return err
	}
	
	_, err = c.Do(req, nil)
	return err
}

// MarkEntriesAsRead marks entries as read (removes from unread)
func (c *Client) MarkEntriesAsRead(entryIDs []int64) error {
	if len(entryIDs) == 0 {
		return nil
	}
	
	unreadReq := &UnreadEntryRequest{
		UnreadEntries: entryIDs,
	}
	
	req, err := c.NewRequest(http.MethodDelete, "/unread_entries.json", unreadReq)
	if err != nil {
		return err
	}
	
	_, err = c.Do(req, nil)
	return err
}

// GetUnreadCount returns the total number of unread entries
func (c *Client) GetUnreadCount() (int, error) {
	unreadIDs, err := c.GetUnreadEntries()
	if err != nil {
		return 0, err
	}
	
	return len(unreadIDs), nil
}

// GetUnreadEntriesByFeed returns a map of feed IDs to their unread entry counts
func (c *Client) GetUnreadEntriesByFeed() (map[int64]int, error) {
	// First get all unread entry IDs
	unreadIDs, err := c.GetUnreadEntries()
	if err != nil {
		return nil, err
	}
	
	if len(unreadIDs) == 0 {
		return make(map[int64]int), nil
	}
	
	// Then get the entries to determine which feed they belong to
	entries, err := c.GetEntriesByIDs(unreadIDs)
	if err != nil {
		return nil, err
	}
	
	// Count entries by feed
	feedCounts := make(map[int64]int)
	for _, entry := range entries {
		feedCounts[entry.FeedID]++
	}
	
	return feedCounts, nil
}
