// Package feedbin provides a client for the Feedbin API v2.
package feedbin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	// DefaultBaseURL is the default Feedbin API endpoint
	DefaultBaseURL = "https://api.feedbin.com/v2/"
	// UserAgent is the user agent string used in requests
	UserAgent = "Go-Feedbin-Client/1.0"
)

// Common errors
var (
	// ErrTooManyIDs is returned when too many IDs are provided in a request
	ErrTooManyIDs = errors.New("too many IDs (maximum 1000 per request)")
)

// Client manages communication with the Feedbin API
type Client struct {
	// HTTP client used to communicate with the API
	client *http.Client

	// Base URL for API requests
	BaseURL *url.URL

	// User agent used for API requests
	UserAgent string

	// Authentication credentials
	username string
	password string

	// Services used for communicating with different API endpoints
	Authentication *AuthenticationService
	Subscriptions  *SubscriptionsService
	Entries        *EntriesService
	Unread         *UnreadService
	Starred        *StarredService
	Tags           *TagsService
	Taggings       *TaggingsService
	SavedSearches  *SavedSearchesService
}

// NewClient creates a new Feedbin API client
func NewClient(username, password string) *Client {
	baseURL, _ := url.Parse(DefaultBaseURL)

	c := &Client{
		client:    http.DefaultClient,
		BaseURL:   baseURL,
		UserAgent: UserAgent,
		username:  username,
		password:  password,
	}

	// Initialize services
	c.Authentication = &AuthenticationService{client: c}
	c.Subscriptions = &SubscriptionsService{client: c}
	c.Entries = &EntriesService{client: c}
	c.Unread = &UnreadService{client: c}
	c.Starred = &StarredService{client: c}
	c.Tags = &TagsService{client: c}
	c.Taggings = &TaggingsService{client: c}
	c.SavedSearches = &SavedSearchesService{client: c}

	return c
}

// NewRequest creates an API request
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		err := json.NewEncoder(buf).Encode(body)
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
	req.Header.Set("User-Agent", c.UserAgent)
	req.SetBasicAuth(c.username, c.password)

	return req, nil
}

// Do sends an API request and returns the API response
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
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
			err = json.NewDecoder(resp.Body).Decode(v)
			if err == io.EOF {
				err = nil // ignore EOF, empty response body
			}
		}
	}

	return resp, err
}

// ErrorResponse represents an error response from the Feedbin API
type ErrorResponse struct {
	Response *http.Response
	Message  string `json:"message"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Message)
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

	return errorResponse
}

// Bool is a helper function that allocates a new bool value to store v and returns a pointer to it.
func Bool(v bool) *bool { return &v }

// Int is a helper function that allocates a new int value to store v and returns a pointer to it.
func Int(v int) *int { return &v }

// String is a helper function that allocates a new string value to store v and returns a pointer to it.
func String(v string) *string { return &v }

// Time is a helper function that allocates a new time.Time value to store v and returns a pointer to it.
func Time(v time.Time) *time.Time { return &v }
