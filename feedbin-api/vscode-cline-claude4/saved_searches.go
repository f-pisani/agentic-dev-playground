package feedbin

import (
	"context"
	"fmt"
)

// GetSavedSearches retrieves all saved searches for the authenticated user.
// Returns a slice of saved searches.
func (c *Client) GetSavedSearches(ctx context.Context) ([]SavedSearch, error) {
	var searches []SavedSearch
	_, err := c.get(ctx, "saved_searches.json", nil, &searches)
	if err != nil {
		return nil, err
	}

	return searches, nil
}

// GetSavedSearch retrieves a specific saved search by ID.
func (c *Client) GetSavedSearch(ctx context.Context, id int) (*SavedSearch, error) {
	path := fmt.Sprintf("saved_searches/%d.json", id)

	var search SavedSearch
	_, err := c.get(ctx, path, nil, &search)
	if err != nil {
		return nil, err
	}

	return &search, nil
}

// CreateSavedSearch creates a new saved search.
// Returns the created saved search or an error.
func (c *Client) CreateSavedSearch(ctx context.Context, name, query string) (*SavedSearch, error) {
	request := CreateSavedSearchRequest{
		Name:  name,
		Query: query,
	}

	var search SavedSearch
	_, err := c.post(ctx, "saved_searches.json", nil, request, &search)
	if err != nil {
		return nil, err
	}

	return &search, nil
}

// UpdateSavedSearch updates a saved search's name and/or query.
func (c *Client) UpdateSavedSearch(ctx context.Context, id int, name, query string) (*SavedSearch, error) {
	path := fmt.Sprintf("saved_searches/%d.json", id)
	request := UpdateSavedSearchRequest{
		Name:  name,
		Query: query,
	}

	var search SavedSearch
	_, err := c.patch(ctx, path, nil, request, &search)
	if err != nil {
		return nil, err
	}

	return &search, nil
}

// DeleteSavedSearch deletes a saved search by ID.
func (c *Client) DeleteSavedSearch(ctx context.Context, id int) error {
	path := fmt.Sprintf("saved_searches/%d.json", id)
	_, err := c.delete(ctx, path, nil, nil)
	return err
}

// GetSavedSearchByName retrieves a saved search by name.
// Returns the first saved search with the matching name, or nil if not found.
func (c *Client) GetSavedSearchByName(ctx context.Context, name string) (*SavedSearch, error) {
	searches, err := c.GetSavedSearches(ctx)
	if err != nil {
		return nil, err
	}

	for _, search := range searches {
		if search.Name == name {
			return &search, nil
		}
	}

	return nil, fmt.Errorf("saved search with name '%s' not found", name)
}

// DeleteSavedSearchByName deletes a saved search by name.
func (c *Client) DeleteSavedSearchByName(ctx context.Context, name string) error {
	search, err := c.GetSavedSearchByName(ctx, name)
	if err != nil {
		return err
	}

	return c.DeleteSavedSearch(ctx, search.ID)
}
