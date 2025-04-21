// Package feedbin provides a Go client for the Feedbin API V2.
package feedbin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// BaseURL is the base URL for the Feedbin API.
	BaseURL = "https://api.feedbin.com"

	// UserAgent is the user agent used for API requests.
	UserAgent = "Feedbin Go Client"

	// DefaultTimeout is the default timeout for API requests.
	DefaultTimeout = 30 * time.Second
)

// Client is a client for the Feedbin API.
type Client struct {
	// HTTP client used to communicate with the API.
	client *http.Client

	// Base URL for API requests.
	baseURL *url.URL

	// User credentials for Basic Auth.
	username string
	password string

	// User agent used when communicating with the API.
	userAgent string

	// API endpoints
	Authentication *AuthenticationService
	Subscriptions  *SubscriptionsService
	Entries        *EntriesService
	UnreadEntries  *UnreadEntriesService
	StarredEntries *StarredEntriesService
	Taggings       *TaggingsService
	Tags           *TagsService
	SavedSearches  *SavedSearchesService
	UpdatedEntries *UpdatedEntriesService
	Icons          *IconsService
	Imports        *ImportsService
	Pages          *PagesService
	Extract        *ExtractService
}

// NewClient returns a new Feedbin API client.
func NewClient(username, password string) *Client {
	baseURL, _ := url.Parse(BaseURL)

	c := &Client{
		client:    &http.Client{Timeout: DefaultTimeout},
		baseURL:   baseURL,
		username:  username,
		password:  password,
		userAgent: UserAgent,
	}

	// Initialize services
	c.Authentication = &AuthenticationService{client: c}
	c.Subscriptions = &SubscriptionsService{client: c}
	c.Entries = &EntriesService{client: c}
	c.UnreadEntries = &UnreadEntriesService{client: c}
	c.StarredEntries = &StarredEntriesService{client: c}
	c.Taggings = &TaggingsService{client: c}
	c.Tags = &TagsService{client: c}
	c.SavedSearches = &SavedSearchesService{client: c}
	c.UpdatedEntries = &UpdatedEntriesService{client: c}
	c.Icons = &IconsService{client: c}
	c.Imports = &ImportsService{client: c}
	c.Pages = &PagesService{client: c}
	c.Extract = &ExtractService{client: c, username: username}

	return c
}

// SetBaseURL sets the base URL for API requests to a custom endpoint.
func (c *Client) SetBaseURL(urlStr string) error {
	baseURL, err := url.Parse(urlStr)
	if err != nil {
		return err
	}
	c.baseURL = baseURL
	return nil
}

// SetUserAgent sets the user agent for API requests.
func (c *Client) SetUserAgent(userAgent string) {
	c.userAgent = userAgent
}

// SetTimeout sets the timeout for API requests.
func (c *Client) SetTimeout(timeout time.Duration) {
	c.client.Timeout = timeout
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// If specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.baseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	// Set basic auth
	req.SetBasicAuth(c.username, c.password)

	// Set headers
	req.Header.Set("User-Agent", c.userAgent)
	if body != nil {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}

	return req, nil
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred.
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
				err = nil // ignore EOF errors caused by empty response body
			}
		}
	}

	return resp, err
}

// ErrorResponse reports an error caused by an API request.
type ErrorResponse struct {
	Response *http.Response // HTTP response that caused this error
	Message  string         // error message
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Message)
}

// CheckResponse checks the API response for errors, and returns them if present.
// A response is considered an error if it has a status code outside the 200 range.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := io.ReadAll(r.Body)
	if err == nil && data != nil {
		errorResponse.Message = string(data)
	} else {
		errorResponse.Message = r.Status
	}

	return errorResponse
}

// PaginationLinks represents the pagination links in the Link header.
type PaginationLinks struct {
	First string
	Prev  string
	Next  string
	Last  string
}

// ParseLinkHeader parses the Link header and returns a PaginationLinks struct.
func ParseLinkHeader(linkHeader string) *PaginationLinks {
	links := &PaginationLinks{}

	if linkHeader == "" {
		return links
	}

	// Split the link header into individual links
	parts := strings.Split(linkHeader, ",")

	for _, part := range parts {
		// Extract the URL and rel values
		urlAndRel := strings.Split(part, ";")
		if len(urlAndRel) != 2 {
			continue
		}

		// Clean up the URL
		url := strings.TrimSpace(urlAndRel[0])
		url = strings.TrimPrefix(url, "<")
		url = strings.TrimSuffix(url, ">")

		// Extract the rel value
		rel := strings.TrimSpace(urlAndRel[1])
		rel = strings.TrimPrefix(rel, `rel="`)
		rel = strings.TrimSuffix(rel, `"`)

		// Set the appropriate link
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

// GetTotalRecordCount extracts the total record count from the X-Feedbin-Record-Count header.
func GetTotalRecordCount(resp *http.Response) int {
	countStr := resp.Header.Get("X-Feedbin-Record-Count")
	if countStr == "" {
		return 0
	}

	var count int
	fmt.Sscanf(countStr, "%d", &count)
	return count
}
