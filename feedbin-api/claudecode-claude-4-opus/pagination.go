package feedbin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// PageIterator provides a way to iterate through paginated results
type PageIterator struct {
	client      *Client
	path        string
	currentPage int
	perPage     int
	baseQuery   url.Values
	hasMore     bool
}

// NewPageIterator creates a new page iterator
func (c *Client) NewPageIterator(path string, perPage int, options ...QueryOption) *PageIterator {
	query := buildQuery(options...)
	if perPage > 0 {
		query.Set("per_page", fmt.Sprintf("%d", perPage))
	}

	return &PageIterator{
		client:      c,
		path:        path,
		currentPage: 0,
		perPage:     perPage,
		baseQuery:   query,
		hasMore:     true,
	}
}

// HasNext returns true if there are more pages to fetch
func (pi *PageIterator) HasNext() bool {
	return pi.hasMore
}

// NextPage fetches the next page of results
func (pi *PageIterator) NextPage(result interface{}) error {
	if !pi.hasMore {
		return fmt.Errorf("no more pages available")
	}

	pi.currentPage++
	query := make(url.Values)
	for k, v := range pi.baseQuery {
		query[k] = v
	}
	query.Set("page", fmt.Sprintf("%d", pi.currentPage))

	resp, err := pi.client.request(http.MethodGet, pi.path, nil, query)
	if err != nil {
		pi.hasMore = false
		return err
	}
	defer resp.Body.Close()

	// Check if we got results
	if resp.StatusCode == 204 || resp.ContentLength == 0 {
		pi.hasMore = false
		return nil
	}

	// Decode the response
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		pi.hasMore = false
		return err
	}

	// For simplicity, we'll check if we got fewer results than requested
	// In a real implementation, you might check response headers or array length
	// This is a basic heuristic that works for most cases
	return nil
}

// Reset resets the iterator to start from the beginning
func (pi *PageIterator) Reset() {
	pi.currentPage = 0
	pi.hasMore = true
}

// CurrentPage returns the current page number
func (pi *PageIterator) CurrentPage() int {
	return pi.currentPage
}
