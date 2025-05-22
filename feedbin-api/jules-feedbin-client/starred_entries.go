package feedbin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-querystring/query" // For potential future options
)

// StarredEntriesService handles operations related to starred entries.
type StarredEntriesService struct {
	client *Client
}

// NewStarredEntriesService creates a new service for starred entry operations.
func NewStarredEntriesService(client *Client) *StarredEntriesService {
	return &StarredEntriesService{client: client}
}

// StarredEntryListOptions specifies optional parameters for listing starred entries.
// The API doc for "GET /v2/starred_entries.json" does not specify any query parameters.
type StarredEntryListOptions struct {
	// PerPage int `url:"per_page,omitempty"` // Example
	// Since string `url:"since,omitempty"`   // Example
}

// List retrieves the IDs of all starred entries.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/content/starred-entries.md#get-starred-entries
func (s *StarredEntriesService) List(opts *StarredEntryListOptions) ([]int64, *http.Response, error) {
	path := "starred_entries.json"
	if opts != nil {
		v, err := query.Values(opts)
		if err != nil {
			return nil, nil, err
		}
		if params := v.Encode(); params != "" {
			path = fmt.Sprintf("%s?%s", path, params)
		}
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var starredEntryIDs []int64
	resp, err := s.client.Do(req, &starredEntryIDs)
	if err != nil {
		return nil, resp, err
	}
	return starredEntryIDs, resp, nil
}

// Create stars a list of entries.
// Expects a JSON array of entry IDs in the request body.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/content/starred-entries.md#star-entries
func (s *StarredEntriesService) Create(entryIDs []int64) ([]int64, *http.Response, error) {
	if len(entryIDs) == 0 {
		return nil, nil, fmt.Errorf("entryIDs cannot be empty")
	}
	path := "starred_entries.json"

	req, err := s.client.NewRequest(http.MethodPost, path, entryIDs)
	if err != nil {
		return nil, nil, err
	}

	var createdStarredEntryIDs []int64
	resp, err := s.client.Do(req, &createdStarredEntryIDs)
	if err != nil {
		return nil, resp, err
	}
	return createdStarredEntryIDs, resp, nil
}

// Delete unstars a list of entries.
// Expects a JSON array of entry IDs in the request body.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/content/starred-entries.md#unstar-entries
func (s *StarredEntriesService) Delete(entryIDs []int64) ([]int64, *http.Response, error) {
	if len(entryIDs) == 0 {
		return nil, nil, fmt.Errorf("entryIDs cannot be empty")
	}
	path := "starred_entries.json"

	req, err := s.client.NewRequest(http.MethodDelete, path, entryIDs)
	if err != nil {
		return nil, nil, err
	}

	var deletedStarredEntryIDs []int64
	resp, err := s.client.Do(req, &deletedStarredEntryIDs)
	if err != nil {
		return nil, resp, err
	}
	return deletedStarredEntryIDs, resp, nil
}
