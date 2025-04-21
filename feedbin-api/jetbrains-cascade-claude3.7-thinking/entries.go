package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// EntriesService handles communication with the entry related
// methods of the Feedbin API
type EntriesService struct {
	client *Client
}

// EntryListOptions specifies the optional parameters to the
// EntriesService.List method
type EntryListOptions struct {
	Page               int       // Page number
	Since              time.Time // Get entries created after this time
	Ids                []int64   // Specific entry IDs to retrieve (max 100)
	Read               *bool     // Filter by read status
	Starred            *bool     // Filter by starred status
	PerPage            int       // Number of entries per page
	Mode               string    // Extended mode for more metadata
	IncludeOriginal    bool      // Include original entry data if updated
	IncludeEnclosure   bool      // Include podcast/RSS enclosure data
	IncludeContentDiff bool      // Include a diff of changed content
}

// PaginationInfo holds information about current page and links to other pages
type PaginationInfo struct {
	FirstLink  string
	PrevLink   string
	NextLink   string
	LastLink   string
	TotalCount int
}

// List returns entries based on the given options
func (s *EntriesService) List(opts *EntryListOptions) ([]*Entry, *PaginationInfo, error) {
	u := "entries.json"
	u, err := addEntryOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var entries []*Entry
	resp, err := s.client.do(req, &entries)
	if err != nil {
		return nil, nil, err
	}

	// Process pagination
	paginationInfo := &PaginationInfo{
		TotalCount: extractRecordCount(resp),
	}

	// Parse Link header for pagination URLs
	if linkHeader := resp.Header.Get("Link"); linkHeader != "" {
		links := parseLinkHeader(linkHeader)
		paginationInfo.FirstLink = links["first"]
		paginationInfo.PrevLink = links["prev"]
		paginationInfo.NextLink = links["next"]
		paginationInfo.LastLink = links["last"]
	}

	return entries, paginationInfo, nil
}

// ListFeedEntries returns entries for a specific feed
func (s *EntriesService) ListFeedEntries(feedID int64, opts *EntryListOptions) ([]*Entry, *PaginationInfo, error) {
	u := fmt.Sprintf("feeds/%d/entries.json", feedID)
	u, err := addEntryOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var entries []*Entry
	resp, err := s.client.do(req, &entries)
	if err != nil {
		return nil, nil, err
	}

	// Process pagination
	paginationInfo := &PaginationInfo{
		TotalCount: extractRecordCount(resp),
	}

	// Parse Link header for pagination URLs
	if linkHeader := resp.Header.Get("Link"); linkHeader != "" {
		links := parseLinkHeader(linkHeader)
		paginationInfo.FirstLink = links["first"]
		paginationInfo.PrevLink = links["prev"]
		paginationInfo.NextLink = links["next"]
		paginationInfo.LastLink = links["last"]
	}

	return entries, paginationInfo, nil
}

// Get returns a single entry
func (s *EntriesService) Get(id int64) (*Entry, error) {
	u := fmt.Sprintf("entries/%d.json", id)

	req, err := s.client.newRequest("GET", u, nil)
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

// addEntryOptions adds the parameters in opts as URL query parameters to s.
// opts is a pointer to an EntryListOptions struct.
func addEntryOptions(s string, opts *EntryListOptions) (string, error) {
	if opts == nil {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	q := u.Query()

	if opts.Page > 0 {
		q.Set("page", strconv.Itoa(opts.Page))
	}

	if !opts.Since.IsZero() {
		q.Set("since", opts.Since.Format(time.RFC3339Nano))
	}

	if len(opts.Ids) > 0 {
		if len(opts.Ids) > 100 {
			return s, fmt.Errorf("maximum of 100 IDs can be requested at once")
		}
		var ids []string
		for _, id := range opts.Ids {
			ids = append(ids, strconv.FormatInt(id, 10))
		}
		q.Set("ids", strings.Join(ids, ","))
	}

	if opts.Read != nil {
		q.Set("read", strconv.FormatBool(*opts.Read))
	}

	if opts.Starred != nil {
		q.Set("starred", strconv.FormatBool(*opts.Starred))
	}

	if opts.PerPage > 0 {
		q.Set("per_page", strconv.Itoa(opts.PerPage))
	}

	if opts.Mode != "" {
		q.Set("mode", opts.Mode)
	}

	if opts.IncludeOriginal {
		q.Set("include_original", "true")
	}

	if opts.IncludeEnclosure {
		q.Set("include_enclosure", "true")
	}

	if opts.IncludeContentDiff {
		q.Set("include_content_diff", "true")
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}

// GetExtractedContent fetches the full content extraction for an entry
func (s *EntriesService) GetExtractedContent(extractedContentURL string) (*ExtractedContent, error) {
	// For this method, we'll make a direct request to the URL provided in the entry
	req, err := http.NewRequest("GET", extractedContentURL, nil)
	if err != nil {
		return nil, err
	}

	// Set the same auth headers
	req.SetBasicAuth(s.client.username, s.client.password)

	// Set user agent
	req.Header.Set("User-Agent", s.client.userAgent)

	resp, err := s.client.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch extracted content: %s", resp.Status)
	}

	// Parse the response
	var extractedContent *ExtractedContent
	if err := json.NewDecoder(resp.Body).Decode(&extractedContent); err != nil {
		return nil, err
	}

	return extractedContent, nil
}

// ExtractedContent represents the Mercury Parser extraction result
type ExtractedContent struct {
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Author        string    `json:"author"`
	DatePublished time.Time `json:"date_published"`
	LeadImageURL  string    `json:"lead_image_url"`
	Dek           string    `json:"dek"`
	NextPageURL   string    `json:"next_page_url"`
	URL           string    `json:"url"`
	Domain        string    `json:"domain"`
	Excerpt       string    `json:"excerpt"`
	WordCount     int       `json:"word_count"`
	Direction     string    `json:"direction"`
	TotalPages    int       `json:"total_pages"`
	RenderedPages int       `json:"rendered_pages"`
}
