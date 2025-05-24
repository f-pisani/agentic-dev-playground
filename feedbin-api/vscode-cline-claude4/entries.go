package feedbin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// GetEntries retrieves entries for the authenticated user with optional filtering.
// Returns a slice of entries and pagination information.
func (c *Client) GetEntries(ctx context.Context, opts *EntryOptions) ([]Entry, *PaginationInfo, error) {
	params := c.buildEntryParams(opts)

	var entries []Entry
	resp, err := c.get(ctx, "entries.json", params, &entries)
	if err != nil {
		return nil, nil, err
	}

	pagination := c.parsePagination(resp)
	return entries, pagination, nil
}

// GetEntry retrieves a specific entry by ID.
func (c *Client) GetEntry(ctx context.Context, id int) (*Entry, error) {
	path := fmt.Sprintf("entries/%d.json", id)

	var entry Entry
	_, err := c.get(ctx, path, nil, &entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

// GetFeedEntries retrieves entries for a specific feed.
// Returns a slice of entries and pagination information.
func (c *Client) GetFeedEntries(ctx context.Context, feedID int, opts *EntryOptions) ([]Entry, *PaginationInfo, error) {
	path := fmt.Sprintf("feeds/%d/entries.json", feedID)
	params := c.buildEntryParams(opts)

	var entries []Entry
	resp, err := c.get(ctx, path, params, &entries)
	if err != nil {
		return nil, nil, err
	}

	pagination := c.parsePagination(resp)
	return entries, pagination, nil
}

// GetEntriesByIDs retrieves specific entries by their IDs.
// Maximum of 100 entries can be requested at once.
func (c *Client) GetEntriesByIDs(ctx context.Context, ids []int) ([]Entry, error) {
	if len(ids) == 0 {
		return []Entry{}, nil
	}

	if len(ids) > MaxEntriesPerRequest {
		return nil, fmt.Errorf("too many entry IDs: maximum %d allowed, got %d", MaxEntriesPerRequest, len(ids))
	}

	// Convert IDs to comma-separated string
	idStrings := make([]string, len(ids))
	for i, id := range ids {
		idStrings[i] = strconv.Itoa(id)
	}

	params := url.Values{}
	params.Set("ids", strings.Join(idStrings, ","))

	var entries []Entry
	_, err := c.get(ctx, "entries.json", params, &entries)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

// GetEntriesExtended retrieves entries with extended metadata.
// This is a convenience method that calls GetEntries with mode="extended".
func (c *Client) GetEntriesExtended(ctx context.Context, opts *EntryOptions) ([]Entry, *PaginationInfo, error) {
	if opts == nil {
		opts = &EntryOptions{}
	}
	opts.Mode = "extended"
	return c.GetEntries(ctx, opts)
}

// GetUnreadEntriesWithDetails retrieves unread entries with full entry details.
// This is a convenience method that calls GetEntries with read=false.
func (c *Client) GetUnreadEntriesWithDetails(ctx context.Context, opts *EntryOptions) ([]Entry, *PaginationInfo, error) {
	if opts == nil {
		opts = &EntryOptions{}
	}
	read := false
	opts.Read = &read
	return c.GetEntries(ctx, opts)
}

// GetStarredEntriesWithDetails retrieves starred entries with full entry details.
// This is a convenience method that calls GetEntries with starred=true.
func (c *Client) GetStarredEntriesWithDetails(ctx context.Context, opts *EntryOptions) ([]Entry, *PaginationInfo, error) {
	if opts == nil {
		opts = &EntryOptions{}
	}
	starred := true
	opts.Starred = &starred
	return c.GetEntries(ctx, opts)
}

// buildEntryParams constructs URL parameters from EntryOptions
func (c *Client) buildEntryParams(opts *EntryOptions) url.Values {
	params := url.Values{}

	if opts == nil {
		return params
	}

	if opts.Page != nil {
		params.Set("page", strconv.Itoa(*opts.Page))
	}

	if opts.Since != nil {
		params.Set("since", opts.Since.Format("2006-01-02T15:04:05.000000Z"))
	}

	if len(opts.IDs) > 0 {
		if len(opts.IDs) > MaxEntriesPerRequest {
			// Truncate to maximum allowed
			opts.IDs = opts.IDs[:MaxEntriesPerRequest]
		}
		idStrings := make([]string, len(opts.IDs))
		for i, id := range opts.IDs {
			idStrings[i] = strconv.Itoa(id)
		}
		params.Set("ids", strings.Join(idStrings, ","))
	}

	if opts.Read != nil {
		params.Set("read", strconv.FormatBool(*opts.Read))
	}

	if opts.Starred != nil {
		params.Set("starred", strconv.FormatBool(*opts.Starred))
	}

	if opts.PerPage != nil {
		params.Set("per_page", strconv.Itoa(*opts.PerPage))
	}

	if opts.Mode != "" {
		params.Set("mode", opts.Mode)
	}

	if opts.IncludeOriginal {
		params.Set("include_original", "true")
	}

	if opts.IncludeEnclosure {
		params.Set("include_enclosure", "true")
	}

	if opts.IncludeContentDiff {
		params.Set("include_content_diff", "true")
	}

	return params
}

// EntryIterator provides a way to iterate through all entries with automatic pagination
type EntryIterator struct {
	client  *Client
	opts    *EntryOptions
	nextURL string
	hasMore bool
}

// NewEntryIterator creates a new entry iterator
func (c *Client) NewEntryIterator(ctx context.Context, opts *EntryOptions) *EntryIterator {
	return &EntryIterator{
		client:  c,
		opts:    opts,
		hasMore: true,
	}
}

// Next retrieves the next page of entries
func (iter *EntryIterator) Next(ctx context.Context) ([]Entry, error) {
	if !iter.hasMore {
		return nil, nil
	}

	var entries []Entry
	var resp *http.Response
	var err error

	if iter.nextURL != "" {
		// Use the next URL from pagination
		req, reqErr := http.NewRequestWithContext(ctx, "GET", iter.nextURL, nil)
		if reqErr != nil {
			return nil, reqErr
		}
		req.SetBasicAuth(iter.client.username, iter.client.password)
		req.Header.Set("User-Agent", iter.client.userAgent)

		resp, err = iter.client.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			body, _ := io.ReadAll(resp.Body)
			return nil, iter.client.handleErrorResponse(resp, body)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(body, &entries); err != nil {
			return nil, err
		}
	} else {
		// First request
		entries, pagination, err := iter.client.GetEntries(ctx, iter.opts)
		if err != nil {
			return nil, err
		}

		iter.nextURL = pagination.NextURL
		iter.hasMore = iter.nextURL != ""
		return entries, nil
	}

	// Parse pagination from response
	pagination := iter.client.parsePagination(resp)
	iter.nextURL = pagination.NextURL
	iter.hasMore = iter.nextURL != ""

	return entries, nil
}

// HasMore returns true if there are more entries to fetch
func (iter *EntryIterator) HasMore() bool {
	return iter.hasMore
}
