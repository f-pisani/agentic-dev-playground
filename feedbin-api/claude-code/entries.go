package feedbin

import (
	"fmt"
	"net/http"
	neturl "net/url"
	"strconv"
	"strings"
)

// EntriesService handles communication with the entries related
// endpoints of the Feedbin API
type EntriesService struct {
	client *Client
}

// EntryListOptions specifies the optional parameters to the
// EntriesService.List method
type EntryListOptions struct {
	ListOptions
	IDs                []int  // Filter by entry IDs
	Read               *bool  // Filter by read status
	Starred            *bool  // Filter by starred status
	PerPage            int    // Limit results to this many per page
	Mode               string // "extended" mode is available
	IncludeOriginal    bool   // Include original entry data
	IncludeEnclosure   bool   // Include podcast/RSS enclosure data
	IncludeContentDiff bool   // Include diff of changed content
}

// List returns entries
func (s *EntriesService) List(opts *EntryListOptions) ([]*Entry, *PaginationInfo, *http.Response, error) {
	url := "entries.json"
	params := neturl.Values{}

	if opts != nil {
		if opts.Page > 0 {
			params.Add("page", strconv.Itoa(opts.Page))
		}
		if !opts.Since.IsZero() {
			params.Add("since", opts.Since.Format("2006-01-02T15:04:05.000000Z"))
		}
		if len(opts.IDs) > 0 {
			ids := make([]string, len(opts.IDs))
			for i, id := range opts.IDs {
				ids[i] = strconv.Itoa(id)
			}
			params.Add("ids", strings.Join(ids, ","))
		}
		if opts.Read != nil {
			if *opts.Read {
				params.Add("read", "true")
			} else {
				params.Add("read", "false")
			}
		}
		if opts.Starred != nil {
			if *opts.Starred {
				params.Add("starred", "true")
			} else {
				params.Add("starred", "false")
			}
		}
		if opts.PerPage > 0 {
			params.Add("per_page", strconv.Itoa(opts.PerPage))
		}
		if opts.Mode != "" {
			params.Add("mode", opts.Mode)
		}
		if opts.IncludeOriginal {
			params.Add("include_original", "true")
		}
		if opts.IncludeEnclosure {
			params.Add("include_enclosure", "true")
		}
		if opts.IncludeContentDiff {
			params.Add("include_content_diff", "true")
		}
	}

	if len(params) > 0 {
		url += "?" + params.Encode()
	}

	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	var entries []*Entry
	resp, err := s.client.Do(req, &entries)
	if err != nil {
		return nil, nil, resp, err
	}

	pagination := s.client.GetPagination(resp)

	return entries, pagination, resp, nil
}

// ListByFeed retrieves entries for a specific feed
func (s *EntriesService) ListByFeed(feedID int, opts *EntryListOptions) ([]*Entry, *PaginationInfo, *http.Response, error) {
	url := fmt.Sprintf("feeds/%d/entries.json", feedID)
	params := neturl.Values{}

	if opts != nil {
		if opts.Page > 0 {
			params.Add("page", strconv.Itoa(opts.Page))
		}
		if !opts.Since.IsZero() {
			params.Add("since", opts.Since.Format("2006-01-02T15:04:05.000000Z"))
		}
		if opts.Read != nil {
			if *opts.Read {
				params.Add("read", "true")
			} else {
				params.Add("read", "false")
			}
		}
		if opts.Starred != nil {
			if *opts.Starred {
				params.Add("starred", "true")
			} else {
				params.Add("starred", "false")
			}
		}
		if opts.PerPage > 0 {
			params.Add("per_page", strconv.Itoa(opts.PerPage))
		}
		if opts.Mode != "" {
			params.Add("mode", opts.Mode)
		}
		if opts.IncludeOriginal {
			params.Add("include_original", "true")
		}
		if opts.IncludeEnclosure {
			params.Add("include_enclosure", "true")
		}
		if opts.IncludeContentDiff {
			params.Add("include_content_diff", "true")
		}
	}

	if len(params) > 0 {
		url += "?" + params.Encode()
	}

	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	var entries []*Entry
	resp, err := s.client.Do(req, &entries)
	if err != nil {
		return nil, nil, resp, err
	}

	pagination := s.client.GetPagination(resp)

	return entries, pagination, resp, nil
}

// Get gets a single entry by ID
func (s *EntriesService) Get(id int, opts *EntryListOptions) (*Entry, *http.Response, error) {
	url := fmt.Sprintf("entries/%d.json", id)
	params := neturl.Values{}

	if opts != nil {
		if opts.Mode != "" {
			params.Add("mode", opts.Mode)
		}
		if opts.IncludeOriginal {
			params.Add("include_original", "true")
		}
		if opts.IncludeEnclosure {
			params.Add("include_enclosure", "true")
		}
		if opts.IncludeContentDiff {
			params.Add("include_content_diff", "true")
		}
	}

	if len(params) > 0 {
		url += "?" + params.Encode()
	}

	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	entry := new(Entry)
	resp, err := s.client.Do(req, entry)
	if err != nil {
		return nil, resp, err
	}

	return entry, resp, nil
}
