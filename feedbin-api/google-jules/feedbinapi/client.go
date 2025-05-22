package feedbinapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	defaultBaseURL = "https://api.feedbin.com/v2/"
	userAgent      = "Go-FeedbinAPI-Client/0.1 (google-jules)"
)

// Client manages communication with the Feedbin API.
type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
	username   string
	password   string

	// Services used for talking to different parts of the Feedbin API.
	// These will be initialized later.
	Entries *EntriesService
	Feeds *FeedsService
	// Subscriptions *SubscriptionsService
	// ... and so on for all resources
}

// NewClient creates a new Feedbin API client.
// A username and password are required for authentication.
func NewClient(username, password string, opts ...ClientOption) (*Client, error) {
	if username == "" {
		return nil, fmt.Errorf("username is required")
	}
	if password == "" {
		return nil, fmt.Errorf("password is required")
	}

	baseURL, _ := url.Parse(defaultBaseURL) // Should not fail for a valid constant

	c := &Client{
		baseURL:    baseURL,
		httpClient: http.DefaultClient,
		username:   username,
		password:   password,
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	// Initialize services here later
	c.Entries = &EntriesService{client: c}
	c.Feeds = &FeedsService{client: c}
	// ...

	return c, nil
}

// Response is a Feedbin API response. This wraps the standard http.Response.
type Response struct {
	*http.Response
	Pagination *PaginationInfo
}

// newRequest creates an API request.
// A relative URL pathStr can be provided in pathStr, in which case it is resolved relative to the baseURL of the Client.
// Relative URLs should always be specified without a preceding slash.
// If body is not nil, it will be JSON encoded and included in the request.
func (c *Client) newRequest(ctx context.Context, method, pathStr string, queryParams url.Values, body interface{}) (*http.Request, error) {
	relURL, err := url.Parse(pathStr)
	if err != nil {
		return nil, fmt.Errorf("parsing path: %w", err)
	}

	u := c.baseURL.ResolveReference(relURL)
	if queryParams != nil {
		u.RawQuery = queryParams.Encode()
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		if err := enc.Encode(body); err != nil {
			return nil, fmt.Errorf("encoding body: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("User-Agent", userAgent)
	if body != nil {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}
	req.Header.Set("Accept", "application/json")

	return req, nil
}

// do sends an API request and returns the API response.
// The API response is JSON decoded and stored in the value pointed to by v,
// or returned as an error if an API error has occurred.
// If v implements the io.Writer interface, the raw response body will be written to v,
// without attempting to decode it.
func (c *Client) do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-req.Context().Done():
			return nil, req.Context().Err()
		default:
		}
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	response := &Response{Response: resp}
	response.Pagination = extractPaginationInfo(resp.Header)

	// Check for API errors.
	if code := resp.StatusCode; code < 200 || code > 299 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return response, newAPIError(resp, bodyBytes) // Uses the new constructor from errors.go
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				return response, fmt.Errorf("copying response body to writer: %w", err)
			}
		} else {
			decErr := json.NewDecoder(resp.Body).Decode(v)
			if decErr == io.EOF {
				decErr = nil // ignore EOF errors caused by empty response body
			}
			if decErr != nil {
				err = fmt.Errorf("decoding response body: %w", decErr)
			}
		}
	}

	return response, err
}

// VerifyCredentials checks if the provided username and password are valid.
// It makes a request to the GET /v2/authentication.json endpoint.
func (c *Client) VerifyCredentials(ctx context.Context) error {
	req, err := c.newRequest(ctx, http.MethodGet, "authentication.json", nil, nil)
	if err != nil {
		return fmt.Errorf("creating VerifyCredentials request: %w", err)
	}

	resp, err := c.do(req, nil)
	if err != nil {
		apiErr, ok := err.(*APIError) // Check for the new APIError type
		if ok {
			if apiErr.IsStatus(http.StatusUnauthorized) { // Use IsStatus method
				return fmt.Errorf("authentication failed: invalid credentials (status %d)", apiErr.Response.StatusCode)
			}
			// Or, if you want to return a wrapped typed error:
			// if apiErr.IsStatus(http.StatusUnauthorized) {
			//    return newErrUnauthorized(apiErr)
			// }
		}
		return fmt.Errorf("VerifyCredentials API error: %w", err) // Wrap original error
	}

	if resp.StatusCode != http.StatusOK { // This check might be redundant if do() always returns APIError for non-2xx
		// but good as a safeguard or if do() could return other errors for 2xx.
		return fmt.Errorf("VerifyCredentials request failed: expected status %d, got %d (this might indicate an issue if APIError was not returned for a non-OK status)", http.StatusOK, resp.StatusCode)
	}

	return nil
}

// Helper function to add query parameters
func addQueryParams(baseURL string, params map[string]string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	q := u.Query()
	for key, value := range params {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

// Helper function to handle PATCH alternative via POST
func (c *Client) patchViaPOST(ctx context.Context, urlStr string, body interface{}, v interface{}) (*Response, error) {
    // Ensure the URL ends with /update.json or similar, or append if necessary
    if !strings.HasSuffix(urlStr, "/update.json") { // This check might need to be more robust
        urlStr += "/update.json" // Or handle this based on specific endpoint needs
    }
	req, err := c.newRequest(ctx, http.MethodPost, urlStr, nil, body)
	if err != nil {
		return nil, err
	}
	return c.do(req, v)
}
