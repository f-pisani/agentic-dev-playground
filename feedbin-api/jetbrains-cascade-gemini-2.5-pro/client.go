package feedbin

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://api.feedbin.com/v2/"
	extractBaseURL = "https://extract.feedbin.com/parser/"
	userAgent      = "cascade-feedbin-go-client/1.0"
)

// Client manages communication with the Feedbin API.
type Client struct {
	client     *http.Client // HTTP client used to communicate with the API.
	BaseURL    *url.URL     // Base URL for API requests.
	ExtractURL *url.URL     // Base URL for Extract API requests.
	Username   string       // Feedbin username for authentication.
	password   string       // Feedbin password for authentication.

	// Services used for talking to different parts of the Feedbin API.
	Authentication *AuthenticationService
	Entries        *EntriesService
	Feeds          *FeedsService
	Icons          *IconsService
	Imports        *ImportsService
	Pages          *PagesService
	RecentlyRead   *RecentlyReadService
	SavedSearches  *SavedSearchesService
	StarredEntries *StarredEntriesService
	Subscriptions  *SubscriptionsService
	Taggings       *TaggingsService
	Tags           *TagsService
	UnreadEntries  *UnreadEntriesService
	UpdatedEntries *UpdatedEntriesService
}

// NewClient returns a new Feedbin API client.
// If a nil httpClient is provided, http.DefaultClient will be used.
func NewClient(username, password string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(defaultBaseURL)
	extractURL, _ := url.Parse(extractBaseURL)

	c := &Client{
		client:     httpClient,
		BaseURL:    baseURL,
		ExtractURL: extractURL,
		Username:   username,
		password:   password,
	}

	// Initialize services
	c.Authentication = &AuthenticationService{client: c}
	c.Entries = &EntriesService{client: c}
	c.Feeds = &FeedsService{client: c}
	c.Icons = &IconsService{client: c}
	c.Imports = &ImportsService{client: c}
	c.Pages = &PagesService{client: c}
	c.RecentlyRead = &RecentlyReadService{client: c}
	c.SavedSearches = &SavedSearchesService{client: c}
	c.StarredEntries = &StarredEntriesService{client: c}
	c.Subscriptions = &SubscriptionsService{client: c}
	c.Taggings = &TaggingsService{client: c}
	c.Tags = &TagsService{client: c}
	c.UnreadEntries = &UnreadEntriesService{client: c}
	c.UpdatedEntries = &UpdatedEntriesService{client: c}

	return c
}

// newRequest creates an API request. A relative URL path can be provided in
// path, in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash.
func (c *Client) newRequest(method, path string, body io.Reader) (*http.Request, error) {
	rel, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.Username, c.password)
	req.Header.Set("User-Agent", userAgent)
	if body != nil {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}
	req.Header.Set("Accept", "application/json")

	return req, nil
}

// do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := checkResponse(resp); err != nil {
		// even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		return resp, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err == io.EOF {
				err = nil // ignore EOF errors caused by empty response body
			}
		}
	}

	return resp, err
}

// checkResponse checks the API response for errors, and returns them if
// present. A response is considered an error if it has a status code outside
// the 200 range. API error responses are expected to have either no body,
// or a JSON response body that maps to ErrorResponse. Any other response
// body will be silently ignored.
func checkResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		// Error reading body, create a basic error
		return &APIError{
			StatusCode: r.StatusCode,
			Message:    fmt.Sprintf("API error (status code %d), could not read response body: %v", r.StatusCode, err),
			Response:   r,
		}
	}

	// Attempt to reset the response body so it can be read again later if needed
	r.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))

	apiErr := &APIError{
		StatusCode: r.StatusCode,
		Body:       string(bodyBytes),
		Response:   r,
	}

	// Try to unmarshal the error response if body is not empty
	if len(bodyBytes) > 0 {
		// Feedbin errors sometimes are simple strings, sometimes JSON arrays/objects
		// We'll just store the raw body string for now.
		apiErr.Message = fmt.Sprintf("API error (status code %d): %s", r.StatusCode, string(bodyBytes))
	}

	return apiErr
}

// Bool is a helper routine that allocates a new bool value
// to store v and returns a pointer to it.
func Bool(v bool) *bool { return &v }

// Int is a helper routine that allocates a new int value
// to store v and returns a pointer to it.
func Int(v int) *int { return &v }

// String is a helper routine that allocates a new string value
// to store v and returns a pointer to it.
func String(v string) *string { return &v }

// Time is a helper routine that allocates a new time.Time value
// to store v and returns a pointer to it.
func Time(v time.Time) *time.Time { return &v }

// addOptions adds the parameters in opt as URL query parameters to s.
// opt must be a struct whose fields may contain "url" tags.
func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs := url.Values{}

	// Use reflection or a dedicated library if complex options are needed.
	// For now, handle specific option types manually.

	switch o := opt.(type) {
	case *ListEntriesOptions:
		if o != nil {
			if o.Page != nil {
				qs.Set("page", fmt.Sprintf("%d", *o.Page))
			}
			if o.Since != nil {
				qs.Set("since", o.Since.Format(time.RFC3339Nano))
			}
			if len(o.IDs) > 0 {
				qs.Set("ids", joinInts(o.IDs, ","))
			}
			if o.Read != nil {
				qs.Set("read", fmt.Sprintf("%t", *o.Read))
			}
			if o.Starred != nil {
				qs.Set("starred", fmt.Sprintf("%t", *o.Starred))
			}
			if o.PerPage != nil {
				qs.Set("per_page", fmt.Sprintf("%d", *o.PerPage))
			}
			if o.Mode != nil {
				qs.Set("mode", *o.Mode)
			}
			if o.IncludeOriginal != nil {
				qs.Set("include_original", fmt.Sprintf("%t", *o.IncludeOriginal))
			}
			if o.IncludeEnclosure != nil {
				qs.Set("include_enclosure", fmt.Sprintf("%t", *o.IncludeEnclosure))
			}
			if o.IncludeContentDiff != nil {
				qs.Set("include_content_diff", fmt.Sprintf("%t", *o.IncludeContentDiff))
			}
		}
	case *ListFeedEntriesOptions:
		if o != nil {
			if o.Page != nil {
				qs.Set("page", fmt.Sprintf("%d", *o.Page))
			}
			if o.Since != nil {
				qs.Set("since", o.Since.Format(time.RFC3339Nano))
			}
			if o.Read != nil {
				qs.Set("read", fmt.Sprintf("%t", *o.Read))
			}
			if o.Starred != nil {
				qs.Set("starred", fmt.Sprintf("%t", *o.Starred))
			}
			if o.PerPage != nil {
				qs.Set("per_page", fmt.Sprintf("%d", *o.PerPage))
			}
			if o.Mode != nil {
				qs.Set("mode", *o.Mode)
			}
			if o.IncludeOriginal != nil {
				qs.Set("include_original", fmt.Sprintf("%t", *o.IncludeOriginal))
			}
			if o.IncludeEnclosure != nil {
				qs.Set("include_enclosure", fmt.Sprintf("%t", *o.IncludeEnclosure))
			}
			if o.IncludeContentDiff != nil {
				qs.Set("include_content_diff", fmt.Sprintf("%t", *o.IncludeContentDiff))
			}
		}
	case *ListSubscriptionsOptions:
		if o != nil {
			if o.Mode != nil {
				qs.Set("mode", *o.Mode)
			}
		}
	case *GetSavedSearchEntriesOptions:
		if o != nil {
			if o.IncludeEntries != nil {
				qs.Set("include_entries", fmt.Sprintf("%t", *o.IncludeEntries))
			}
			if o.Page != nil {
				qs.Set("page", fmt.Sprintf("%d", *o.Page))
			}
		}
	default:
		// No options or unknown type
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

// joinInts converts a slice of ints to a comma-separated string.
func joinInts(ints []int64, sep string) string {
	strs := make([]string, len(ints))
	for i, v := range ints {
		strs[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(strs, sep)
}

// generateHmacSignature generates the HMAC-SHA1 signature for the extract URL.
func generateHmacSignature(secret, data string) string {
	h := hmac.New(sha1.New, []byte(secret))
	h.Write([]byte(data))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// encodeBase64URLSafe encodes a string using URL-safe Base64 encoding.
func encodeBase64URLSafe(data string) string {
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte(data))
}

// --- Services ---

type service struct {
	client *Client
}
