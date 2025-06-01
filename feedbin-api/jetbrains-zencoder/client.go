// Package feedbin provides a client for the Feedbin API V2.
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

	// Username is the Feedbin account username.
	Username string

	// Password is the Feedbin account password.
	Password string
}

// NewClient creates a new Feedbin API client with the given credentials.
func NewClient(username, password string) *Client {
	return &Client{
		BaseURL:    BaseURL,
		HTTPClient: &http.Client{Timeout: DefaultTimeout},
		Username:   username,
		Password:   password,
	}
}

// NewClientWithOptions creates a new Feedbin API client with custom options.
func NewClientWithOptions(username, password string, baseURL string, timeout time.Duration) *Client {
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{Timeout: timeout},
		Username:   username,
		Password:   password,
	}
}

// doRequest performs an HTTP request and returns the response.
func (c *Client) doRequest(method, path string, body interface{}, query url.Values) (*http.Response, error) {
	// Construct the full URL
	reqURL, err := url.Parse(c.BaseURL + path)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// Add query parameters if provided
	if query != nil {
		reqURL.RawQuery = query.Encode()
	}

	// Create request body if provided
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	// Create the request
	req, err := http.NewRequest(method, reqURL.String(), reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set basic auth
	req.SetBasicAuth(c.Username, c.Password)

	// Set content type for POST, PUT, and PATCH requests
	if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}

	// Perform the request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

// parseResponse parses the response body into the given target.
func parseResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()

	// Check for non-2xx status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error: %s (status code: %d, body: %s)",
			resp.Status, resp.StatusCode, string(body))
	}

	// If no target is provided, just return
	if target == nil {
		return nil
	}

	// Parse the response body
	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	return nil
}
