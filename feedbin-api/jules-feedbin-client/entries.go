package feedbin

import (
	"fmt"
	"net/http"
	// "net/url" // Not strictly needed if using go-querystring for all params

	"github.com/google/go-querystring/query"
)

// EntriesService handles operations related to feed entries.
type EntriesService struct {
	client *Client
}

// NewEntriesService creates a new service for entry related operations.
func NewEntriesService(client *Client) *EntriesService {
	return &EntriesService{client: client}
}

// EntryListOptions specifies the optional parameters to the EntriesService.List method.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/content/entries.md#get-entries
type EntryListOptions struct {
	ListOptions     // Embeds Page, PerPage, Since
	Read            *bool  `url:"read,omitempty"`             // Filter by read status
	Starred         *bool  `url:"starred,omitempty"`          // Filter by starred status
	Mode            string `url:"mode,omitempty"`             // "extended" for more details
	IncludeOriginal *bool  `url:"include_original,omitempty"` // Include original entry data if updated
	IDs             []int64 `url:"ids,comma,omitempty"`       // Retrieve specific entries by ID
	MinID           int64  `url:"min_id,omitempty"`           // Get entries with ID greater than min_id
	MaxID           int64  `url:"max_id,omitempty"`           // Get entries with ID less than max_id
	// Additional params like `keywords`, `tags` might be considered if API supports them on this endpoint.
}

// List retrieves all entries, optionally filtered.
// This is the main endpoint for entries.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/content/entries.md#get-entries
func (s *EntriesService) List(opts *EntryListOptions) ([]Entry, *http.Response, error) {
	path := "entries.json"
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

	var entries []Entry
	resp, err := s.client.Do(req, &entries)
	if err != nil {
		return nil, resp, err
	}
	return entries, resp, nil
}

// EntryGetOptions specifies optional parameters for getting a single entry.
type EntryGetOptions struct {
	Mode            string `url:"mode,omitempty"`             // "extended"
	IncludeOriginal *bool  `url:"include_original,omitempty"`
}

// Get retrieves a single entry by its ID.
// Docs: https://github.com/feedbin/feedbin-api/blob/master/content/entries.md#get-entry
func (s *EntriesService) Get(id int64, opts *EntryGetOptions) (*Entry, *http.Response, error) {
	path := fmt.Sprintf("entries/%d.json", id)
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

	var entry Entry
	resp, err := s.client.Do(req, &entry)
	if err != nil {
		return nil, resp, err
	}
	return &entry, resp, nil
}

// ListByFeed retrieves entries for a specific feed.
// Note: The documentation implies this is done via `/v2/feeds/{feed_id}/entries.json`.
// It also mentions that `entries.json` can be filtered by `feed_id` if that's how Feedbin implemented it,
// but the path structure is more common. Assuming the path structure for now.
// Let's check existing clients or assume this is equivalent to List with FeedID filter if not explicitly separate.
// The spec `entries.md` does not list `feed_id` as a query param for `entries.json`.
// It *does* say: "You can also get a list of entries for a single feed: GET /v2/feeds/{feed_id}/entries.json"
// So, we'll use that specific path.
func (s *EntriesService) ListByFeed(feedID int64, opts *EntryListOptions) ([]Entry, *http.Response, error) {
	// Remove FeedID if present in opts, as it's part of the path
	var effectiveOpts *EntryListOptions
	if opts != nil {
		// Clone opts to avoid modifying the original
		clonedOpts := *opts
		// No direct FeedID field in EntryListOptions to filter by path, so this is fine.
		effectiveOpts = &clonedOpts
	}


	path := fmt.Sprintf("feeds/%d/entries.json", feedID)
	if effectiveOpts != nil {
		v, err := query.Values(effectiveOpts)
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

	var entries []Entry
	resp, err := s.client.Do(req, &entries)
	if err != nil {
		return nil, resp, err
	}
	return entries, resp, nil
}
