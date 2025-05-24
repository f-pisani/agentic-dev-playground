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
	defaultBaseURL = "https://api.feedbin.com/v2"
	defaultTimeout = 30 * time.Second
	userAgent      = "feedbin-go-client/1.0"
	contentType    = "application/json; charset=utf-8"
)

// Client represents a Feedbin API client
type Client struct {
	baseURL     *url.URL
	httpClient  *http.Client
	credentials *Credentials
}

// NewClient creates a new Feedbin API client
func NewClient(email, password string) *Client {
	baseURL, _ := url.Parse(defaultBaseURL)

	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		credentials: &Credentials{
			Email:    email,
			Password: password,
		},
	}
}

// NewClientWithHTTPClient creates a new client with a custom HTTP client
func NewClientWithHTTPClient(email, password string, httpClient *http.Client) *Client {
	baseURL, _ := url.Parse(defaultBaseURL)

	return &Client{
		baseURL:    baseURL,
		httpClient: httpClient,
		credentials: &Credentials{
			Email:    email,
			Password: password,
		},
	}
}

// SetBaseURL sets a custom base URL for the API
func (c *Client) SetBaseURL(baseURL string) error {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return fmt.Errorf("invalid base URL: %v", err)
	}
	c.baseURL = parsedURL
	return nil
}

// SetTimeout sets the HTTP client timeout
func (c *Client) SetTimeout(timeout time.Duration) {
	c.httpClient.Timeout = timeout
}

// newRequest creates a new HTTP request with proper headers and authentication
func (c *Client) newRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	if err := c.credentials.Validate(); err != nil {
		return nil, err
	}

	u, err := c.baseURL.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL path: %v", err)
	}

	var buf io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %v", err)
		}
		buf = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set required headers
	req.Header.Set("User-Agent", userAgent)
	if body != nil {
		req.Header.Set("Content-Type", contentType)
	}

	// Set authentication
	setBasicAuth(req, c.credentials)

	return req, nil
}

// do executes an HTTP request and handles the response
func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &NetworkError{
			Operation:  req.Method,
			URL:        req.URL.String(),
			Underlying: err,
		}
	}

	if err := checkResponse(resp); err != nil {
		return resp, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
		if err != nil {
			resp.Body.Close()
			return resp, fmt.Errorf("failed to decode response: %v", err)
		}
	}

	return resp, nil
}

// buildQueryParams builds URL query parameters from options
func buildQueryParams(opts interface{}) url.Values {
	params := url.Values{}

	if opts == nil {
		return params
	}

	switch o := opts.(type) {
	case *SubscriptionOptions:
		if o.Since != nil {
			params.Set("since", o.Since.Format(time.RFC3339Nano))
		}
		if o.Mode != nil {
			params.Set("mode", *o.Mode)
		}
	case *EntryOptions:
		if o.Page != nil {
			params.Set("page", strconv.Itoa(*o.Page))
		}
		if o.Since != nil {
			params.Set("since", o.Since.Format(time.RFC3339Nano))
		}
		if o.IDs != nil && len(o.IDs) > 0 {
			ids := make([]string, len(o.IDs))
			for i, id := range o.IDs {
				ids[i] = strconv.Itoa(id)
			}
			params.Set("ids", strings.Join(ids, ","))
		}
		if o.Read != nil {
			params.Set("read", strconv.FormatBool(*o.Read))
		}
		if o.Starred != nil {
			params.Set("starred", strconv.FormatBool(*o.Starred))
		}
		if o.PerPage != nil {
			params.Set("per_page", strconv.Itoa(*o.PerPage))
		}
		if o.Mode != nil {
			params.Set("mode", *o.Mode)
		}
		if o.IncludeOriginal != nil {
			params.Set("include_original", strconv.FormatBool(*o.IncludeOriginal))
		}
		if o.IncludeEnclosure != nil {
			params.Set("include_enclosure", strconv.FormatBool(*o.IncludeEnclosure))
		}
		if o.IncludeContentDiff != nil {
			params.Set("include_content_diff", strconv.FormatBool(*o.IncludeContentDiff))
		}
	case *SavedSearchOptions:
		if o.IncludeEntries != nil {
			params.Set("include_entries", strconv.FormatBool(*o.IncludeEntries))
		}
		if o.Page != nil {
			params.Set("page", strconv.Itoa(*o.Page))
		}
	case *UpdatedEntriesOptions:
		if o.Since != nil {
			params.Set("since", o.Since.Format(time.RFC3339Nano))
		}
	}

	return params
}

// addQueryParams adds query parameters to a URL path
func addQueryParams(path string, params url.Values) string {
	if len(params) == 0 {
		return path
	}
	return path + "?" + params.Encode()
}

// makeRequest is a helper method that combines newRequest and do
func (c *Client) makeRequest(ctx context.Context, method, path string, body interface{}, result interface{}) (*http.Response, error) {
	req, err := c.newRequest(ctx, method, path, body)
	if err != nil {
		return nil, err
	}

	return c.do(req, result)
}

// GetFromURL makes a GET request to a full URL (used for pagination)
func (c *Client) GetFromURL(ctx context.Context, fullURL string, result interface{}) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("User-Agent", userAgent)
	setBasicAuth(req, c.credentials)

	return c.do(req, result)
}
