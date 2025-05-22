package feedbin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-querystring/query" // For potential future options
)

// UnreadEntriesService handles operations related to unread entries.
type UnreadEntriesService struct {
	client *Client
}

// NewUnreadEntriesService creates a new service for unread entry operations.
func NewUnreadEntriesService(client *Client) *UnreadEntriesService {
	return &UnreadEntriesService{client: client}
}

// UnreadEntryListOptions specifies optional parameters for listing unread entries.
// Currently, the API doc for "GET /v2/unread_entries.json" does not specify any query parameters.
// Adding for consistency and potential future use.
type UnreadEntryListOptions struct {
	// PerPage int `url:"per_page,omitempty"` // Example if pagination was supported
	// Since string `url:"since,omitempty"` // Example if since was supported
}

// List retrieves the IDs of all unread entries.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/content/unread-entries.md#get-unread-entries
func (s *UnreadEntriesService) List(opts *UnreadEntryListOptions) ([]int64, *http.Response, error) {
	path := "unread_entries.json"
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

	var unreadEntryIDs []int64
	resp, err := s.client.Do(req, &unreadEntryIDs)
	if err != nil {
		return nil, resp, err
	}
	return unreadEntryIDs, resp, nil
}

// Create marks a list of entries as unread.
// The body should be a JSON array of entry IDs. e.g., `{"unread_entries": [1,2,3]}`
// However, the spec says "POST /v2/unread_entries.json Expects an array of entry IDs"
// and shows `curl -d '[1,2,3]'`. This implies the root object IS the array.
// Let's follow the curl example.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/content/unread-entries.md#mark-entries-as-unread
func (s *UnreadEntriesService) Create(entryIDs []int64) ([]int64, *http.Response, error) {
	if len(entryIDs) == 0 {
		return nil, nil, fmt.Errorf("entryIDs cannot be empty")
	}
	path := "unread_entries.json"

	// The API expects a direct JSON array of IDs, not a struct.
	req, err := s.client.NewRequest(http.MethodPost, path, entryIDs)
	if err != nil {
		return nil, nil, err
	}

	var createdUnreadEntryIDs []int64 // The response is also an array of IDs
	resp, err := s.client.Do(req, &createdUnreadEntryIDs)
	if err != nil {
		return nil, resp, err
	}
	return createdUnreadEntryIDs, resp, nil
}

// Delete marks a list of entries as read (i.e., removes them from unread).
// The body should be a JSON array of entry IDs. e.g., `{"unread_entries": [1,2,3]}`
// Similar to Create, spec says "DELETE /v2/unread_entries.json Expects an array of entry IDs"
// and shows `curl -X DELETE -d '[1,2,3]'`.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/content/unread-entries.md#mark-entries-as-read
func (s *UnreadEntriesService) Delete(entryIDs []int64) ([]int64, *http.Response, error) {
	if len(entryIDs) == 0 {
		return nil, nil, fmt.Errorf("entryIDs cannot be empty")
	}
	path := "unread_entries.json"

	// The API expects a direct JSON array of IDs for the body.
	req, err := s.client.NewRequest(http.MethodDelete, path, entryIDs)
	if err != nil {
		return nil, nil, err
	}

	var deletedUnreadEntryIDs []int64 // The response is also an array of IDs
	resp, err := s.client.Do(req, &deletedUnreadEntryIDs)
	if err != nil {
		// A 204 No Content might also be a valid response for DELETE operations
		// if nothing is returned in the body. The client.Do method handles 204s.
		// If the API *always* returns the list of IDs, this is fine.
		// If it can return 204, then `deletedUnreadEntryIDs` might be empty/nil.
		return nil, resp, err
	}
	return deletedUnreadEntryIDs, resp, nil
}
