// Package feedbin provides a Go client for the Feedbin API v2.
package feedbin

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// EntryService provides access to the entry-related endpoints of the Feedbin API.
type EntryService struct {
	client *Client
}

// EntryOptions represents options for retrieving entries.
type EntryOptions struct {
	PageOptions
	Read               *bool  `url:"read,omitempty"`
	Starred            *bool  `url:"starred,omitempty"`
	IDs                []int  `url:"ids,omitempty"`
	Mode               string `url:"mode,omitempty"`
	IncludeOriginal    bool   `url:"include_original,omitempty"`
	IncludeEnclosure   bool   `url:"include_enclosure,omitempty"`
	IncludeContentDiff bool   `url:"include_content_diff,omitempty"`
}

// GetEntries retrieves all entries for the authenticated user.
func (s *EntryService) GetEntries(options *EntryOptions) ([]Entry, *http.Response, error) {
	path := "/entries.json"
	path = s.addEntryParams(path, options)

	req, err := s.client.NewRequest("GET", path, nil)
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

// GetEntry retrieves a specific entry by ID.
func (s *EntryService) GetEntry(id int, options *EntryOptions) (*Entry, *http.Response, error) {
	path := fmt.Sprintf("/entries/%d.json", id)
	path = s.addEntryParams(path, options)

	req, err := s.client.NewRequest("GET", path, nil)
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

// GetFeedEntries retrieves all entries for a specific feed.
func (s *EntryService) GetFeedEntries(feedID int, options *EntryOptions) ([]Entry, *http.Response, error) {
	path := fmt.Sprintf("/feeds/%d/entries.json", feedID)
	path = s.addEntryParams(path, options)

	req, err := s.client.NewRequest("GET", path, nil)
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

// GetEntriesByIDs retrieves entries with the specified IDs.
func (s *EntryService) GetEntriesByIDs(ids []int, options *EntryOptions) ([]Entry, *http.Response, error) {
	if len(ids) == 0 {
		return nil, nil, fmt.Errorf("no entry IDs provided")
	}

	if len(ids) > 100 {
		return nil, nil, fmt.Errorf("maximum of 100 entry IDs allowed")
	}

	if options == nil {
		options = &EntryOptions{}
	}
	options.IDs = ids

	return s.GetEntries(options)
}

// GetUnreadEntries retrieves all unread entries.
func (s *EntryService) GetUnreadEntries(options *EntryOptions) ([]Entry, *http.Response, error) {
	if options == nil {
		options = &EntryOptions{}
	}
	read := false
	options.Read = &read

	return s.GetEntries(options)
}

// GetStarredEntries retrieves all starred entries.
func (s *EntryService) GetStarredEntries(options *EntryOptions) ([]Entry, *http.Response, error) {
	if options == nil {
		options = &EntryOptions{}
	}
	starred := true
	options.Starred = &starred

	return s.GetEntries(options)
}

// GetEntriesSince retrieves all entries created after the specified time.
func (s *EntryService) GetEntriesSince(since time.Time, options *EntryOptions) ([]Entry, *http.Response, error) {
	if options == nil {
		options = &EntryOptions{}
	}
	options.Since = since.Format(time.RFC3339Nano)

	return s.GetEntries(options)
}

// GetEntriesExtended retrieves all entries with extended metadata.
func (s *EntryService) GetEntriesExtended(options *EntryOptions) ([]Entry, *http.Response, error) {
	if options == nil {
		options = &EntryOptions{}
	}
	options.Mode = "extended"

	return s.GetEntries(options)
}

// GetEntriesWithOriginal retrieves all entries including original entry data if updated.
func (s *EntryService) GetEntriesWithOriginal(options *EntryOptions) ([]Entry, *http.Response, error) {
	if options == nil {
		options = &EntryOptions{}
	}
	options.IncludeOriginal = true

	return s.GetEntries(options)
}

// GetEntriesWithEnclosure retrieves all entries including podcast/RSS enclosure data.
func (s *EntryService) GetEntriesWithEnclosure(options *EntryOptions) ([]Entry, *http.Response, error) {
	if options == nil {
		options = &EntryOptions{}
	}
	options.IncludeEnclosure = true

	return s.GetEntries(options)
}

// GetEntriesWithContentDiff retrieves all entries including a diff of changed content.
func (s *EntryService) GetEntriesWithContentDiff(options *EntryOptions) ([]Entry, *http.Response, error) {
	if options == nil {
		options = &EntryOptions{}
	}
	options.IncludeContentDiff = true

	return s.GetEntries(options)
}

// addEntryParams adds entry-specific parameters to a URL.
func (s *EntryService) addEntryParams(url string, options *EntryOptions) string {
	if options == nil {
		return url
	}

	params := []string{}

	// Add PageOptions parameters
	if options.Page > 0 {
		params = append(params, "page="+strconv.Itoa(options.Page))
	}

	if options.PerPage > 0 {
		params = append(params, "per_page="+strconv.Itoa(options.PerPage))
	}

	if options.Since != "" {
		params = append(params, "since="+options.Since)
	}

	// Add EntryOptions parameters
	if options.Read != nil {
		params = append(params, "read="+strconv.FormatBool(*options.Read))
	}

	if options.Starred != nil {
		params = append(params, "starred="+strconv.FormatBool(*options.Starred))
	}

	if len(options.IDs) > 0 {
		idStrs := make([]string, len(options.IDs))
		for i, id := range options.IDs {
			idStrs[i] = strconv.Itoa(id)
		}
		params = append(params, "ids="+strings.Join(idStrs, ","))
	}

	if options.Mode != "" {
		params = append(params, "mode="+options.Mode)
	}

	if options.IncludeOriginal {
		params = append(params, "include_original=true")
	}

	if options.IncludeEnclosure {
		params = append(params, "include_enclosure=true")
	}

	if options.IncludeContentDiff {
		params = append(params, "include_content_diff=true")
	}

	if len(params) == 0 {
		return url
	}

	if strings.Contains(url, "?") {
		return url + "&" + strings.Join(params, "&")
	}

	return url + "?" + strings.Join(params, "&")
}
