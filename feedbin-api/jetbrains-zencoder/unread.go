package feedbin

import (
	"net/http"
)

// GetUnreadEntryIDs retrieves the IDs of all unread entries.
func (c *Client) GetUnreadEntryIDs() ([]int, error) {
	// Make the request
	resp, err := c.doRequest(http.MethodGet, "/unread_entries.json", nil, nil)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var ids []int
	if err := parseResponse(resp, &ids); err != nil {
		return nil, err
	}

	return ids, nil
}

// MarkEntriesAsRead marks the specified entries as read.
func (c *Client) MarkEntriesAsRead(entryIDs []int) error {
	// Create the request body
	reqBody := EntryIDs{
		IDs: entryIDs,
	}

	// Make the request
	resp, err := c.doRequest(http.MethodPost, "/unread_entries/delete.json", reqBody, nil)
	if err != nil {
		return err
	}

	// Check for success (200 OK)
	if resp.StatusCode != http.StatusOK {
		return parseResponse(resp, nil)
	}

	return nil
}

// MarkEntriesAsUnread marks the specified entries as unread.
func (c *Client) MarkEntriesAsUnread(entryIDs []int) error {
	// Create the request body
	reqBody := EntryIDs{
		IDs: entryIDs,
	}

	// Make the request
	resp, err := c.doRequest(http.MethodPost, "/unread_entries.json", reqBody, nil)
	if err != nil {
		return err
	}

	// Check for success (200 OK)
	if resp.StatusCode != http.StatusOK {
		return parseResponse(resp, nil)
	}

	return nil
}
