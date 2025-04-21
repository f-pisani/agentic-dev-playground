package feedbin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://api.feedbin.com/v2/"
	userAgent      = "go-feedbin-client/0.1" // Consider making this configurable or dynamic
	contentType    = "application/json; charset=utf-8"
)

// Client manages communication with the Feedbin API.
type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
	username   string
	password   string
}

// NewClient creates a new Feedbin API client.
// If a nil httpClient is provided, http.DefaultClient will be used.
func NewClient(username, password string, httpClient *http.Client) (*Client, error) {
	if username == "" || password == "" {
		return nil, fmt.Errorf("username and password must be provided")
	}
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	parsedBaseURL, err := url.Parse(defaultBaseURL)
	if err != nil {
		// This should not happen with a hardcoded URL
		return nil, fmt.Errorf("error parsing base URL: %w", err)
	}

	return &Client{
		baseURL:    parsedBaseURL,
		httpClient: httpClient,
		username:   username,
		password:   password,
	}, nil
}

// Option represents a functional option for configuring API requests.
type Option func(*url.Values)

// --- Request Options ---

// ListSubscriptionsOptions specifies optional parameters for the ListSubscriptions method.
type ListSubscriptionsOptions struct {
	Since *time.Time
	Mode  *string // e.g., "extended"
}

// GetSubscriptionOptions specifies optional parameters for the GetSubscription method.
type GetSubscriptionOptions struct {
	Mode *string // e.g., "extended"
}

// ListEntriesOptions specifies optional parameters for the ListEntries and ListFeedEntries methods.
type ListEntriesOptions struct {
	Page               *int
	Since              *time.Time
	IDs                []int64 // Max 100
	Read               *bool
	Starred            *bool
	PerPage            *int
	Mode               *string // e.g., "extended"
	IncludeOriginal    *bool
	IncludeEnclosure   *bool
	IncludeContentDiff *bool
}

// GetEntryOptions specifies optional parameters for the GetEntry method.
type GetEntryOptions struct {
	Mode               *string // e.g., "extended"
	IncludeOriginal    *bool
	IncludeEnclosure   *bool
	IncludeContentDiff *bool
}

// GetSavedSearchResultsOptions specifies optional parameters for the GetSavedSearchResults method.
type GetSavedSearchResultsOptions struct {
	IncludeEntries *bool
	Page           *int
}

// --- Helper Functions ---

// addOptions adds query parameters to a URL based on the provided options struct.
func addOptions(baseURL string, opts interface{}) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return baseURL, err
	}
	q := u.Query()

	switch v := opts.(type) {
	case *ListSubscriptionsOptions:
		if v != nil {
			if v.Since != nil {
				q.Set("since", v.Since.Format(feedbinTimeFormat))
			}
			if v.Mode != nil {
				q.Set("mode", *v.Mode)
			}
		}
	case *GetSubscriptionOptions:
		if v != nil {
			if v.Mode != nil {
				q.Set("mode", *v.Mode)
			}
		}
	case *ListEntriesOptions:
		if v != nil {
			if v.Page != nil {
				q.Set("page", strconv.Itoa(*v.Page))
			}
			if v.Since != nil {
				q.Set("since", v.Since.Format(feedbinTimeFormat))
			}
			if len(v.IDs) > 0 {
				if len(v.IDs) > 100 {
					return "", fmt.Errorf("maximum of 100 IDs can be requested at once")
				}
				idsStr := make([]string, len(v.IDs))
				for i, id := range v.IDs {
					idsStr[i] = strconv.FormatInt(id, 10)
				}
				q.Set("ids", strings.Join(idsStr, ","))
			}
			if v.Read != nil {
				q.Set("read", strconv.FormatBool(*v.Read))
			}
			if v.Starred != nil {
				q.Set("starred", strconv.FormatBool(*v.Starred))
			}
			if v.PerPage != nil {
				q.Set("per_page", strconv.Itoa(*v.PerPage))
			}
			if v.Mode != nil {
				q.Set("mode", *v.Mode)
			}
			if v.IncludeOriginal != nil {
				q.Set("include_original", strconv.FormatBool(*v.IncludeOriginal))
			}
			if v.IncludeEnclosure != nil {
				q.Set("include_enclosure", strconv.FormatBool(*v.IncludeEnclosure))
			}
			if v.IncludeContentDiff != nil {
				q.Set("include_content_diff", strconv.FormatBool(*v.IncludeContentDiff))
			}
		}
	case *GetEntryOptions:
		if v != nil {
			if v.Mode != nil {
				q.Set("mode", *v.Mode)
			}
			if v.IncludeOriginal != nil {
				q.Set("include_original", strconv.FormatBool(*v.IncludeOriginal))
			}
			if v.IncludeEnclosure != nil {
				q.Set("include_enclosure", strconv.FormatBool(*v.IncludeEnclosure))
			}
			if v.IncludeContentDiff != nil {
				q.Set("include_content_diff", strconv.FormatBool(*v.IncludeContentDiff))
			}
		}
	case *GetSavedSearchResultsOptions:
		if v != nil {
			if v.IncludeEntries != nil {
				q.Set("include_entries", strconv.FormatBool(*v.IncludeEntries))
			}
			if v.Page != nil {
				q.Set("page", strconv.Itoa(*v.Page))
			}
		}
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}

// doRequest makes an HTTP request to the Feedbin API.
func (c *Client) doRequest(ctx context.Context, method, path string, params url.Values, body interface{}) (*http.Response, error) {
	rel := &url.URL{Path: path}
	u := c.baseURL.ResolveReference(rel)

	// Add query parameters
	q := u.Query()
	for k, v := range params {
		for _, vv := range v {
			q.Add(k, vv)
		}
	}
	u.RawQuery = q.Encode()

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		if err := enc.Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", contentType)
	if body != nil {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		return nil, err
	}

	return resp, nil
}

// handleResponse checks the API response for errors, and decodes the body
// into the provided value v.
func (c *Client) handleResponse(resp *http.Response, v interface{}) (*PaginationInfo, error) {
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, newAPIError(resp)
	}

	// Check for specific success codes where we don't expect a body
	if resp.StatusCode == http.StatusNoContent { // 204
		return nil, nil
	}

	// Check for 302 Found (Subscription exists) - body is the subscription
	// Check for 300 Multiple Choices (Multiple feeds found) - body is the choices
	// For these, we proceed to decode normally.

	paginationInfo := parsePaginationHeaders(resp)

	if v != nil {
		// If v implements io.Writer, write response body to it directly.
		// Otherwise, decode JSON into v.
		if w, ok := v.(io.Writer); ok {
			_, err := io.Copy(w, resp.Body)
			if err != nil {
				return paginationInfo, fmt.Errorf("error copying response body: %w", err)
			}
		} else {
			err := json.NewDecoder(resp.Body).Decode(v)
			if err != nil && err != io.EOF { // EOF is fine if the body is empty
				return paginationInfo, fmt.Errorf("error decoding response JSON: %w", err)
			}
		}
	}

	return paginationInfo, nil
}

// --- Pagination Parsing ---

var linkHeaderRegex = regexp.MustCompile(`<([^>]+)>;\s*rel="([^"]+)"`)

func parsePaginationHeaders(resp *http.Response) *PaginationInfo {
	info := &PaginationInfo{}

	// Parse Link header
	linkHeader := resp.Header.Get("Link")
	if linkHeader != "" {
		matches := linkHeaderRegex.FindAllStringSubmatch(linkHeader, -1)
		for _, match := range matches {
			if len(match) == 3 {
				url := match[1]
				rel := match[2]
				switch rel {
				case "next":
					info.NextPageURL = url
				case "last":
					info.LastPageURL = url
				case "first":
					info.FirstPageURL = url
				case "prev":
					info.PrevPageURL = url
				}
			}
		}
	}

	// Parse X-Feedbin-Record-Count header
	countHeader := resp.Header.Get("X-Feedbin-Record-Count")
	if countHeader != "" {
		count, err := strconv.Atoi(countHeader)
		if err == nil {
			info.TotalRecords = count
		}
	}

	return info
}

// --- API Methods ---

// VerifyCredentials checks if the provided credentials are valid.
func (c *Client) VerifyCredentials(ctx context.Context) (bool, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "authentication.json", nil, nil)
	if err != nil {
		// Check if the error is specifically a 401 Unauthorized
		if apiErr, ok := err.(*APIError); ok && apiErr.StatusCode == http.StatusUnauthorized {
			return false, nil // Valid response indicating invalid credentials
		}
		return false, err // Other network or unexpected errors
	}
	defer resp.Body.Close()

	// Any 2xx status code means success
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return true, nil
	}

	// This path shouldn't be reached due to handleResponse logic in typical methods,
	// but handle explicitly for VerifyCredentials.
	return false, newAPIError(resp)
}

// ListSubscriptions retrieves all subscriptions for the authenticated user.
func (c *Client) ListSubscriptions(ctx context.Context, opts *ListSubscriptionsOptions) ([]Subscription, error) {
	path, err := addOptions("subscriptions.json", opts)
	if err != nil {
		return nil, fmt.Errorf("error building query parameters: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}

	var subscriptions []Subscription
	_, err = c.handleResponse(resp, &subscriptions)
	if err != nil {
		return nil, err
	}

	return subscriptions, nil
}

// GetSubscription retrieves a specific subscription by its ID.
func (c *Client) GetSubscription(ctx context.Context, id int64, opts *GetSubscriptionOptions) (*Subscription, error) {
	path := fmt.Sprintf("subscriptions/%d.json", id)
	pathWithOptions, err := addOptions(path, opts)
	if err != nil {
		return nil, fmt.Errorf("error building query parameters: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodGet, pathWithOptions, nil, nil)
	if err != nil {
		return nil, err
	}

	var subscription Subscription
	_, err = c.handleResponse(resp, &subscription)
	if err != nil {
		return nil, err
	}

	return &subscription, nil
}

// CreateSubscription creates a new subscription to the specified feed URL.
// It returns the created or existing subscription.
// If the feed URL resolves to multiple feeds, it returns a list of choices and a nil subscription.
// An APIError with status 300 indicates multiple choices were found.
func (c *Client) CreateSubscription(ctx context.Context, feedURL string) (*Subscription, []FeedChoice, error) {
	body := map[string]string{"feed_url": feedURL}
	resp, err := c.doRequest(ctx, http.MethodPost, "subscriptions.json", nil, body)
	if err != nil {
		return nil, nil, err // Network error, etc.
	}
	defer resp.Body.Close()

	// Handle specific status codes
	switch resp.StatusCode {
	case http.StatusCreated, http.StatusFound: // 201, 302
		var subscription Subscription
		_, err = c.handleResponse(resp, &subscription)
		if err != nil {
			// Need to reset body reading potentially
			return nil, nil, fmt.Errorf("error decoding subscription response (status %d): %w", resp.StatusCode, err)
		}
		return &subscription, nil, nil
	case http.StatusMultipleChoices: // 300
		var choices []FeedChoice
		_, err = c.handleResponse(resp, &choices)
		if err != nil {
			// Need to reset body reading potentially
			return nil, nil, fmt.Errorf("error decoding feed choices response (status %d): %w", resp.StatusCode, err)
		}
		// Return the choices and a specific error type might be useful here,
		// but for now, just return the choices and nil error, matching the plan.
		// Caller should check if subscription is nil.
		return nil, choices, nil
	default:
		// Handle 404 Not Found, 415 Unsupported Media Type, etc.
		return nil, nil, newAPIError(resp)
	}
}

// UpdateSubscription updates the title of a subscription.
func (c *Client) UpdateSubscription(ctx context.Context, id int64, title string) (*Subscription, error) {
	path := fmt.Sprintf("subscriptions/%d.json", id)
	body := map[string]string{"title": title}

	resp, err := c.doRequest(ctx, http.MethodPatch, path, nil, body)
	if err != nil {
		return nil, err
	}

	var subscription Subscription
	_, err = c.handleResponse(resp, &subscription)
	if err != nil {
		return nil, err
	}

	return &subscription, nil
}

// DeleteSubscription deletes a subscription by its ID.
func (c *Client) DeleteSubscription(ctx context.Context, id int64) error {
	path := fmt.Sprintf("subscriptions/%d.json", id)
	resp, err := c.doRequest(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return err
	}

	// Expecting 204 No Content on success
	_, err = c.handleResponse(resp, nil)
	return err
}

// ListEntries retrieves all entries for the authenticated user, paginated.
func (c *Client) ListEntries(ctx context.Context, opts *ListEntriesOptions) ([]Entry, *PaginationInfo, error) {
	path, err := addOptions("entries.json", opts)
	if err != nil {
		return nil, nil, fmt.Errorf("error building query parameters: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var entries []Entry
	paginationInfo, err := c.handleResponse(resp, &entries)
	if err != nil {
		return nil, paginationInfo, err // Return partial pagination info on error if available
	}

	return entries, paginationInfo, nil
}

// ListFeedEntries retrieves entries for a specific feed, paginated.
func (c *Client) ListFeedEntries(ctx context.Context, feedID int64, opts *ListEntriesOptions) ([]Entry, *PaginationInfo, error) {
	basePath := fmt.Sprintf("feeds/%d/entries.json", feedID)
	path, err := addOptions(basePath, opts)
	if err != nil {
		return nil, nil, fmt.Errorf("error building query parameters: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var entries []Entry
	paginationInfo, err := c.handleResponse(resp, &entries)
	if err != nil {
		return nil, paginationInfo, err
	}

	return entries, paginationInfo, nil
}

// GetEntry retrieves a single entry by its ID.
func (c *Client) GetEntry(ctx context.Context, id int64, opts *GetEntryOptions) (*Entry, error) {
	basePath := fmt.Sprintf("entries/%d.json", id)
	path, err := addOptions(basePath, opts)
	if err != nil {
		return nil, fmt.Errorf("error building query parameters: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}

	var entry Entry
	_, err = c.handleResponse(resp, &entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

// --- Unread Entries ---

// ListUnreadEntries retrieves the IDs of all unread entries.
// To get the full entry details, use ListEntries with the returned IDs.
func (c *Client) ListUnreadEntries(ctx context.Context) ([]int64, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "unread_entries.json", nil, nil)
	if err != nil {
		return nil, err
	}

	var entryIDs []int64
	_, err = c.handleResponse(resp, &entryIDs)
	if err != nil {
		return nil, err
	}
	return entryIDs, nil
}

// MarkEntriesAsUnread marks the specified entry IDs as unread.
// Limit: 1000 IDs per request.
// Returns the list of IDs successfully marked as unread.
func (c *Client) MarkEntriesAsUnread(ctx context.Context, entryIDs []int64) ([]int64, error) {
	if len(entryIDs) > 1000 {
		return nil, fmt.Errorf("maximum of 1000 entry IDs allowed per request")
	}
	if len(entryIDs) == 0 {
		return []int64{}, nil // No-op
	}

	body := map[string][]int64{"unread_entries": entryIDs}
	resp, err := c.doRequest(ctx, http.MethodPost, "unread_entries.json", nil, body)
	if err != nil {
		return nil, err
	}

	var resultIDs []int64
	_, err = c.handleResponse(resp, &resultIDs)
	if err != nil {
		return nil, err
	}
	return resultIDs, nil
}

// MarkEntriesAsRead marks the specified entry IDs as read.
// Limit: 1000 IDs per request.
// Returns the list of IDs successfully marked as read.
func (c *Client) MarkEntriesAsRead(ctx context.Context, entryIDs []int64) ([]int64, error) {
	if len(entryIDs) > 1000 {
		return nil, fmt.Errorf("maximum of 1000 entry IDs allowed per request")
	}
	if len(entryIDs) == 0 {
		return []int64{}, nil // No-op
	}

	body := map[string][]int64{"unread_entries": entryIDs}
	resp, err := c.doRequest(ctx, http.MethodDelete, "unread_entries.json", nil, body)
	if err != nil {
		// Check for alternative POST endpoint if DELETE with body fails (though stdlib handles it)
		// Consider if a client option is needed to force POST alternative.
		return nil, err
	}

	var resultIDs []int64
	// Expect 200 OK with a body containing the IDs successfully marked
	_, err = c.handleResponse(resp, &resultIDs)
	if err != nil {
		return nil, err
	}
	return resultIDs, nil
}

// MarkEntriesAsReadAlt uses the alternative POST endpoint to mark entries as read.
// Useful for clients that have issues with DELETE requests containing a body.
func (c *Client) MarkEntriesAsReadAlt(ctx context.Context, entryIDs []int64) ([]int64, error) {
	if len(entryIDs) > 1000 {
		return nil, fmt.Errorf("maximum of 1000 entry IDs allowed per request")
	}
	if len(entryIDs) == 0 {
		return []int64{}, nil // No-op
	}

	body := map[string][]int64{"unread_entries": entryIDs}
	resp, err := c.doRequest(ctx, http.MethodPost, "unread_entries/delete.json", nil, body)
	if err != nil {
		return nil, err
	}

	var resultIDs []int64
	_, err = c.handleResponse(resp, &resultIDs)
	if err != nil {
		return nil, err
	}
	return resultIDs, nil
}

// --- Starred Entries ---

// ListStarredEntries retrieves the IDs of all starred entries.
// To get the full entry details, use ListEntries with the returned IDs.
func (c *Client) ListStarredEntries(ctx context.Context) ([]int64, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "starred_entries.json", nil, nil)
	if err != nil {
		return nil, err
	}

	var entryIDs []int64
	_, err = c.handleResponse(resp, &entryIDs)
	if err != nil {
		return nil, err
	}
	return entryIDs, nil
}

// StarEntries marks the specified entry IDs as starred.
// Limit: 1000 IDs per request.
// Returns the list of IDs successfully starred.
func (c *Client) StarEntries(ctx context.Context, entryIDs []int64) ([]int64, error) {
	if len(entryIDs) > 1000 {
		return nil, fmt.Errorf("maximum of 1000 entry IDs allowed per request")
	}
	if len(entryIDs) == 0 {
		return []int64{}, nil // No-op
	}

	body := map[string][]int64{"starred_entries": entryIDs}
	resp, err := c.doRequest(ctx, http.MethodPost, "starred_entries.json", nil, body)
	if err != nil {
		return nil, err
	}

	var resultIDs []int64
	_, err = c.handleResponse(resp, &resultIDs)
	if err != nil {
		return nil, err
	}
	return resultIDs, nil
}

// UnstarEntries removes the star from the specified entry IDs.
// Limit: 1000 IDs per request.
// Returns the list of IDs successfully unstarred.
func (c *Client) UnstarEntries(ctx context.Context, entryIDs []int64) ([]int64, error) {
	if len(entryIDs) > 1000 {
		return nil, fmt.Errorf("maximum of 1000 entry IDs allowed per request")
	}
	if len(entryIDs) == 0 {
		return []int64{}, nil // No-op
	}

	body := map[string][]int64{"starred_entries": entryIDs}
	resp, err := c.doRequest(ctx, http.MethodDelete, "starred_entries.json", nil, body)
	if err != nil {
		return nil, err
	}

	var resultIDs []int64
	// Expect 200 OK with a body containing the IDs successfully unstarred
	_, err = c.handleResponse(resp, &resultIDs)
	if err != nil {
		return nil, err
	}
	return resultIDs, nil
}

// UnstarEntriesAlt uses the alternative POST endpoint to unstar entries.
// Useful for clients that have issues with DELETE requests containing a body.
func (c *Client) UnstarEntriesAlt(ctx context.Context, entryIDs []int64) ([]int64, error) {
	if len(entryIDs) > 1000 {
		return nil, fmt.Errorf("maximum of 1000 entry IDs allowed per request")
	}
	if len(entryIDs) == 0 {
		return []int64{}, nil // No-op
	}

	body := map[string][]int64{"starred_entries": entryIDs}
	resp, err := c.doRequest(ctx, http.MethodPost, "starred_entries/delete.json", nil, body)
	if err != nil {
		return nil, err
	}

	var resultIDs []int64
	_, err = c.handleResponse(resp, &resultIDs)
	if err != nil {
		return nil, err
	}
	return resultIDs, nil
}

// --- Taggings ---

// ListTaggings retrieves all taggings for the authenticated user.
func (c *Client) ListTaggings(ctx context.Context) ([]Tagging, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "taggings.json", nil, nil)
	if err != nil {
		return nil, err
	}

	var taggings []Tagging
	_, err = c.handleResponse(resp, &taggings)
	if err != nil {
		return nil, err
	}
	return taggings, nil
}

// GetTagging retrieves a specific tagging by its ID.
func (c *Client) GetTagging(ctx context.Context, id int64) (*Tagging, error) {
	path := fmt.Sprintf("taggings/%d.json", id)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}

	var tagging Tagging
	_, err = c.handleResponse(resp, &tagging)
	if err != nil {
		return nil, err
	}
	return &tagging, nil
}

// CreateTagging assigns a tag name to a specific feed ID.
// Returns the created or existing tagging.
func (c *Client) CreateTagging(ctx context.Context, feedID int64, name string) (*Tagging, error) {
	body := map[string]interface{}{
		"feed_id": feedID,
		"name":    name,
	}
	resp, err := c.doRequest(ctx, http.MethodPost, "taggings.json", nil, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusCreated, http.StatusFound: // 201, 302
		var tagging Tagging
		// Need to re-decode body after checking status
		err = json.NewDecoder(resp.Body).Decode(&tagging)
		if err != nil {
			return nil, fmt.Errorf("error decoding tagging response (status %d): %w", resp.StatusCode, err)
		}
		return &tagging, nil
	default:
		return nil, newAPIError(resp)
	}
}

// DeleteTagging removes a specific tag assignment from a feed.
func (c *Client) DeleteTagging(ctx context.Context, id int64) error {
	path := fmt.Sprintf("taggings/%d.json", id)
	resp, err := c.doRequest(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return err
	}

	// Expecting 204 No Content on success
	_, err = c.handleResponse(resp, nil)
	return err
}

// --- Tags ---

// RenameTag changes the name of a tag across all feeds it's applied to.
// Returns the updated list of taggings affected by the rename.
func (c *Client) RenameTag(ctx context.Context, oldName, newName string) ([]Tagging, error) {
	body := map[string]string{
		"old_name": oldName,
		"new_name": newName,
	}
	resp, err := c.doRequest(ctx, http.MethodPost, "tags.json", nil, body)
	if err != nil {
		return nil, err
	}

	var taggings []Tagging
	_, err = c.handleResponse(resp, &taggings)
	if err != nil {
		return nil, err
	}
	return taggings, nil
}

// DeleteTag removes a tag entirely from all feeds it's applied to.
// Returns the list of remaining taggings (for the user, not just affected ones based on spec).
func (c *Client) DeleteTag(ctx context.Context, name string) ([]Tagging, error) {
	body := map[string]string{"name": name}
	resp, err := c.doRequest(ctx, http.MethodDelete, "tags.json", nil, body)
	if err != nil {
		return nil, err
	}

	var taggings []Tagging
	_, err = c.handleResponse(resp, &taggings)
	if err != nil {
		return nil, err
	}
	return taggings, nil
}

// --- Saved Searches ---

// ListSavedSearches retrieves all saved searches for the user.
func (c *Client) ListSavedSearches(ctx context.Context) ([]SavedSearch, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "saved_searches.json", nil, nil)
	if err != nil {
		return nil, err
	}

	var searches []SavedSearch
	_, err = c.handleResponse(resp, &searches)
	if err != nil {
		return nil, err
	}
	return searches, nil
}

// GetSavedSearchResults retrieves the results of a specific saved search.
// By default, it returns a list of entry IDs. Use options to include full entries and paginate.
func (c *Client) GetSavedSearchResults(ctx context.Context, id int64, opts *GetSavedSearchResultsOptions) ([]int64, []Entry, *PaginationInfo, error) {
	basePath := fmt.Sprintf("saved_searches/%d.json", id)
	path, err := addOptions(basePath, opts)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error building query parameters: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	includeEntries := opts != nil && opts.IncludeEntries != nil && *opts.IncludeEntries

	if includeEntries {
		var entries []Entry
		pagingInfo, err := c.handleResponse(resp, &entries)
		if err != nil {
			return nil, nil, pagingInfo, err
		}
		return nil, entries, pagingInfo, nil
	} else {
		var entryIDs []int64
		pagingInfo, err := c.handleResponse(resp, &entryIDs)
		if err != nil {
			return nil, nil, pagingInfo, err
		}
		return entryIDs, nil, pagingInfo, nil
	}
}

// CreateSavedSearch creates a new saved search.
func (c *Client) CreateSavedSearch(ctx context.Context, name, query string) (*SavedSearch, error) {
	body := map[string]string{
		"name":  name,
		"query": query,
	}
	resp, err := c.doRequest(ctx, http.MethodPost, "saved_searches.json", nil, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, newAPIError(resp)
	}

	// The spec doesn't explicitly state the response body on 201,
	// but it's reasonable to assume it returns the created object.
	var search SavedSearch
	err = json.NewDecoder(resp.Body).Decode(&search)
	if err != nil {
		return nil, fmt.Errorf("error decoding saved search response (status %d): %w", resp.StatusCode, err)
	}
	return &search, nil
}

// UpdateSavedSearch updates the name and/or query of an existing saved search.
func (c *Client) UpdateSavedSearch(ctx context.Context, id int64, name *string, query *string) (*SavedSearch, error) {
	if name == nil && query == nil {
		return nil, fmt.Errorf("at least one field (name or query) must be provided for update")
	}

	body := make(map[string]string)
	if name != nil {
		body["name"] = *name
	}
	if query != nil {
		body["query"] = *query
	}

	path := fmt.Sprintf("saved_searches/%d.json", id)
	resp, err := c.doRequest(ctx, http.MethodPatch, path, nil, body)
	if err != nil {
		// Consider implementing POST alternative: POST /v2/saved_searches/1/update.json
		return nil, err
	}

	var search SavedSearch
	_, err = c.handleResponse(resp, &search)
	if err != nil {
		return nil, err
	}
	return &search, nil
}

// DeleteSavedSearch deletes a saved search by its ID.
func (c *Client) DeleteSavedSearch(ctx context.Context, id int64) error {
	path := fmt.Sprintf("saved_searches/%d.json", id)
	resp, err := c.doRequest(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return err
	}

	// Expecting 204 No Content on success
	_, err = c.handleResponse(resp, nil)
	return err
}

// --- Recently Read Entries ---

// ListRecentlyReadEntries retrieves the IDs of recently read entries, ordered by recency.
func (c *Client) ListRecentlyReadEntries(ctx context.Context) ([]int64, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "recently_read_entries.json", nil, nil)
	if err != nil {
		return nil, err
	}

	var entryIDs []int64
	_, err = c.handleResponse(resp, &entryIDs)
	if err != nil {
		return nil, err
	}
	return entryIDs, nil
}

// CreateRecentlyReadEntries adds entry IDs to the recently read list.
// Returns the list of IDs successfully added.
func (c *Client) CreateRecentlyReadEntries(ctx context.Context, entryIDs []int64) ([]int64, error) {
	if len(entryIDs) == 0 {
		return []int64{}, nil // No-op
	}

	body := map[string][]int64{"recently_read_entries": entryIDs}
	resp, err := c.doRequest(ctx, http.MethodPost, "recently_read_entries.json", nil, body)
	if err != nil {
		return nil, err
	}

	var resultIDs []int64
	_, err = c.handleResponse(resp, &resultIDs)
	if err != nil {
		return nil, err
	}
	return resultIDs, nil
}

// --- Updated Entries ---

// ListUpdatedEntriesOptions specifies optional parameters for the ListUpdatedEntries method.
type ListUpdatedEntriesOptions struct {
	Since *time.Time
}

// ListUpdatedEntries retrieves the IDs of entries that have been updated since publication.
func (c *Client) ListUpdatedEntries(ctx context.Context, opts *ListUpdatedEntriesOptions) ([]int64, error) {
	path := "updated_entries.json"
	if opts != nil && opts.Since != nil {
		q := url.Values{}
		q.Set("since", opts.Since.Format(feedbinTimeFormat))
		path = path + "?" + q.Encode()
	}

	resp, err := c.doRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}

	var entryIDs []int64
	_, err = c.handleResponse(resp, &entryIDs)
	if err != nil {
		return nil, err
	}
	return entryIDs, nil
}

// DeleteUpdatedEntries removes entry IDs from the updated entries list (marks as seen).
// Returns the list of IDs successfully removed.
func (c *Client) DeleteUpdatedEntries(ctx context.Context, entryIDs []int64) ([]int64, error) {
	if len(entryIDs) == 0 {
		return []int64{}, nil // No-op
	}

	body := map[string][]int64{"updated_entries": entryIDs}
	resp, err := c.doRequest(ctx, http.MethodDelete, "updated_entries.json", nil, body)
	if err != nil {
		return nil, err
	}

	var resultIDs []int64
	_, err = c.handleResponse(resp, &resultIDs)
	if err != nil {
		return nil, err
	}
	return resultIDs, nil
}

// DeleteUpdatedEntriesAlt uses the alternative POST endpoint to remove updated entries.
func (c *Client) DeleteUpdatedEntriesAlt(ctx context.Context, entryIDs []int64) ([]int64, error) {
	if len(entryIDs) == 0 {
		return []int64{}, nil // No-op
	}

	body := map[string][]int64{"updated_entries": entryIDs}
	resp, err := c.doRequest(ctx, http.MethodPost, "updated_entries/delete.json", nil, body)
	if err != nil {
		return nil, err
	}

	var resultIDs []int64
	_, err = c.handleResponse(resp, &resultIDs)
	if err != nil {
		return nil, err
	}
	return resultIDs, nil
}

// --- Icons ---

// ListIcons retrieves the favicon URLs for all subscribed feeds.
func (c *Client) ListIcons(ctx context.Context) ([]Icon, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "icons.json", nil, nil)
	if err != nil {
		return nil, err
	}

	var icons []Icon
	_, err = c.handleResponse(resp, &icons)
	if err != nil {
		return nil, err
	}
	return icons, nil
}

// --- Imports ---

// CreateImport starts a new import job from an OPML file.
// opmlData should be an io.Reader containing the OPML XML content.
// Returns the initial status of the import job.
func (c *Client) CreateImport(ctx context.Context, opmlData io.Reader) (*Import, error) {
	rel := &url.URL{Path: "imports.json"}
	u := c.baseURL.ResolveReference(rel)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), opmlData)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", contentType)      // Expect JSON response
	req.Header.Set("Content-Type", "text/xml") // Sending XML

	resp, err := c.httpClient.Do(req)
	if err != nil {
		// Context error handling
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		return nil, err
	}

	var imp Import
	_, err = c.handleResponse(resp, &imp) // handleResponse checks for >= 400 errors
	if err != nil {
		return nil, err
	}
	return &imp, nil
}

// ListImports retrieves all import jobs for the user.
// Note: Does not include detailed import items.
func (c *Client) ListImports(ctx context.Context) ([]Import, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "imports.json", nil, nil)
	if err != nil {
		return nil, err
	}

	var imports []Import
	_, err = c.handleResponse(resp, &imports)
	if err != nil {
		return nil, err
	}
	return imports, nil
}

// GetImport retrieves the status of a specific import job, including item details.
func (c *Client) GetImport(ctx context.Context, id int64) (*Import, error) {
	path := fmt.Sprintf("imports/%d.json", id)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}

	var imp Import
	_, err = c.handleResponse(resp, &imp)
	if err != nil {
		return nil, err
	}
	return &imp, nil
}

// --- Pages ---

// CreatePage requests Feedbin to fetch the content from a given URL and create a new entry for it.
// An optional title can be provided, which is used if Feedbin cannot determine the title from the content.
// Returns the newly created Entry on success.
func (c *Client) CreatePage(ctx context.Context, pageURL string, title *string) (*Entry, error) {
	body := map[string]interface{}{
		"url": pageURL,
	}
	if title != nil {
		body["title"] = *title
	}

	resp, err := c.doRequest(ctx, http.MethodPost, "pages.json", nil, body)
	if err != nil {
		return nil, err
	}

	var entry Entry
	_, err = c.handleResponse(resp, &entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}
