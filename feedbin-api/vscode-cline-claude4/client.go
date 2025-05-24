// Package feedbin provides a comprehensive Go client for the Feedbin API v2.
//
// The client supports all Feedbin API endpoints with proper error handling,
// pagination support, HTTP caching, and follows Go best practices.
//
// Example usage:
//
//	client := feedbin.NewClient(&feedbin.Config{
//		Username: "your-email@example.com",
//		Password: "your-password",
//	})
//
//	ctx := context.Background()
//	subscriptions, err := client.GetSubscriptions(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
package feedbin

import (
	"bytes"
	"context"
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
	// DefaultBaseURL is the default Feedbin API base URL
	DefaultBaseURL = "https://api.feedbin.com/v2/"

	// DefaultUserAgent is the default user agent string
	DefaultUserAgent = "feedbin-go-client/1.0"

	// ContentTypeJSON is the JSON content type header value
	ContentTypeJSON = "application/json; charset=utf-8"

	// MaxBulkOperationSize is the maximum number of IDs allowed in bulk operations
	MaxBulkOperationSize = 1000

	// MaxEntriesPerRequest is the maximum number of entries that can be requested at once
	MaxEntriesPerRequest = 100
)

// Client represents a Feedbin API client
type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
	username   string
	password   string
	userAgent  string
	cache      CacheManager
}

// Config holds configuration options for the Feedbin client
type Config struct {
	// Username is the Feedbin account email
	Username string

	// Password is the Feedbin account password
	Password string

	// BaseURL is the API base URL (optional, defaults to DefaultBaseURL)
	BaseURL string

	// HTTPClient is a custom HTTP client (optional)
	HTTPClient *http.Client

	// UserAgent is a custom user agent string (optional)
	UserAgent string

	// EnableCache enables HTTP caching with ETag/Last-Modified headers
	EnableCache bool
}

// NewClient creates a new Feedbin API client with the given configuration
func NewClient(config *Config) *Client {
	if config == nil {
		panic("config cannot be nil")
	}

	if config.Username == "" || config.Password == "" {
		panic("username and password are required")
	}

	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		panic(fmt.Sprintf("invalid base URL: %v", err))
	}

	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	userAgent := config.UserAgent
	if userAgent == "" {
		userAgent = DefaultUserAgent
	}

	var cache CacheManager
	if config.EnableCache {
		cache = NewMemoryCache()
	}

	return &Client{
		baseURL:    parsedURL,
		httpClient: httpClient,
		username:   config.Username,
		password:   config.Password,
		userAgent:  userAgent,
		cache:      cache,
	}
}

// buildURL constructs a full URL from the given path and query parameters
func (c *Client) buildURL(path string, params url.Values) string {
	u := *c.baseURL
	u.Path = strings.TrimSuffix(u.Path, "/") + "/" + strings.TrimPrefix(path, "/")

	if params != nil {
		u.RawQuery = params.Encode()
	}

	return u.String()
}

// newRequest creates a new HTTP request with proper authentication and headers
func (c *Client) newRequest(ctx context.Context, method, path string, params url.Values, body interface{}) (*http.Request, error) {
	url := c.buildURL(path, params)

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set authentication
	req.SetBasicAuth(c.username, c.password)

	// Set headers
	req.Header.Set("User-Agent", c.userAgent)
	if body != nil {
		req.Header.Set("Content-Type", ContentTypeJSON)
	}

	return req, nil
}

// doRequest executes an HTTP request and handles the response
func (c *Client) doRequest(req *http.Request, result interface{}) (*http.Response, error) {
	// Check cache for GET requests
	if req.Method == http.MethodGet && c.cache != nil {
		if cached, found := c.cache.Get(req.URL.String()); found {
			if cached.ETag != "" {
				req.Header.Set("If-None-Match", cached.ETag)
			}
			if !cached.LastModified.IsZero() {
				req.Header.Set("If-Modified-Since", cached.LastModified.Format(http.TimeFormat))
			}
		}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Handle 304 Not Modified
	if resp.StatusCode == http.StatusNotModified && c.cache != nil {
		if cached, found := c.cache.Get(req.URL.String()); found {
			if result != nil {
				return resp, json.Unmarshal(cached.Data, result)
			}
			return resp, nil
		}
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, fmt.Errorf("failed to read response body: %w", err)
	}

	// Handle non-2xx status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp, c.handleErrorResponse(resp, body)
	}

	// Cache successful GET responses
	if req.Method == http.MethodGet && c.cache != nil && resp.StatusCode == http.StatusOK {
		cached := &CachedResponse{
			Data: body,
		}

		if etag := resp.Header.Get("ETag"); etag != "" {
			cached.ETag = etag
		}

		if lastModified := resp.Header.Get("Last-Modified"); lastModified != "" {
			if t, err := time.Parse(http.TimeFormat, lastModified); err == nil {
				cached.LastModified = t
			}
		}

		c.cache.Set(req.URL.String(), cached)
	}

	// Parse JSON response
	if result != nil && len(body) > 0 {
		if err := json.Unmarshal(body, result); err != nil {
			return resp, fmt.Errorf("failed to parse response JSON: %w", err)
		}
	}

	return resp, nil
}

// handleErrorResponse creates appropriate error types based on the HTTP response
func (c *Client) handleErrorResponse(resp *http.Response, body []byte) error {
	apiErr := &APIError{
		StatusCode: resp.StatusCode,
		Response:   resp,
	}

	// Try to parse error message from response body
	var errorResponse struct {
		Message string              `json:"message"`
		Error   string              `json:"error"`
		Errors  map[string][]string `json:"errors"`
	}

	if len(body) > 0 {
		if err := json.Unmarshal(body, &errorResponse); err == nil {
			if errorResponse.Message != "" {
				apiErr.Message = errorResponse.Message
			} else if errorResponse.Error != "" {
				apiErr.Message = errorResponse.Error
			} else if len(errorResponse.Errors) > 0 {
				// Convert validation errors
				var messages []string
				for field, errs := range errorResponse.Errors {
					for _, err := range errs {
						messages = append(messages, fmt.Sprintf("%s: %s", field, err))
					}
				}
				apiErr.Message = strings.Join(messages, ", ")
			}
		}
	}

	// Set default message if none found
	if apiErr.Message == "" {
		apiErr.Message = resp.Status
	}

	// Determine if error is retryable
	apiErr.Retryable = resp.StatusCode >= 500 || resp.StatusCode == 429

	return apiErr
}

// get performs a GET request
func (c *Client) get(ctx context.Context, path string, params url.Values, result interface{}) (*http.Response, error) {
	req, err := c.newRequest(ctx, http.MethodGet, path, params, nil)
	if err != nil {
		return nil, err
	}

	return c.doRequest(req, result)
}

// post performs a POST request
func (c *Client) post(ctx context.Context, path string, params url.Values, body interface{}, result interface{}) (*http.Response, error) {
	req, err := c.newRequest(ctx, http.MethodPost, path, params, body)
	if err != nil {
		return nil, err
	}

	return c.doRequest(req, result)
}

// patch performs a PATCH request
func (c *Client) patch(ctx context.Context, path string, params url.Values, body interface{}, result interface{}) (*http.Response, error) {
	req, err := c.newRequest(ctx, http.MethodPatch, path, params, body)
	if err != nil {
		return nil, err
	}

	return c.doRequest(req, result)
}

// delete performs a DELETE request
func (c *Client) delete(ctx context.Context, path string, params url.Values, body interface{}) (*http.Response, error) {
	req, err := c.newRequest(ctx, http.MethodDelete, path, params, body)
	if err != nil {
		return nil, err
	}

	return c.doRequest(req, nil)
}

// parsePagination extracts pagination information from response headers
func (c *Client) parsePagination(resp *http.Response) *PaginationInfo {
	pagination := &PaginationInfo{}

	// Parse X-Feedbin-Record-Count header
	if countStr := resp.Header.Get("X-Feedbin-Record-Count"); countStr != "" {
		if count, err := strconv.Atoi(countStr); err == nil {
			pagination.TotalCount = count
		}
	}

	// Parse Link header
	if linkHeader := resp.Header.Get("Link"); linkHeader != "" {
		links := parseLinkHeader(linkHeader)
		pagination.NextURL = links["next"]
		pagination.PrevURL = links["prev"]
		pagination.FirstURL = links["first"]
		pagination.LastURL = links["last"]
	}

	return pagination
}

// parseLinkHeader parses the Link header and returns a map of rel -> URL
func parseLinkHeader(header string) map[string]string {
	links := make(map[string]string)

	// Split by comma to get individual links
	parts := strings.Split(header, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)

		// Extract URL and rel
		if idx := strings.Index(part, ">; rel=\""); idx > 0 {
			url := strings.TrimPrefix(part[:idx], "<")
			relStart := idx + 8 // len(">; rel=\"")
			relEnd := strings.Index(part[relStart:], "\"")
			if relEnd > 0 {
				rel := part[relStart : relStart+relEnd]
				links[rel] = url
			}
		}
	}

	return links
}
