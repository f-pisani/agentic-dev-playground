package feedbin

import (
	"context"
	"fmt"
)

// GetSavedSearches retrieves all saved searches for the authenticated user
func (c *Client) GetSavedSearches(ctx context.Context) ([]SavedSearch, error) {
	var savedSearches []SavedSearch
	_, err := c.makeRequest(ctx, "GET", "/saved_searches.json", nil, &savedSearches)
	if err != nil {
		return nil, err
	}

	return savedSearches, nil
}

// GetSavedSearch retrieves results for a specific saved search by ID
func (c *Client) GetSavedSearch(ctx context.Context, id int, opts *SavedSearchOptions) (interface{}, error) {
	path := fmt.Sprintf("/saved_searches/%d.json", id)
	if opts != nil {
		params := buildQueryParams(opts)
		path = addQueryParams(path, params)
	}

	if opts != nil && opts.IncludeEntries != nil && *opts.IncludeEntries {
		var entries []Entry
		_, err := c.makeRequest(ctx, "GET", path, nil, &entries)
		if err != nil {
			return nil, err
		}
		return entries, nil
	}

	var entryIDs []int
	_, err := c.makeRequest(ctx, "GET", path, nil, &entryIDs)
	if err != nil {
		return nil, err
	}

	return entryIDs, nil
}

// GetSavedSearchEntryIDs retrieves entry IDs for a specific saved search
func (c *Client) GetSavedSearchEntryIDs(ctx context.Context, id int, opts *SavedSearchOptions) ([]int, error) {
	result, err := c.GetSavedSearch(ctx, id, opts)
	if err != nil {
		return nil, err
	}

	if entryIDs, ok := result.([]int); ok {
		return entryIDs, nil
	}

	return nil, fmt.Errorf("expected entry IDs but got entries")
}

// GetSavedSearchEntries retrieves full entries for a specific saved search
func (c *Client) GetSavedSearchEntries(ctx context.Context, id int, opts *SavedSearchOptions) ([]Entry, error) {
	if opts == nil {
		opts = &SavedSearchOptions{}
	}
	opts.IncludeEntries = Bool(true)

	result, err := c.GetSavedSearch(ctx, id, opts)
	if err != nil {
		return nil, err
	}

	if entries, ok := result.([]Entry); ok {
		return entries, nil
	}

	return nil, fmt.Errorf("expected entries but got entry IDs")
}

// CreateSavedSearch creates a new saved search
func (c *Client) CreateSavedSearch(ctx context.Context, name, query string) (*SavedSearch, error) {
	if name == "" {
		return nil, &ValidationError{
			Field:   "name",
			Message: "name is required",
		}
	}

	if query == "" {
		return nil, &ValidationError{
			Field:   "query",
			Message: "query is required",
		}
	}

	req := &CreateSavedSearchRequest{
		Name:  name,
		Query: query,
	}

	var savedSearch SavedSearch
	_, err := c.makeRequest(ctx, "POST", "/saved_searches.json", req, &savedSearch)
	if err != nil {
		return nil, err
	}

	return &savedSearch, nil
}

// UpdateSavedSearch updates a saved search's name
func (c *Client) UpdateSavedSearch(ctx context.Context, id int, name string) (*SavedSearch, error) {
	if name == "" {
		return nil, &ValidationError{
			Field:   "name",
			Message: "name is required",
		}
	}

	path := fmt.Sprintf("/saved_searches/%d.json", id)
	req := &UpdateSavedSearchRequest{
		Name: name,
	}

	var savedSearch SavedSearch
	_, err := c.makeRequest(ctx, "PATCH", path, req, &savedSearch)
	if err != nil {
		return nil, err
	}

	return &savedSearch, nil
}

// UpdateSavedSearchPOST updates a saved search using the POST alternative endpoint
func (c *Client) UpdateSavedSearchPOST(ctx context.Context, id int, name string) (*SavedSearch, error) {
	if name == "" {
		return nil, &ValidationError{
			Field:   "name",
			Message: "name is required",
		}
	}

	path := fmt.Sprintf("/saved_searches/%d/update.json", id)
	req := &UpdateSavedSearchRequest{
		Name: name,
	}

	var savedSearch SavedSearch
	_, err := c.makeRequest(ctx, "POST", path, req, &savedSearch)
	if err != nil {
		return nil, err
	}

	return &savedSearch, nil
}

// DeleteSavedSearch deletes a saved search
func (c *Client) DeleteSavedSearch(ctx context.Context, id int) error {
	path := fmt.Sprintf("/saved_searches/%d.json", id)

	_, err := c.makeRequest(ctx, "DELETE", path, nil, nil)
	return err
}

// Convenience methods without context (use background context)

func (c *Client) GetSavedSearchesWithoutContext() ([]SavedSearch, error) {
	return c.GetSavedSearches(context.Background())
}

func (c *Client) GetSavedSearchWithoutContext(id int, opts *SavedSearchOptions) (interface{}, error) {
	return c.GetSavedSearch(context.Background(), id, opts)
}

func (c *Client) GetSavedSearchEntryIDsWithoutContext(id int, opts *SavedSearchOptions) ([]int, error) {
	return c.GetSavedSearchEntryIDs(context.Background(), id, opts)
}

func (c *Client) GetSavedSearchEntriesWithoutContext(id int, opts *SavedSearchOptions) ([]Entry, error) {
	return c.GetSavedSearchEntries(context.Background(), id, opts)
}

func (c *Client) CreateSavedSearchWithoutContext(name, query string) (*SavedSearch, error) {
	return c.CreateSavedSearch(context.Background(), name, query)
}

func (c *Client) UpdateSavedSearchWithoutContext(id int, name string) (*SavedSearch, error) {
	return c.UpdateSavedSearch(context.Background(), id, name)
}

func (c *Client) DeleteSavedSearchWithoutContext(id int) error {
	return c.DeleteSavedSearch(context.Background(), id)
}
