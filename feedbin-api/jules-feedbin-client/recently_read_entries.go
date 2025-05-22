package feedbin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-querystring/query"
)

// RecentlyReadEntriesService handles operations related to recently read entries.
type RecentlyReadEntriesService struct {
	client *Client
}

// NewRecentlyReadEntriesService creates a new service for recently read entry operations.
func NewRecentlyReadEntriesService(client *Client) *RecentlyReadEntriesService {
	return &RecentlyReadEntriesService{client: client}
}

// RecentlyReadEntryListOptions specifies optional parameters for listing recently read entries.
// The API doc for "GET /v2/recently_read_entries.json" does not specify query parameters.
type RecentlyReadEntryListOptions struct {
	ListOptions // Embeds Page, PerPage, Since (though not specified in docs)
}

// List retrieves recently read entries.
// The response is an array of `RecentlyReadEntry` objects.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/content/recently-read-entries.md#get-recently-read-entries
func (s *RecentlyReadEntriesService) List(opts *RecentlyReadEntryListOptions) ([]RecentlyReadEntry, *http.Response, error) {
	path := "recently_read_entries.json"
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

	var entries []RecentlyReadEntry
	resp, err := s.client.Do(req, &entries)
	if err != nil {
		return nil, resp, err
	}
	return entries, resp, nil
}

// CreateRecentlyReadEntryOptions defines the structure for the POST request.
// The API expects an array of objects, but the example shows a single object for a single entry.
// Let's assume for now it's a single object post, or the API handles an array of these if multiple are sent.
// The spec says "POST /v2/recently_read_entries.json" and expects `entry_id` and optionally `interaction`.
// It's more likely this endpoint records *one* interaction at a time or a batch.
// Let's make a method for a single interaction, and one for batch if needed.
// For now, sticking to simpler single interaction.
// The example shows `curl -d '{"entry_id": 1, "interaction": "mark_as_read"}'`
// This implies the endpoint itself is for creating *one* such record.
// The response is the created record.

// Create marks an entry as recently read or interacted with.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/content/recently-read-entries.md#create-recently-read-entry
func (s *RecentlyReadEntriesService) Create(entryID int64, interaction *string) (*RecentlyReadEntry, *http.Response, error) {
	if entryID == 0 {
		return nil, nil, fmt.Errorf("entryID is required")
	}

	body := struct {
		EntryID     int64   `json:"entry_id"`
		Interaction *string `json:"interaction,omitempty"`
	}{
		EntryID:     entryID,
		Interaction: interaction,
	}

	path := "recently_read_entries.json"
	req, err := s.client.NewRequest(http.MethodPost, path, body)
	if err != nil {
		return nil, nil, err
	}

	var createdEntry RecentlyReadEntry
	resp, err := s.client.Do(req, &createdEntry)
	if err != nil {
		return nil, resp, err
	}
	return &createdEntry, resp, nil
}

// Note: The API spec for "Create" shows creating a single recently_read_entry.
// If batch creation is needed and supported differently, that would be a separate method.
// For example, if it expected an array of these objects `[]RecentlyReadEntryPayload`.
