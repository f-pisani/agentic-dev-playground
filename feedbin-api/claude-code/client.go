// Package feedbin provides a client for the Feedbin API (V2).
package feedbin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	// DefaultBaseURL is the base URL for Feedbin API V2
	DefaultBaseURL = "https://api.feedbin.com/v2/"
	// UserAgent is the user agent used for API requests
	UserAgent = "GoFeedbinClient/1.0"
)

// Client represents a Feedbin API client
type Client struct {
	// HTTP client used to communicate with the API
	client *http.Client

	// Base URL for API requests
	BaseURL *url.URL

	// User credentials
	Username string
	Password string

	// Services used for communicating with different parts of the Feedbin API
	Authentication *AuthenticationService
	Subscriptions  *SubscriptionsService
	Entries        *EntriesService
	UnreadEntries  *UnreadEntriesService
	StarredEntries *StarredEntriesService
	Tags           *TagsService
	Taggings       *TaggingsService
	SavedSearches  *SavedSearchesService
}

// ClientOption allows customizing the Feedbin client
type ClientOption func(*Client)

// WithHTTPClient sets a custom HTTP client for the Feedbin client
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.client = httpClient
	}
}

// WithBaseURL sets a custom base URL for the Feedbin client
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		parsedURL, err := url.Parse(baseURL)
		if err == nil {
			c.BaseURL = parsedURL
		}
	}
}

// NewClient creates a new Feedbin API client.
func NewClient(username, password string, options ...ClientOption) *Client {
	baseURL, _ := url.Parse(DefaultBaseURL)

	c := &Client{
		client:   http.DefaultClient,
		BaseURL:  baseURL,
		Username: username,
		Password: password,
	}

	for _, option := range options {
		option(c)
	}

	// Initialize services
	c.Authentication = &AuthenticationService{client: c}
	c.Subscriptions = &SubscriptionsService{client: c}
	c.Entries = &EntriesService{client: c}
	c.UnreadEntries = &UnreadEntriesService{client: c}
	c.StarredEntries = &StarredEntriesService{client: c}
	c.Tags = &TagsService{client: c}
	c.Taggings = &TaggingsService{client: c}
	c.SavedSearches = &SavedSearchesService{client: c}

	return c
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

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

	req.SetBasicAuth(c.Username, c.Password)

	if body != nil {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", UserAgent)

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

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp, fmt.Errorf("api error: %s", resp.Status)
	}

	if v != nil && resp.StatusCode != http.StatusNoContent {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return resp, err
		}
	}

	return resp, nil
}

// PaginationInfo represents pagination information from Link headers
type PaginationInfo struct {
	NextPage     *int
	PreviousPage *int
	FirstPage    *int
	LastPage     *int
	TotalCount   int
}

// parseLink parses Link header
func parseLink(linkHeader string) *PaginationInfo {
	if linkHeader == "" {
		return nil
	}

	info := &PaginationInfo{}
	linkRegex := regexp.MustCompile(`<([^>]+)>;\s*rel="([^"]+)"`)
	pageRegex := regexp.MustCompile(`[&?]page=(\d+)`)

	for _, link := range strings.Split(linkHeader, ",") {
		matches := linkRegex.FindStringSubmatch(link)
		if len(matches) < 3 {
			continue
		}

		url := matches[1]
		rel := matches[2]

		pageMatches := pageRegex.FindStringSubmatch(url)
		if len(pageMatches) < 2 {
			continue
		}

		pageNum, _ := strconv.Atoi(pageMatches[1])

		switch rel {
		case "next":
			info.NextPage = &pageNum
		case "prev":
			info.PreviousPage = &pageNum
		case "first":
			info.FirstPage = &pageNum
		case "last":
			info.LastPage = &pageNum
		}
	}

	return info
}

// GetPagination extracts pagination information from response
func (c *Client) GetPagination(resp *http.Response) *PaginationInfo {
	info := parseLink(resp.Header.Get("Link"))
	if info != nil {
		countStr := resp.Header.Get("X-Feedbin-Record-Count")
		if countStr != "" {
			count, err := strconv.Atoi(countStr)
			if err == nil {
				info.TotalCount = count
			}
		}
	}
	return info
}

// ListOptions specifies the common parameters to various list API methods
type ListOptions struct {
	// For paginated results
	Page int `url:"page,omitempty"`

	// For time-based filtering
	Since time.Time `url:"since,omitempty"`
}

// AddQueryParams adds the parameters to the URL string
func addQueryParams(s string, options *ListOptions) (string, error) {
	if options == nil {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	q := u.Query()

	if options.Page > 0 {
		q.Add("page", strconv.Itoa(options.Page))
	}

	if !options.Since.IsZero() {
		q.Add("since", options.Since.Format(time.RFC3339Nano))
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}
