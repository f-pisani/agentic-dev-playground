package feedbin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// EntriesService handles communication with the entries related
// endpoints of the Feedbin API
type EntriesService struct {
	client *Client
}

// EntryOptions specifies optional parameters for entry requests
type EntryOptions struct {
	PageOptions
	IDs                []int  `url:"ids,omitempty"`
	Read               *bool  `url:"read,omitempty"`
	Starred            *bool  `url:"starred,omitempty"`
	Mode               string `url:"mode,omitempty"` // only valid value is "extended"
	IncludeOriginal    *bool  `url:"include_original,omitempty"`
	IncludeEnclosure   *bool  `url:"include_enclosure,omitempty"`
	IncludeContentDiff *bool  `url:"include_content_diff,omitempty"`
}

// List returns all entries
// https://github.com/feedbin/feedbin-api/blob/master/content/entries.md#get-v2entriesjson
func (s *EntriesService) List(opts *EntryOptions) ([]*Entry, *http.Response, error) {
	url := "entries.json"
	url, err := AddQueryParams(url, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, url, nil)
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

// ListByFeed returns entries for a specific feed
// https://github.com/feedbin/feedbin-api/blob/master/content/entries.md#get-v2feeds203entriesjson
func (s *EntriesService) ListByFeed(feedID int, opts *EntryOptions) ([]*Entry, *http.Response, error) {
	url := fmt.Sprintf("feeds/%d/entries.json", feedID)
	url, err := AddQueryParams(url, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, url, nil)
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

// Get returns a specific entry
// https://github.com/feedbin/feedbin-api/blob/master/content/entries.md#get-v2entries3648json
func (s *EntriesService) Get(id int) (*Entry, *http.Response, error) {
	url := fmt.Sprintf("entries/%d.json", id)

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

// GetByIDs returns entries with the specified IDs
func (s *EntriesService) GetByIDs(ids []int) ([]*Entry, *http.Response, error) {
	if len(ids) > 100 {
		return nil, nil, fmt.Errorf("maximum of 100 IDs allowed per request")
	}

	idsStr := make([]string, len(ids))
	for i, id := range ids {
		idsStr[i] = strconv.Itoa(id)
	}

	url := fmt.Sprintf("entries.json?ids=%s", strings.Join(idsStr, ","))

	req, err := s.client.NewRequest(http.MethodGet, url, nil)
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

// GetExtractedContent fetches the extracted full content of an entry
func (s *EntriesService) GetExtractedContent(extractedContentURL string) (*ExtractedContent, *http.Response, error) {
	// The API provides a full URL for extracted content, so we need to use it directly
	// rather than constructing a URL relative to the API base URL
	req, err := http.NewRequest(http.MethodGet, extractedContentURL, nil)
	if err != nil {
		return nil, nil, err
	}

	// Add authentication
	req.SetBasicAuth(s.client.username, s.client.password)
	req.Header.Set("User-Agent", s.client.UserAgent)

	content := new(ExtractedContent)
	resp, err := s.client.client.Do(req)
	if err != nil {
		return nil, resp, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, resp, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	err = CheckResponse(resp)
	if err != nil {
		return nil, resp, err
	}

	err = decodeJSON(resp, content)
	if err != nil {
		return nil, resp, err
	}

	return content, resp, nil
}

// Helper function to decode JSON response
func decodeJSON(resp *http.Response, v interface{}) error {
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(v)
}
