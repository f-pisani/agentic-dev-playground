// Package feedbin provides a client for the Feedbin REST API v2.
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
	// BaseURL is the base URL for the Feedbin API
	BaseURL = "https://api.feedbin.com/v2/"

	// UserAgent is the user agent used for requests
	UserAgent = "Go-Feedbin-Client/1.0"
)

// Client represents a Feedbin API client
type Client struct {
	// BaseURL is the base URL for API requests
	BaseURL *url.URL

	// HTTPClient is the HTTP client used for requests
	HTTPClient *http.Client

	// Username is the Feedbin username
	Username string

	// Password is the Feedbin password
	Password string

	// LastETag stores the last ETag received from the API
	LastETag string

	// LastModified stores the last Modified timestamp received from the API
	LastModified string
}

// NewClient creates a new Feedbin API client
func NewClient(username, password string) *Client {
	baseURL, _ := url.Parse(BaseURL)

	return &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
		Username:   username,
		Password:   password,
	}
}

// NewRequest creates a new HTTP request
func (c *Client) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	u, err := c.BaseURL.Parse(path)
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

	req.SetBasicAuth(c.Username, c.Password)

	if body != nil {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", UserAgent)

	// Add caching headers if available
	if c.LastETag != "" {
		req.Header.Set("If-None-Match", c.LastETag)
	}

	if c.LastModified != "" {
		req.Header.Set("If-Modified-Since", c.LastModified)
	}

	return req, nil
}

// Do sends an HTTP request and returns an HTTP response
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Store caching headers
	if etag := resp.Header.Get("ETag"); etag != "" {
		c.LastETag = etag
	}

	if lastMod := resp.Header.Get("Last-Modified"); lastMod != "" {
		c.LastModified = lastMod
	}

	err = CheckResponse(resp)
	if err != nil {
		return resp, err
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

// Error represents an error returned by the Feedbin API
type Error struct {
	// StatusCode is the HTTP status code
	StatusCode int

	// Message is the error message
	Message string
}

// Error returns the error message
func (e *Error) Error() string {
	return fmt.Sprintf("feedbin: %d - %s", e.StatusCode, e.Message)
}

// CheckResponse checks the API response for errors
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	errorMessage := fmt.Sprintf("API request failed with status code: %d", r.StatusCode)

	return &Error{
		StatusCode: r.StatusCode,
		Message:    errorMessage,
	}
}

// Bool returns a pointer to the bool value passed in
func Bool(v bool) *bool {
	return &v
}

// Int returns a pointer to the int value passed in
func Int(v int) *int {
	return &v
}

// String returns a pointer to the string value passed in
func String(v string) *string {
	return &v
}
