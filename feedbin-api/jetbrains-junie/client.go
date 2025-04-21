// Package feedbin provides a Go client for the Feedbin API v2.
package feedbin

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	// BaseURL is the base URL for the Feedbin API.
	BaseURL = "https://api.feedbin.com/v2"

	// DefaultTimeout is the default timeout for API requests.
	DefaultTimeout = 30 * time.Second
)

// Client represents a Feedbin API client.
type Client struct {
	// BaseURL is the base URL for API requests.
	BaseURL string

	// HTTPClient is the HTTP client used to make requests.
	HTTPClient *http.Client

	// Username is the Feedbin username.
	Username string

	// Password is the Feedbin password.
	Password string

	// LastETag stores the last ETag received from the API.
	LastETag string

	// LastModified stores the last modified timestamp received from the API.
	LastModified string

	// Services
	Subscriptions *SubscriptionService
	Entries       *EntryService
	Unread        *UnreadService
	Starred       *StarredService
	Taggings      *TaggingService
	Tags          *TagService
	SavedSearches *SavedSearchService
	RecentlyRead  *RecentlyReadService
	Updated       *UpdatedService
	Icons         *IconService
	Imports       *ImportService
	Pages         *PageService
}

// NewClient creates a new Feedbin API client with the given credentials.
func NewClient(username, password string) *Client {
	c := &Client{
		BaseURL:    BaseURL,
		HTTPClient: &http.Client{Timeout: DefaultTimeout},
		Username:   username,
		Password:   password,
	}

	// Initialize services
	c.Subscriptions = &SubscriptionService{client: c}
	c.Entries = &EntryService{client: c}
	c.Unread = &UnreadService{client: c}
	c.Starred = &StarredService{client: c}
	c.Taggings = &TaggingService{client: c}
	c.Tags = &TagService{client: c}
	c.SavedSearches = &SavedSearchService{client: c}
	c.RecentlyRead = &RecentlyReadService{client: c}
	c.Updated = &UpdatedService{client: c}
	c.Icons = &IconService{client: c}
	c.Imports = &ImportService{client: c}
	c.Pages = &PageService{client: c}

	return c
}

// NewRequest creates a new HTTP request with the given method, path, and body.
func (c *Client) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	u, err := url.Parse(c.BaseURL + path)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	// Set basic auth
	req.SetBasicAuth(c.Username, c.Password)

	// Set headers
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}

	// Set caching headers if available
	if c.LastETag != "" {
		req.Header.Set("If-None-Match", c.LastETag)
	}
	if c.LastModified != "" {
		req.Header.Set("If-Modified-Since", c.LastModified)
	}

	return req, nil
}

// Do sends an HTTP request and returns an HTTP response.
// It handles error responses and updates caching headers.
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Update caching headers
	if etag := resp.Header.Get("ETag"); etag != "" {
		c.LastETag = etag
	}
	if lastMod := resp.Header.Get("Last-Modified"); lastMod != "" {
		c.LastModified = lastMod
	}

	// Handle error responses
	if resp.StatusCode >= 400 {
		apiErr := NewAPIError(resp)
		return resp, apiErr
	}

	// Handle 304 Not Modified
	if resp.StatusCode == http.StatusNotModified {
		return resp, nil
	}

	// Parse response body if a target was provided
	if v != nil && resp.StatusCode != http.StatusNoContent {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return resp, err
		}
	}

	return resp, nil
}

// CheckAuth verifies the client's credentials.
func (c *Client) CheckAuth() (bool, error) {
	req, err := c.NewRequest("GET", "/authentication.json", nil)
	if err != nil {
		return false, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}
