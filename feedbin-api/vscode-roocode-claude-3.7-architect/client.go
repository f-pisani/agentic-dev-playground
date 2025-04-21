// Package feedbin provides a Go client for the Feedbin REST API.
package feedbin

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	// DefaultBaseURL is the default Feedbin API base URL
	DefaultBaseURL = "https://api.feedbin.com/v2/"
	// DefaultTimeout is the default timeout for HTTP requests
	DefaultTimeout = 30 * time.Second
)

// Client represents a Feedbin API client
type Client struct {
	// BaseURL is the base URL for API requests
	BaseURL string

	// HTTPClient is the HTTP client used for making requests
	HTTPClient *http.Client

	// Authentication credentials
	Email    string
	Password string

	// API endpoints
	Authentication *AuthenticationService
	Subscriptions  *SubscriptionsService
	Entries        *EntriesService
	UnreadEntries  *UnreadEntriesService
	StarredEntries *StarredEntriesService
	Tags           *TagsService
	Taggings       *TaggingsService
	SavedSearches  *SavedSearchesService
}

// ClientOption is a function that configures a Client
type ClientOption func(*Client)

// NewClient creates a new Feedbin API client
func NewClient(email, password string, options ...ClientOption) *Client {
	c := &Client{
		BaseURL:    DefaultBaseURL,
		HTTPClient: &http.Client{Timeout: DefaultTimeout},
		Email:      email,
		Password:   password,
	}

	// Apply options
	for _, option := range options {
		option(c)
	}

	// Initialize services
	c.Authentication = &AuthenticationService{client: c}
	c.Subscriptions = &SubscriptionsService{client: c}
	c.Entries = &EntriesService{client: c}
	c.UnreadEntries = &UnreadEntriesService{client: c}
	c.StarredEntries = &StarredEntriesService{client: c}
	c.Tags = &TagsService{client: c}
	c.Taggings = &TaggingsService{client: c}
	c.SavedSearches = &SavedSearchesService{client: c}

	return c
}

// WithBaseURL sets a custom base URL
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.BaseURL = baseURL
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.HTTPClient = httpClient
	}
}

// newRequest creates a new HTTP request
func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	url := c.BaseURL + strings.TrimPrefix(path, "/")

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.Email, c.Password)

	if body != nil {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}

	return req, nil
}

// do sends an HTTP request and returns the response
func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check for non-2xx status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp, NewAPIError(resp)
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}

	return resp, err
}

// PaginationLinks represents pagination links from the API
type PaginationLinks struct {
	First string
	Prev  string
	Next  string
	Last  string
}

// parsePaginationLinks parses the Link header for pagination links
func parsePaginationLinks(header string) PaginationLinks {
	links := PaginationLinks{}

	if header == "" {
		return links
	}

	// Regular expression to parse Link header
	re := regexp.MustCompile(`<([^>]+)>;\s*rel="([^"]+)"`)
	matches := re.FindAllStringSubmatch(header, -1)

	for _, match := range matches {
		if len(match) != 3 {
			continue
		}

		url := match[1]
		rel := match[2]

		switch rel {
		case "first":
			links.First = url
		case "prev":
			links.Prev = url
		case "next":
			links.Next = url
		case "last":
			links.Last = url
		}
	}

	return links
}

// getTotalRecords extracts the total record count from the X-Feedbin-Record-Count header
func getTotalRecords(header string) (int, error) {
	if header == "" {
		return 0, nil
	}

	return strconv.Atoi(header)
}

// Bool returns a pointer to the given bool value
func Bool(v bool) *bool {
	return &v
}

// Int returns a pointer to the given int value
func Int(v int) *int {
	return &v
}

// String returns a pointer to the given string value
func String(v string) *string {
	return &v
}
