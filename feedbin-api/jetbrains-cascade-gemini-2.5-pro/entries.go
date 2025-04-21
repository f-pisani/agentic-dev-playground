package feedbin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// EntriesService handles communication with the entry related
// methods of the Feedbin API.
type EntriesService service

// ListEntriesOptions specifies the optional parameters to the EntriesService.List method.
type ListEntriesOptions struct {
	Page               *int       `url:"page,omitempty"`
	Since              *time.Time `url:"since,omitempty"`     // ISO 8601 format
	IDs                []int64    `url:"ids,omitempty,comma"` // Comma-separated, max 100
	Read               *bool      `url:"read,omitempty"`
	Starred            *bool      `url:"starred,omitempty"`
	PerPage            *int       `url:"per_page,omitempty"`
	Mode               *string    `url:"mode,omitempty"` // Only "extended" is valid
	IncludeOriginal    *bool      `url:"include_original,omitempty"`
	IncludeEnclosure   *bool      `url:"include_enclosure,omitempty"`
	IncludeContentDiff *bool      `url:"include_content_diff,omitempty"`
}

// List retrieves all entries for the authenticated user, optionally filtered.
// Feedbin API docs: https://github.com/feedbin/feedbin-api/blob/master/content/entries.md#get-v2entriesjson
func (s *EntriesService) List(opt *ListEntriesOptions) ([]*Entry, *http.Response, error) {
	u, err := addOptions("entries.json", opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.newRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var entries []*Entry
	resp, err := s.client.do(req, &entries)
	if err != nil {
		return nil, resp, err
	}

	return entries, resp, nil
}

// ListFeedEntriesOptions specifies the optional parameters to the EntriesService.ListFeedEntries method.
// It excludes the 'ids' parameter which is not supported for feed-specific entries.
type ListFeedEntriesOptions struct {
	Page               *int       `url:"page,omitempty"`
	Since              *time.Time `url:"since,omitempty"` // ISO 8601 format
	Read               *bool      `url:"read,omitempty"`
	Starred            *bool      `url:"starred,omitempty"`
	PerPage            *int       `url:"per_page,omitempty"`
	Mode               *string    `url:"mode,omitempty"` // Only "extended" is valid
	IncludeOriginal    *bool      `url:"include_original,omitempty"`
	IncludeEnclosure   *bool      `url:"include_enclosure,omitempty"`
	IncludeContentDiff *bool      `url:"include_content_diff,omitempty"`
}

// ListFeedEntries retrieves all entries for a specific feed ID, optionally filtered.
// Feedbin API docs: https://github.com/feedbin/feedbin-api/blob/master/content/entries.md#get-v2feeds203entriesjson
func (s *EntriesService) ListFeedEntries(feedID int64, opt *ListFeedEntriesOptions) ([]*Entry, *http.Response, error) {
	path := fmt.Sprintf("feeds/%d/entries.json", feedID)
	u, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.newRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var entries []*Entry
	resp, err := s.client.do(req, &entries)
	if err != nil {
		return nil, resp, err
	}

	return entries, resp, nil
}

// Get retrieves a single entry by its ID.
// Feedbin API docs: https://github.com/feedbin/feedbin-api/blob/master/content/entries.md#get-v2entries3648json
func (s *EntriesService) Get(entryID int64) (*Entry, *http.Response, error) {
	path := fmt.Sprintf("entries/%d.json", entryID)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var entry Entry
	resp, err := s.client.do(req, &entry)
	if err != nil {
		return nil, resp, err
	}

	return &entry, resp, nil
}

// GetExtractedContent retrieves the extracted content for a given URL using the Feedbin Extract service.
// Note: This uses the client's username and password for authentication and signing.
// Feedbin API docs: https://github.com/feedbin/feedbin-api/blob/master/content/extract-full-content.md
func (s *EntriesService) GetExtractedContent(targetURL string) (*ExtractedContent, *http.Response, error) {
	// Input validation
	if targetURL == "" {
		return nil, nil, fmt.Errorf("targetURL cannot be empty")
	}
	if s.client.Username == "" || s.client.password == "" {
		return nil, nil, fmt.Errorf("username and password (used as secret key) must be set in the client for content extraction")
	}

	// Generate signature and base64 URL
	signature := generateHmacSignature(s.client.password, targetURL)
	base64URL := encodeBase64URLSafe(targetURL)

	// Construct the extract API URL
	relPath := fmt.Sprintf("%s/%s", s.client.Username, signature)
	rel, err := url.Parse(relPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse relative path for extract API: %w", err)
	}

	extractURL := s.client.ExtractURL.ResolveReference(rel)
	query := extractURL.Query()
	query.Set("base64_url", base64URL)
	extractURL.RawQuery = query.Encode()

	// Create the request (no basic auth needed for extract API, authentication is via signature)
	req, err := http.NewRequest(http.MethodGet, extractURL.String(), nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create extract API request: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")

	// Execute the request using the standard HTTP client directly
	// We don't use client.do because it adds Basic Auth which is not needed/wanted here
	// and error handling might differ slightly.
	resp, err := s.client.client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to execute extract API request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status (similar to client.checkResponse but without basic auth context)
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		// Attempt to reset body for potential future reads
		resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		apiErr := &APIError{
			StatusCode: resp.StatusCode,
			Response:   resp,
			Body:       string(bodyBytes),
			Message:    fmt.Sprintf("Extract API error (status %d): %s", resp.StatusCode, string(bodyBytes)),
		}
		return nil, resp, apiErr
	}

	// Decode the successful response
	var extractedContent ExtractedContent
	if err := json.NewDecoder(resp.Body).Decode(&extractedContent); err != nil {
		// Check for EOF which might mean empty but successful response
		if err == io.EOF {
			return &extractedContent, resp, nil // Return empty struct on EOF
		}
		return nil, resp, fmt.Errorf("failed to decode extract API response: %w", err)
	}

	return &extractedContent, resp, nil
}
