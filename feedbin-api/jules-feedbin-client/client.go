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
	// BaseURL is the base URL for the Feedbin API
	BaseURL   = "https://api.feedbin.com/v2/" // Ensure trailing slash
	UserAgent = "Jules Feedbin Go Client/1.0"
	// DefaultTimeout is the default timeout for API requests.
	DefaultTimeout = 30 * time.Second
)

// Client represents a Feedbin API client
type Client struct {
	client    *http.Client
	baseURL   *url.URL
	username  string
	password  string
	UserAgent string
	Authentication *AuthenticationService
	Subscriptions  *SubscriptionsService
	Entries        *EntriesService
	UnreadEntries  *UnreadEntriesService
	StarredEntries      *StarredEntriesService
	Taggings            *TaggingsService
	Tags                *TagsService
	SavedSearches       *SavedSearchesService
	RecentlyReadEntries *RecentlyReadEntriesService
	UpdatedEntries         *UpdatedEntriesService
	Icons                  *IconsService
	Imports                *ImportsService
	Pages                  *PagesService
	// ... other services will be added here (e.g. Extract)
}

// NewClient returns a new Feedbin API client
func NewClient(username, password string) *Client {
	baseURL, _ := url.Parse(BaseURL)
	c := &Client{
		client:    &http.Client{Timeout: DefaultTimeout},
		baseURL:   baseURL,
		username:  username,
		password:  password,
		UserAgent: UserAgent,
	}
	c.initServices()
	return c
}

// initServices initializes all the API services for the client.
func (c *Client) initServices() {
	c.Authentication = NewAuthenticationService(c)
	c.Subscriptions = NewSubscriptionsService(c)
	c.Entries = NewEntriesService(c)
	c.UnreadEntries = NewUnreadEntriesService(c)
	c.StarredEntries = NewStarredEntriesService(c)
	c.Taggings = NewTaggingsService(c)
	c.Tags = NewTagsService(c)
	c.SavedSearches = NewSavedSearchesService(c)
	c.RecentlyReadEntries = NewRecentlyReadEntriesService(c)
	c.UpdatedEntries = NewUpdatedEntriesService(c)
	c.Icons = NewIconsService(c)
	c.Imports = NewImportsService(c)
	c.Pages = NewPagesService(c)
	// ... initialize other services
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
	c.UserAgent = userAgent
}

// SetTimeout sets the timeout for API requests.
func (c *Client) SetTimeout(timeout time.Duration) {
	c.client.Timeout = timeout
}

// NewRequest creates an API request. A relative URL can be provided in path,
// in which case it is resolved relative to the BaseURL of the Client.
// If specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	// Ensure path is relative and does not start with a slash if baseURL already ends with one.
	// If baseURL does not end with a slash, ensure path starts with one if it's not empty.
	if strings.HasSuffix(c.baseURL.Path, "/") && strings.HasPrefix(path, "/") {
		path = strings.TrimPrefix(path, "/")
	} else if !strings.HasSuffix(c.baseURL.Path, "/") && !strings.HasPrefix(path, "/") && path != "" {
		// This case might need refinement based on how paths are constructed.
		// For now, assume paths are relative to the BaseURL's path.
	}

	rel, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("parsing path: %w", err)
	}

	u := c.baseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, fmt.Errorf("encoding body: %w", err)
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Accept", "application/json")

	if body != nil {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}

	return req, nil
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
		// Attempt to unmarshal into a structured error if possible,
		// otherwise use the raw data.
		// For now, just use the raw string.
		errorResponse.Message = string(data)
	} else {
		errorResponse.Message = http.StatusText(r.StatusCode) // Use status text if body is unreadable
	}
	// It's good practice to close the body if you've read it,
	// but 'Do' method's defer will handle it. If reading here,
	// need to be careful about double-closing or not closing.
	// Since ReadAll consumes the body, it needs to be replaced if further reading is needed.
	// However, for an error, we typically just consume and report.
	r.Body = io.NopCloser(bytes.NewBuffer(data)) // Replace body so it can be read again if necessary

	return errorResponse
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred.
// If v is nil, the response body is not decoded.
// If v is an io.Writer, the raw response body will be written to v.
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if err := CheckResponse(resp); err != nil {
		return resp, err // Return response along with the error
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				return resp, fmt.Errorf("writing response to writer: %w", err)
			}
		} else {
			// If status is 204 No Content, don't try to decode.
			if resp.StatusCode == http.StatusNoContent {
				return resp, nil
			}
			err = json.NewDecoder(resp.Body).Decode(v)
			if err == io.EOF {
				// Treat EOF as a non-error, signifies empty body which is fine for some responses.
				// Or, it could mean the JSON was malformed if it's unexpected.
				// For now, let's consider it not an error if the status code was successful.
				return resp, nil // Or return a specific error if EOF is always unexpected.
			}
			if err != nil {
				return resp, fmt.Errorf("decoding response: %w", err)
			}
		}
	}

	return resp, nil
}

// PaginationLinks represents the pagination links in the Link header.
type PaginationLinks struct {
	First string
	Prev  string
	Next  string
	Last  string
}

// GetPaginationLinks extracts pagination links from the Link header
func GetPaginationLinks(resp *http.Response) *PaginationLinks {
	links := &PaginationLinks{}
	linkHeader := resp.Header.Get("Link")
	if linkHeader == "" {
		return links
	}

	parts := strings.Split(linkHeader, ",")
	for _, part := range parts {
		segments := strings.Split(strings.TrimSpace(part), ";")
		if len(segments) < 2 {
			continue
		}

		urlPart := strings.Trim(segments[0], "<>")
		relPart := strings.TrimSpace(segments[1])
		rel := strings.Trim(relPart, `rel=""`) // Trim rel=" and "

		switch rel {
		case "first":
			links.First = urlPart
		case "prev":
			links.Prev = urlPart
		case "next":
			links.Next = urlPart
		case "last":
			links.Last = urlPart
		}
	}
	return links
}

// GetTotalCount extracts the total record count from the X-Feedbin-Record-Count header
func GetTotalCount(resp *http.Response) int {
	countStr := resp.Header.Get("X-Feedbin-Record-Count")
	if countStr == "" {
		return 0
	}
	count, err := strconv.Atoi(countStr)
	if err != nil {
		return 0 // Or handle error more explicitly
	}
	return count
}

// ParseFeedbinTime parses a time string in Feedbin's ISO 8601 format
func ParseFeedbinTime(timeStr string) (time.Time, error) {
	// Feedbin uses ISO 8601, which can have variations.
	// Go's time.RFC3339Nano is a common variant.
	// Example from docs: 2013-02-19T15:33:38.449047Z
	// Example from docs: 2013-02-19T07:33:38.449047-08:00
	layouts := []string{
		"2006-01-02T15:04:05.999999999Z07:00", // RFC3339Nano
		"2006-01-02T15:04:05.999999999Z",      // UTC variant of RFC3339Nano
		"2006-01-02T15:04:05Z07:00",           // Without fractional seconds
		"2006-01-02T15:04:05Z",                // UTC without fractional seconds
		// Adding layouts based on observed client implementations
		"2006-01-02T15:04:05.999999Z",      // Used in aider-claude-3.7
		"2006-01-02T15:04:05.999999-07:00", // Used in aider-claude-3.7
	}

	var t time.Time
	var err error
	for _, layout := range layouts {
		t, err = time.Parse(layout, timeStr)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse time string '%s' with known layouts: %w", timeStr, err)
}

// FormatFeedbinTime formats a time.Time as a string in Feedbin's ISO 8601 UTC format
// with microsecond precision, as commonly expected.
func FormatFeedbinTime(t time.Time) string {
	// Format to "2006-01-02T15:04:05.999999Z"
	return t.UTC().Format("2006-01-02T15:04:05.999999Z")
}

// Bool is a helper routine that allocates a new bool value
// to store v and returns a pointer to it.
func Bool(v bool) *bool {
	p := new(bool)
	*p = v
	return p
}

// Int is a helper routine that allocates a new int value
// to store v and returns a pointer to it.
func Int(v int) *int {
	p := new(int)
	*p = v
	return p
}

// String is a helper routine that allocates a new string value
// to store v and returns a pointer to it.
func String(v string) *string {
	p := new(string)
	*p = v
	return p
}
