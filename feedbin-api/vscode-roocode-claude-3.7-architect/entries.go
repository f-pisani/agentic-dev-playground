package feedbin

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// EntriesService handles entry-related operations
type EntriesService struct {
	client *Client
}

// ListEntriesOptions represents options for listing entries
type ListEntriesOptions struct {
	Page               *int
	Since              *time.Time
	IDs                []int
	Read               *bool
	Starred            *bool
	PerPage            *int
	Mode               string
	IncludeOriginal    bool
	IncludeEnclosure   bool
	IncludeContentDiff bool
}

// List returns entries based on the provided options
func (s *EntriesService) List(opts *ListEntriesOptions) ([]Entry, *PaginationLinks, error) {
	path := "/entries.json"

	if opts != nil {
		query := url.Values{}

		if opts.Page != nil {
			query.Set("page", strconv.Itoa(*opts.Page))
		}

		if opts.Since != nil {
			query.Set("since", opts.Since.Format(time.RFC3339Nano))
		}

		if len(opts.IDs) > 0 {
			if len(opts.IDs) > 100 {
				return nil, nil, fmt.Errorf("maximum of 100 IDs allowed")
			}

			ids := make([]string, len(opts.IDs))
			for i, id := range opts.IDs {
				ids[i] = strconv.Itoa(id)
			}

			query.Set("ids", strings.Join(ids, ","))
		}

		if opts.Read != nil {
			query.Set("read", strconv.FormatBool(*opts.Read))
		}

		if opts.Starred != nil {
			query.Set("starred", strconv.FormatBool(*opts.Starred))
		}

		if opts.PerPage != nil {
			query.Set("per_page", strconv.Itoa(*opts.PerPage))
		}

		if opts.Mode != "" {
			query.Set("mode", opts.Mode)
		}

		if opts.IncludeOriginal {
			query.Set("include_original", "true")
		}

		if opts.IncludeEnclosure {
			query.Set("include_enclosure", "true")
		}

		if opts.IncludeContentDiff {
			query.Set("include_content_diff", "true")
		}

		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var entries []Entry
	resp, err := s.client.do(req, &entries)
	if err != nil {
		return nil, nil, err
	}

	// Parse pagination links
	links := parsePaginationLinks(resp.Header.Get("Link"))

	return entries, &links, nil
}

// ListByFeed returns entries for a specific feed
func (s *EntriesService) ListByFeed(feedID int, opts *ListEntriesOptions) ([]Entry, *PaginationLinks, error) {
	path := fmt.Sprintf("/feeds/%d/entries.json", feedID)

	if opts != nil {
		query := url.Values{}

		if opts.Page != nil {
			query.Set("page", strconv.Itoa(*opts.Page))
		}

		if opts.Since != nil {
			query.Set("since", opts.Since.Format(time.RFC3339Nano))
		}

		if opts.Read != nil {
			query.Set("read", strconv.FormatBool(*opts.Read))
		}

		if opts.Starred != nil {
			query.Set("starred", strconv.FormatBool(*opts.Starred))
		}

		if opts.PerPage != nil {
			query.Set("per_page", strconv.Itoa(*opts.PerPage))
		}

		if opts.Mode != "" {
			query.Set("mode", opts.Mode)
		}

		if opts.IncludeOriginal {
			query.Set("include_original", "true")
		}

		if opts.IncludeEnclosure {
			query.Set("include_enclosure", "true")
		}

		if opts.IncludeContentDiff {
			query.Set("include_content_diff", "true")
		}

		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var entries []Entry
	resp, err := s.client.do(req, &entries)
	if err != nil {
		return nil, nil, err
	}

	// Parse pagination links
	links := parsePaginationLinks(resp.Header.Get("Link"))

	return entries, &links, nil
}

// Get returns a specific entry
func (s *EntriesService) Get(id int, opts *ListEntriesOptions) (*Entry, error) {
	path := fmt.Sprintf("/entries/%d.json", id)

	if opts != nil {
		query := url.Values{}

		if opts.Mode != "" {
			query.Set("mode", opts.Mode)
		}

		if opts.IncludeOriginal {
			query.Set("include_original", "true")
		}

		if opts.IncludeEnclosure {
			query.Set("include_enclosure", "true")
		}

		if opts.IncludeContentDiff {
			query.Set("include_content_diff", "true")
		}

		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	entry := new(Entry)
	_, err = s.client.do(req, entry)
	if err != nil {
		return nil, err
	}

	return entry, nil
}
