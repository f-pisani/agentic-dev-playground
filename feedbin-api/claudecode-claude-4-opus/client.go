package feedbin

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
	defaultBaseURL = "https://api.feedbin.com"
	apiVersion     = "v2"
	userAgent      = "feedbin-go-client/1.0"
)

// Client is the main Feedbin API client
type Client struct {
	BaseURL    string
	Email      string
	Password   string
	HTTPClient *http.Client
}

// NewClient creates a new Feedbin API client
func NewClient(email, password string) *Client {
	return &Client{
		BaseURL:    defaultBaseURL,
		Email:      email,
		Password:   password,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// NewClientWithHTTP creates a new Feedbin API client with a custom HTTP client
func NewClientWithHTTP(email, password string, httpClient *http.Client) *Client {
	return &Client{
		BaseURL:    defaultBaseURL,
		Email:      email,
		Password:   password,
		HTTPClient: httpClient,
	}
}

// request performs an HTTP request with authentication
func (c *Client) request(method, path string, body interface{}, query url.Values) (*http.Response, error) {
	// Build URL
	u, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, err
	}
	u.Path = fmt.Sprintf("/%s%s", apiVersion, path)
	if query != nil {
		u.RawQuery = query.Encode()
	}

	// Prepare body
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	// Create request
	req, err := http.NewRequest(method, u.String(), bodyReader)
	if err != nil {
		return nil, err
	}

	// Set headers
	req.SetBasicAuth(c.Email, c.Password)
	req.Header.Set("User-Agent", userAgent)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	// Execute request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// Check for errors
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, newAPIError(resp, string(bodyBytes))
	}

	return resp, nil
}

// get performs a GET request
func (c *Client) get(path string, query url.Values, result interface{}) error {
	resp, err := c.request(http.MethodGet, path, nil, query)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if result != nil && resp.StatusCode != http.StatusNoContent {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// post performs a POST request
func (c *Client) post(path string, body interface{}, result interface{}) error {
	resp, err := c.request(http.MethodPost, path, body, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if result != nil && resp.StatusCode != http.StatusNoContent {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// patch performs a PATCH request
func (c *Client) patch(path string, body interface{}, result interface{}) error {
	resp, err := c.request(http.MethodPatch, path, body, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if result != nil && resp.StatusCode != http.StatusNoContent {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// delete performs a DELETE request with optional body
func (c *Client) delete(path string, body interface{}) error {
	resp, err := c.request(http.MethodDelete, path, body, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// buildQuery builds query parameters from various options
func buildQuery(options ...QueryOption) url.Values {
	q := url.Values{}
	for _, opt := range options {
		opt(q)
	}
	return q
}

// QueryOption is a function that modifies query parameters
type QueryOption func(url.Values)

// WithPage sets the page number for pagination
func WithPage(page int) QueryOption {
	return func(q url.Values) {
		if page > 0 {
			q.Set("page", strconv.Itoa(page))
		}
	}
}

// WithPerPage sets the number of items per page
func WithPerPage(perPage int) QueryOption {
	return func(q url.Values) {
		if perPage > 0 {
			q.Set("per_page", strconv.Itoa(perPage))
		}
	}
}

// WithSince filters results to items created after the given time
func WithSince(since time.Time) QueryOption {
	return func(q url.Values) {
		if !since.IsZero() {
			q.Set("since", since.Format(time.RFC3339))
		}
	}
}

// WithRead filters entries by read status
func WithRead(read bool) QueryOption {
	return func(q url.Values) {
		q.Set("read", strconv.FormatBool(read))
	}
}

// WithStarred filters entries by starred status
func WithStarred(starred bool) QueryOption {
	return func(q url.Values) {
		q.Set("starred", strconv.FormatBool(starred))
	}
}

// WithIDs filters results to specific IDs
func WithIDs(ids []int) QueryOption {
	return func(q url.Values) {
		if len(ids) > 0 {
			idStrs := make([]string, len(ids))
			for i, id := range ids {
				idStrs[i] = strconv.Itoa(id)
			}
			q.Set("ids", strings.Join(idStrs, ","))
		}
	}
}

// WithMode sets the response mode (e.g., "extended")
func WithMode(mode string) QueryOption {
	return func(q url.Values) {
		if mode != "" {
			q.Set("mode", mode)
		}
	}
}

// WithIncludeOriginal includes original entry data
func WithIncludeOriginal(include bool) QueryOption {
	return func(q url.Values) {
		if include {
			q.Set("include_original", "true")
		}
	}
}

// WithIncludeEnclosure includes enclosure data
func WithIncludeEnclosure(include bool) QueryOption {
	return func(q url.Values) {
		if include {
			q.Set("include_enclosure", "true")
		}
	}
}

// WithIncludeContentDiff includes content diff for updated entries
func WithIncludeContentDiff(include bool) QueryOption {
	return func(q url.Values) {
		if include {
			q.Set("include_content_diff", "true")
		}
	}
}
