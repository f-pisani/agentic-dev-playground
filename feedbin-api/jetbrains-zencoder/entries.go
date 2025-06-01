package feedbin

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// GetEntries retrieves entries for the authenticated user.
func (c *Client) GetEntries(options map[string]interface{}) ([]Entry, *PaginationInfo, error) {
	// Build query parameters
	query := url.Values{}

	// Add options to query parameters
	for key, value := range options {
		switch v := value.(type) {
		case int:
			query.Set(key, strconv.Itoa(v))
		case bool:
			query.Set(key, strconv.FormatBool(v))
		case string:
			query.Set(key, v)
		case time.Time:
			query.Set(key, v.Format(time.RFC3339Nano))
		case []int:
			// Convert slice of ints to comma-separated string
			strValues := make([]string, len(v))
			for i, id := range v {
				strValues[i] = strconv.Itoa(id)
			}
			query.Set(key, strings.Join(strValues, ","))
		}
	}

	// Make the request
	resp, err := c.doRequest(http.MethodGet, "/entries.json", nil, query)
	if err != nil {
		return nil, nil, err
	}

	// Parse pagination info
	paginationInfo := GetPaginationInfo(resp)

	// Parse the response
	var entries []Entry
	if err := parseResponse(resp, &entries); err != nil {
		return nil, nil, err
	}

	return entries, &paginationInfo, nil
}

// GetFeedEntries retrieves entries for a specific feed.
func (c *Client) GetFeedEntries(feedID int, options map[string]interface{}) ([]Entry, *PaginationInfo, error) {
	// Build query parameters
	query := url.Values{}

	// Add options to query parameters
	for key, value := range options {
		switch v := value.(type) {
		case int:
			query.Set(key, strconv.Itoa(v))
		case bool:
			query.Set(key, strconv.FormatBool(v))
		case string:
			query.Set(key, v)
		case time.Time:
			query.Set(key, v.Format(time.RFC3339Nano))
		}
	}

	// Make the request
	resp, err := c.doRequest(http.MethodGet, fmt.Sprintf("/feeds/%d/entries.json", feedID), nil, query)
	if err != nil {
		return nil, nil, err
	}

	// Parse pagination info
	paginationInfo := GetPaginationInfo(resp)

	// Parse the response
	var entries []Entry
	if err := parseResponse(resp, &entries); err != nil {
		return nil, nil, err
	}

	return entries, &paginationInfo, nil
}

// GetEntry retrieves a specific entry by ID.
func (c *Client) GetEntry(id int) (*Entry, error) {
	// Make the request
	resp, err := c.doRequest(http.MethodGet, fmt.Sprintf("/entries/%d.json", id), nil, nil)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var entry Entry
	if err := parseResponse(resp, &entry); err != nil {
		return nil, err
	}

	return &entry, nil
}

// Helper functions for common entry queries

// GetUnreadEntries retrieves all unread entries.
func (c *Client) GetUnreadEntries(page, perPage int) ([]Entry, *PaginationInfo, error) {
	options := map[string]interface{}{
		"read": false,
	}

	if page > 0 {
		options["page"] = page
	}

	if perPage > 0 {
		options["per_page"] = perPage
	}

	return c.GetEntries(options)
}

// GetStarredEntries retrieves all starred entries.
func (c *Client) GetStarredEntries(page, perPage int) ([]Entry, *PaginationInfo, error) {
	options := map[string]interface{}{
		"starred": true,
	}

	if page > 0 {
		options["page"] = page
	}

	if perPage > 0 {
		options["per_page"] = perPage
	}

	return c.GetEntries(options)
}

// GetEntriesSince retrieves all entries created after the specified time.
func (c *Client) GetEntriesSince(since time.Time, page, perPage int) ([]Entry, *PaginationInfo, error) {
	options := map[string]interface{}{
		"since": since,
	}

	if page > 0 {
		options["page"] = page
	}

	if perPage > 0 {
		options["per_page"] = perPage
	}

	return c.GetEntries(options)
}

// GetEntriesByIDs retrieves entries with the specified IDs.
func (c *Client) GetEntriesByIDs(ids []int) ([]Entry, *PaginationInfo, error) {
	options := map[string]interface{}{
		"ids": ids,
	}

	return c.GetEntries(options)
}
