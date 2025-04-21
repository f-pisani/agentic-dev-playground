// Package feedbin provides a Go client for the Feedbin API v2.
package feedbin

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	// BaseURL is the base URL for the Feedbin API v2
	BaseURL = "https://api.feedbin.com/v2"
	// UserAgent is the user agent used for API requests
	UserAgent = "Feedbin Go Client"
	// ContentType is the content type used for API requests
	ContentType = "application/json; charset=utf-8"
)

// Client represents a Feedbin API client
type Client struct {
	// BaseURL is the base URL for API requests
	BaseURL string

	// HTTPClient is the HTTP client used for API requests
	HTTPClient *http.Client

	// Authentication credentials
	email    string
	password string

	// Services for different API endpoints
	Subscriptions       *SubscriptionsService
	Entries             *EntriesService
	UnreadEntries       *UnreadEntriesService
	StarredEntries      *StarredEntriesService
	Taggings            *TaggingsService
	Tags                *TagsService
	SavedSearches       *SavedSearchesService
	RecentlyReadEntries *RecentlyReadEntriesService
	UpdatedEntries      *UpdatedEntriesService
	Icons               *IconsService
	Imports             *ImportsService
	Pages               *PagesService
}

// NewClient creates a new Feedbin API client
func NewClient(email, password string) *Client {
	c := &Client{
		BaseURL:    BaseURL,
		HTTPClient: &http.Client{},
		email:      email,
		password:   password,
	}

	// Initialize services
	c.Subscriptions = &SubscriptionsService{client: c}
	c.Entries = &EntriesService{client: c}
	c.UnreadEntries = &UnreadEntriesService{client: c}
	c.StarredEntries = &StarredEntriesService{client: c}
	c.Taggings = &TaggingsService{client: c}
	c.Tags = &TagsService{client: c}
	c.SavedSearches = &SavedSearchesService{client: c}
	c.RecentlyReadEntries = &RecentlyReadEntriesService{client: c}
	c.UpdatedEntries = &UpdatedEntriesService{client: c}
	c.Icons = &IconsService{client: c}
	c.Imports = &ImportsService{client: c}
	c.Pages = &PagesService{client: c}

	return c
}

// VerifyAuthentication verifies the authentication credentials
func (c *Client) VerifyAuthentication() error {
	_, err := c.NewRequest("GET", "/authentication.json", nil)
	return err
}

// NewRequest creates a new API request
func (c *Client) NewRequest(method, path string, body interface{}) (*http.Response, error) {
	// Ensure path starts with a slash
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	// Create the full URL
	url := c.BaseURL + path

	// Create the request body if provided
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, fmt.Errorf("error encoding request body: %v", err)
		}
	}

	// Create the request
	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("User-Agent", UserAgent)
	if body != nil {
		req.Header.Set("Content-Type", ContentType)
	}

	// Set basic auth
	req.SetBasicAuth(c.email, c.password)

	// Send the request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}

	// Check for non-2xx status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		return nil, NewAPIError(resp)
	}

	return resp, nil
}

// ParseResponse parses the response body into the provided interface
func (c *Client) ParseResponse(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()

	if v == nil {
		return nil
	}

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return fmt.Errorf("error parsing response: %v", err)
	}

	return nil
}

// AddQueryParams adds query parameters to the URL
func AddQueryParams(baseURL string, params map[string]string) string {
	if len(params) == 0 {
		return baseURL
	}

	u, _ := url.Parse(baseURL)
	q := u.Query()

	for key, value := range params {
		q.Add(key, value)
	}

	u.RawQuery = q.Encode()
	return u.String()
}

// PaginationParams represents common pagination parameters
type PaginationParams struct {
	Page    int
	PerPage int
	Since   time.Time
}

// ToQueryParams converts pagination parameters to query parameters
func (p *PaginationParams) ToQueryParams() map[string]string {
	params := make(map[string]string)

	if p.Page > 0 {
		params["page"] = strconv.Itoa(p.Page)
	}

	if p.PerPage > 0 {
		params["per_page"] = strconv.Itoa(p.PerPage)
	}

	if !p.Since.IsZero() {
		params["since"] = p.Since.Format(time.RFC3339)
	}

	return params
}
