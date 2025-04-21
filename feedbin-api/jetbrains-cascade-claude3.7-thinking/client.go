package client

import (
	"bytes"
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
	// BaseURL is the base URL for the Feedbin API
	defaultBaseURL = "https://api.feedbin.com/v2/"
	userAgent      = "FeedbinGoClient/1.0"
)

// Client manages communication with the Feedbin API
type Client struct {
	// HTTP client used to communicate with the API
	client *http.Client

	// Base URL for API requests
	baseURL *url.URL

	// Authentication credentials
	username string
	password string

	// User agent for client
	userAgent string

	// Services used for communicating with different parts of the Feedbin API
	Subscriptions  *SubscriptionsService
	Entries        *EntriesService
	UnreadEntries  *UnreadEntriesService
	StarredEntries *StarredEntriesService
	Tags           *TagsService
	Taggings       *TaggingsService
}

// ClientOption is a function that configures a Client
type ClientOption func(*Client) error

// NewClient returns a new Feedbin API client
func NewClient(username, password string, options ...ClientOption) (*Client, error) {
	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{
		client:    http.DefaultClient,
		baseURL:   baseURL,
		username:  username,
		password:  password,
		userAgent: userAgent,
	}

	// Apply any custom options
	for _, option := range options {
		if err := option(c); err != nil {
			return nil, err
		}
	}

	// Initialize services
	c.Subscriptions = &SubscriptionsService{client: c}
	c.Entries = &EntriesService{client: c}
	c.UnreadEntries = &UnreadEntriesService{client: c}
	c.StarredEntries = &StarredEntriesService{client: c}
	c.Tags = &TagsService{client: c}
	c.Taggings = &TaggingsService{client: c}

	return c, nil
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) error {
		if httpClient == nil {
			return fmt.Errorf("HTTP client cannot be nil")
		}
		c.client = httpClient
		return nil
	}
}

// WithBaseURL sets a custom base URL
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		u, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.baseURL = u
		return nil
	}
}

// WithTimeout sets a custom timeout for the HTTP client
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) error {
		c.client.Timeout = timeout
		return nil
	}
}

// TestAuthentication tests if the authentication credentials are valid
func (c *Client) TestAuthentication() (bool, error) {
	req, err := c.newRequest("GET", "authentication.json", nil)
	if err != nil {
		return false, err
	}

	resp, err := c.do(req, nil)
	if err != nil {
		return false, err
	}

	return resp.StatusCode == http.StatusOK, nil
}

// newRequest creates an API request
func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	u, err := c.baseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	req.SetBasicAuth(c.username, c.password)

	return req, nil
}

// do sends an API request and returns the API response
func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = CheckResponse(resp)
	if err != nil {
		return resp, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			decErr := json.NewDecoder(resp.Body).Decode(v)
			if decErr == io.EOF {
				decErr = nil // ignore EOF errors caused by empty response body
			}
			if decErr != nil {
				err = decErr
			}
		}
	}

	return resp, err
}

// CheckResponse checks the API response for errors
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := io.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}

	switch r.StatusCode {
	case 401:
		return &AuthError{ErrorResponse: errorResponse}
	case 403:
		return &ForbiddenError{ErrorResponse: errorResponse}
	case 404:
		return &NotFoundError{ErrorResponse: errorResponse}
	case 429:
		return &RateLimitError{ErrorResponse: errorResponse}
	default:
		return errorResponse
	}
}

// ErrorResponse represents an error response from the Feedbin API
type ErrorResponse struct {
	Response *http.Response
	Message  string `json:"error,omitempty"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Message)
}

// AuthError occurs when there are authentication issues
type AuthError struct {
	*ErrorResponse
}

// ForbiddenError occurs when the user doesn't have permission to access a resource
type ForbiddenError struct {
	*ErrorResponse
}

// NotFoundError occurs when a resource doesn't exist
type NotFoundError struct {
	*ErrorResponse
}

// RateLimitError occurs when rate limit is exceeded
type RateLimitError struct {
	*ErrorResponse
}

// parseLinkHeader parses the Link header for pagination
func parseLinkHeader(linkHeader string) map[string]string {
	links := make(map[string]string)

	if linkHeader == "" {
		return links
	}

	// Split the header into individual link parts
	parts := strings.Split(linkHeader, ",")

	for _, part := range parts {
		// Split each part into URL and rel sections
		section := strings.Split(part, ";")
		if len(section) < 2 {
			continue
		}

		// Extract URL and remove < > brackets
		url := strings.TrimSpace(section[0])
		url = strings.Trim(url, "<>")

		// Extract rel value
		relParts := strings.Split(strings.TrimSpace(section[1]), "=")
		if len(relParts) < 2 {
			continue
		}

		rel := strings.Trim(relParts[1], `"`)
		links[rel] = url
	}

	return links
}

// extractRecordCount extracts the X-Feedbin-Record-Count header
func extractRecordCount(r *http.Response) int {
	countStr := r.Header.Get("X-Feedbin-Record-Count")
	if countStr == "" {
		return 0
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		return 0
	}

	return count
}
