package feedbinapi

import (
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

// UpdatedEntriesService handles operations related to updated entries.
type UpdatedEntriesService struct {
	client *Client
}

// NewUpdatedEntriesService creates a new service for updated entry operations.
func NewUpdatedEntriesService(client *Client) *UpdatedEntriesService {
	return &UpdatedEntriesService{client: client}
}

// UpdatedEntryListOptions specifies optional parameters for listing updated entry IDs.
type UpdatedEntryListOptions struct {
	// The API docs state "GET /v2/updated_entries.json?since=YYYY-MM-DDTHH:MM:SS.SSSSSSZ"
	// 'since' is mandatory for this endpoint according to the description.
	Since string `url:"since"` // ISO 8601 date string, e.g., "2015-01-22T15:33:38.449047Z"
	// Other options like Page, PerPage are not mentioned for this endpoint.
}

// List retrieves the IDs of entries that have been updated since a given time.
// The 'since' parameter is required.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/content/updated-entries.md#get-updated-entries
func (s *UpdatedEntriesService) List(opts *UpdatedEntryListOptions) ([]int64, *http.Response, error) {
	if opts == nil || opts.Since == "" {
		return nil, nil, fmt.Errorf("'Since' parameter is required for listing updated entries")
	}

	// Validate that 'since' can be parsed as a time, though the API will ultimately validate it.
	// This is more of a sanity check for the client-side.
	_, err := ParseFeedbinTime(opts.Since)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid 'Since' time format: %w", err)
	}

	path := "updated_entries.json"
	v, err := query.Values(opts)
	if err != nil {
		return nil, nil, err
	}
	path = fmt.Sprintf("%s?%s", path, v.Encode())

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var updatedEntryIDs []int64
	resp, err := s.client.Do(req, &updatedEntryIDs)
	if err != nil {
		return nil, resp, err
	}
	return updatedEntryIDs, resp, nil
}

// Helper to format time for the 'since' parameter, if needed, though client.FormatFeedbinTime can be used directly.
/*
func FormatTimeForSinceParam(t time.Time) string {
	return FormatFeedbinTime(t) // Uses the existing helper from client.go
}
*/
