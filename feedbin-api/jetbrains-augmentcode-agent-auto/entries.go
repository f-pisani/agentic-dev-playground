package feedbin

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// EntriesService handles communication with the entries related
// methods of the Feedbin API.
type EntriesService struct {
	client *Client
}

// EntryListOptions specifies the optional parameters to the
// EntriesService.List method.
type EntryListOptions struct {
	Page               int
	Since              time.Time
	IDs                []int64
	Read               *bool
	Starred            *bool
	PerPage            int
	Mode               string // "extended" is the only available mode
	IncludeOriginal    bool
	IncludeEnclosure   bool
	IncludeContentDiff bool
}

// List returns all entries for the user.
func (s *EntriesService) List(opts *EntryListOptions) ([]*Entry, *http.Response, error) {
	u := "/v2/entries.json"

	if opts != nil {
		params := url.Values{}

		if opts.Page > 0 {
			params.Add("page", strconv.Itoa(opts.Page))
		}

		if !opts.Since.IsZero() {
			params.Add("since", opts.Since.Format(time.RFC3339Nano))
		}

		if len(opts.IDs) > 0 {
			var ids []string
			for _, id := range opts.IDs {
				ids = append(ids, strconv.FormatInt(id, 10))
			}
			params.Add("ids", strings.Join(ids, ","))
		}

		if opts.Read != nil {
			params.Add("read", strconv.FormatBool(*opts.Read))
		}

		if opts.Starred != nil {
			params.Add("starred", strconv.FormatBool(*opts.Starred))
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

		if len(params) > 0 {
			u += "?" + params.Encode()
		}
	}

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var entries []*Entry
	resp, err := s.client.Do(req, &entries)
	if err != nil {
		return nil, resp, err
	}

	return entries, resp, nil
}

// ListByFeed returns all entries for a specific feed.
func (s *EntriesService) ListByFeed(feedID int64, opts *EntryListOptions) ([]*Entry, *http.Response, error) {
	u := fmt.Sprintf("/v2/feeds/%d/entries.json", feedID)

	if opts != nil {
		params := url.Values{}

		if opts.Page > 0 {
			params.Add("page", strconv.Itoa(opts.Page))
		}

		if !opts.Since.IsZero() {
			params.Add("since", opts.Since.Format(time.RFC3339Nano))
		}

		if opts.Read != nil {
			params.Add("read", strconv.FormatBool(*opts.Read))
		}

		if opts.Starred != nil {
			params.Add("starred", strconv.FormatBool(*opts.Starred))
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

		if len(params) > 0 {
			u += "?" + params.Encode()
		}
	}

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var entries []*Entry
	resp, err := s.client.Do(req, &entries)
	if err != nil {
		return nil, resp, err
	}

	return entries, resp, nil
}

// Get returns a single entry.
func (s *EntriesService) Get(id int64, opts *EntryListOptions) (*Entry, *http.Response, error) {
	u := fmt.Sprintf("/v2/entries/%d.json", id)

	if opts != nil {
		params := url.Values{}

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

		if len(params) > 0 {
			u += "?" + params.Encode()
		}
	}

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
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
