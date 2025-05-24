package feedbin

import (
	"fmt"
	"net/url"
)

// GetSavedSearches retrieves all saved searches
func (c *Client) GetSavedSearches() ([]SavedSearch, error) {
	var searches []SavedSearch
	err := c.get("/saved_searches.json", nil, &searches)
	return searches, err
}

// GetSavedSearch retrieves a single saved search by ID
func (c *Client) GetSavedSearch(id int) (*SavedSearch, error) {
	var search SavedSearch
	path := fmt.Sprintf("/saved_searches/%d.json", id)
	err := c.get(path, nil, &search)
	if err != nil {
		return nil, err
	}
	return &search, nil
}

// CreateSavedSearch creates a new saved search
func (c *Client) CreateSavedSearch(name, query string) (*SavedSearch, error) {
	body := struct {
		Name  string `json:"name"`
		Query string `json:"query"`
	}{
		Name:  name,
		Query: query,
	}

	var search SavedSearch
	err := c.post("/saved_searches.json", body, &search)
	if err != nil {
		return nil, err
	}
	return &search, nil
}

// UpdateSavedSearch updates an existing saved search
func (c *Client) UpdateSavedSearch(id int, name, query string) (*SavedSearch, error) {
	body := struct {
		Name  string `json:"name,omitempty"`
		Query string `json:"query,omitempty"`
	}{
		Name:  name,
		Query: query,
	}

	var search SavedSearch
	path := fmt.Sprintf("/saved_searches/%d.json", id)
	err := c.patch(path, body, &search)
	if err != nil {
		return nil, err
	}
	return &search, nil
}

// DeleteSavedSearch deletes a saved search
func (c *Client) DeleteSavedSearch(id int) error {
	path := fmt.Sprintf("/saved_searches/%d.json", id)
	return c.delete(path, nil)
}

// GetSavedSearchResults retrieves entries matching a saved search
func (c *Client) GetSavedSearchResults(id int, options ...QueryOption) ([]Entry, error) {
	search, err := c.GetSavedSearch(id)
	if err != nil {
		return nil, err
	}

	// Add the search query to the options
	searchOpts := append(options, WithQuery(search.Query))
	return c.GetEntries(searchOpts...)
}

// UpdateSavedSearchAlt updates a saved search using POST (alternative to PATCH)
func (c *Client) UpdateSavedSearchAlt(id int, name, query string) (*SavedSearch, error) {
	body := struct {
		Name  string `json:"name,omitempty"`
		Query string `json:"query,omitempty"`
	}{
		Name:  name,
		Query: query,
	}

	var search SavedSearch
	path := fmt.Sprintf("/saved_searches/%d/update.json", id)
	err := c.post(path, body, &search)
	if err != nil {
		return nil, err
	}
	return &search, nil
}

// WithQuery adds a search query to the request
func WithQuery(query string) QueryOption {
	return func(q url.Values) {
		if query != "" {
			q.Set("q", query)
		}
	}
}
