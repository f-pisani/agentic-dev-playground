package feedbin

import (
	"net/http"
)

// GetStarredEntryIDs retrieves the IDs of all starred entries.
func (c *Client) GetStarredEntryIDs() ([]int, error) {
	// Make the request
	resp, err := c.doRequest(http.MethodGet, "/starred_entries.json", nil, nil)
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

// StarEntries marks the specified entries as starred.
func (c *Client) StarEntries(entryIDs []int) error {
	// Create the request body
	reqBody := EntryIDs{
		IDs: entryIDs,
	}

	// Make the request
	resp, err := c.doRequest(http.MethodPost, "/starred_entries.json", reqBody, nil)
	if err != nil {
		return err
	}

	// Check for success (200 OK)
	if resp.StatusCode != http.StatusOK {
		return parseResponse(resp, nil)
	}

	return nil
}

// UnstarEntries removes the star from the specified entries.
func (c *Client) UnstarEntries(entryIDs []int) error {
	// Create the request body
	reqBody := EntryIDs{
		IDs: entryIDs,
	}

	// Make the request
	resp, err := c.doRequest(http.MethodPost, "/starred_entries/delete.json", reqBody, nil)
	if err != nil {
		return err
	}

	// Check for success (200 OK)
	if resp.StatusCode != http.StatusOK {
		return parseResponse(resp, nil)
	}

	return nil
}
