package feedbin

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// EntriesService handles operations related to entries
type EntriesService struct {
	client *Client
}

// List returns all entries
// If options is nil, default options will be used
func (s *EntriesService) List(options *EntryListOptions) ([]Entry, error) {
	path := "entries.json"
	if options != nil {
		path = addEntryOptionsToPath(path, options)
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var entries []Entry
	_, err = s.client.Do(req, &entries)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

// Get returns a single entry by ID
func (s *EntriesService) Get(id int, options *EntryListOptions) (*Entry, error) {
	path := fmt.Sprintf("entries/%d.json", id)
	if options != nil {
		path = addEntryOptionsToPath(path, options)
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	entry := new(Entry)
	_, err = s.client.Do(req, entry)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

// ListByFeed returns all entries for a specific feed
func (s *EntriesService) ListByFeed(feedID int, options *EntryListOptions) ([]Entry, error) {
	path := fmt.Sprintf("feeds/%d/entries.json", feedID)
	if options != nil {
		path = addEntryOptionsToPath(path, options)
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var entries []Entry
	_, err = s.client.Do(req, &entries)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

// Helper function to add entry options to the path
func addEntryOptionsToPath(path string, options *EntryListOptions) string {
	values := url.Values{}

	if options.Page != nil {
		values.Add("page", strconv.Itoa(*options.Page))
	}

	if options.PerPage != nil {
		values.Add("per_page", strconv.Itoa(*options.PerPage))
	}

	if options.Since != nil {
		values.Add("since", options.Since.Format(time.RFC3339Nano))
	}

	if len(options.IDs) > 0 {
		// Convert IDs to strings and join with commas
		idStrings := make([]string, len(options.IDs))
		for i, id := range options.IDs {
			idStrings[i] = strconv.Itoa(id)
		}
		values.Add("ids", strings.Join(idStrings, ","))
	}

	if options.Read != nil {
		values.Add("read", strconv.FormatBool(*options.Read))
	}

	if options.Starred != nil {
		values.Add("starred", strconv.FormatBool(*options.Starred))
	}

	if options.Mode != "" {
		values.Add("mode", options.Mode)
	}

	if options.IncludeOriginal != nil {
		values.Add("include_original", strconv.FormatBool(*options.IncludeOriginal))
	}

	if options.IncludeEnclosure != nil {
		values.Add("include_enclosure", strconv.FormatBool(*options.IncludeEnclosure))
	}

	if options.IncludeContentDiff != nil {
		values.Add("include_content_diff", strconv.FormatBool(*options.IncludeContentDiff))
	}

	if len(values) > 0 {
		return path + "?" + values.Encode()
	}

	return path
}
