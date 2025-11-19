package feedbin

import (
	"context"
	"fmt"
	"strings"
)

// GetSavedSearches fetches all saved searches.
func (c *Client) GetSavedSearches(ctx context.Context) ([]SavedSearch, *Response, error) {
	req, err := c.NewRequest("GET", "saved_searches.json", nil)
	if err != nil {
		return nil, nil, err
	}

	var savedSearches []SavedSearch
	resp, err := c.Do(ctx, req, &savedSearches)
	if err != nil {
		return nil, resp, err
	}

	return savedSearches, resp, nil
}

// GetSavedSearch fetches a single saved search by ID.
// It returns a list of Entry IDs (int64) or full Entry objects if includeEntries is true.
// Since the return type varies, we use interface{}.
func (c *Client) GetSavedSearch(ctx context.Context, id int64, includeEntries bool, opt *PaginationOptions) (interface{}, *Response, error) {
	u := fmt.Sprintf("saved_searches/%d.json", id)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	if includeEntries {
		// Hacky way to add query param since addOptions handles standard pagination
		if strings.Contains(u, "?") {
			u += "&include_entries=true"
		} else {
			u += "?include_entries=true"
		}
	}

	req, err := c.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	if includeEntries {
		var entries []Entry
		resp, err := c.Do(ctx, req, &entries)
		if err != nil {
			return nil, resp, err
		}
		return entries, resp, nil
	} else {
		var ids []int64
		resp, err := c.Do(ctx, req, &ids)
		if err != nil {
			return nil, resp, err
		}
		return ids, resp, nil
	}
}

type savedSearchCreateRequest struct {
	Name  string `json:"name"`
	Query string `json:"query"`
}

// CreateSavedSearch creates a new saved search.
func (c *Client) CreateSavedSearch(ctx context.Context, name, query string) (*SavedSearch, *Response, error) {
	reqBody := &savedSearchCreateRequest{Name: name, Query: query}
	req, err := c.NewRequest("POST", "saved_searches.json", reqBody)
	if err != nil {
		return nil, nil, err
	}

	var savedSearch SavedSearch
	resp, err := c.Do(ctx, req, &savedSearch)
	if err != nil {
		return nil, resp, err
	}

	return &savedSearch, resp, nil
}

type savedSearchUpdateRequest struct {
	Name string `json:"name"`
}

// UpdateSavedSearch updates a saved search's name.
func (c *Client) UpdateSavedSearch(ctx context.Context, id int64, name string) (*SavedSearch, *Response, error) {
	u := fmt.Sprintf("saved_searches/%d.json", id)
	reqBody := &savedSearchUpdateRequest{Name: name}
	req, err := c.NewRequest("PATCH", u, reqBody)
	if err != nil {
		return nil, nil, err
	}

	var savedSearch SavedSearch
	resp, err := c.Do(ctx, req, &savedSearch)
	if err != nil {
		return nil, resp, err
	}

	return &savedSearch, resp, nil
}

// DeleteSavedSearch deletes a saved search.
func (c *Client) DeleteSavedSearch(ctx context.Context, id int64) (*Response, error) {
	u := fmt.Sprintf("saved_searches/%d.json", id)
	req, err := c.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(ctx, req, nil)
}
