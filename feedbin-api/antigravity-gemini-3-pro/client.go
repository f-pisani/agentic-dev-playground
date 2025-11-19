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
	defaultBaseURL = "https://api.feedbin.com/v2/"
	userAgent      = "go-feedbin-client"
)

// Client manages communication with the Feedbin API.
type Client struct {
	// BaseURL is the base URL for the API.
	BaseURL *url.URL

	// User agent used when communicating with the Feedbin API.
	UserAgent string

	// HTTP client used to communicate with the API.
	client *http.Client

	// Authentication credentials
	username string
	password string
}

// NewClient returns a new Feedbin API client.
func NewClient(username, password string) *Client {
	baseURL, _ := url.Parse(defaultBaseURL)

	return &Client{
		BaseURL:   baseURL,
		UserAgent: userAgent,
		client:    http.DefaultClient,
		username:  username,
		password:  password,
	}
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
	}
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.username, c.password)
	if body != nil {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	return req, nil
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response := &Response{Response: resp}
	response.populatePageValues()

	err = CheckResponse(resp)
	if err != nil {
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			dec := json.NewDecoder(resp.Body)
			err = dec.Decode(v)
			if err == io.EOF {
				err = nil // ignore EOF errors caused by empty response body
			}
		}
	}

	return response, err
}

// CheckResponse checks the API response for errors, and returns them if
// present. A response is considered an error if it has a status code outside
// the 200 range.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	errorResponse := &APIError{Response: r}
	data, err := io.ReadAll(r.Body)
	if err == nil && data != nil {
		errorResponse.Message = string(data)
	}
	return errorResponse
}

// APIError represents an error returned by the Feedbin API.
type APIError struct {
	Response *http.Response
	Message  string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API request failed with status code %d: %s. Body: %s", e.Response.StatusCode, e.Response.Status, e.Message)
}

// Response is a Feedbin API response. This wraps the standard http.Response
// returned from Feedbin and provides convenient access to things like
// pagination links.
type Response struct {
	*http.Response

	// These fields are populated from the Link header.
	NextPage  int
	PrevPage  int
	FirstPage int
	LastPage  int
}

func (r *Response) populatePageValues() {
	if links, ok := r.Header["Link"]; ok && len(links) > 0 {
		for _, link := range strings.Split(links[0], ",") {
			segments := strings.Split(strings.TrimSpace(link), ";")
			if len(segments) < 2 {
				continue
			}
			// link part: <https://api.feedbin.com/v2/entries.json?page=2>
			urlStr := strings.Trim(segments[0], "<>")
			// rel part: rel="next"
			rel := strings.TrimSpace(segments[1])

			u, err := url.Parse(urlStr)
			if err != nil {
				continue
			}
			pageStr := u.Query().Get("page")
			if pageStr == "" {
				continue
			}
			page, err := strconv.Atoi(pageStr)
			if err != nil {
				continue
			}

			switch {
			case strings.Contains(rel, `"next"`):
				r.NextPage = page
			case strings.Contains(rel, `"prev"`):
				r.PrevPage = page
			case strings.Contains(rel, `"first"`):
				r.FirstPage = page
			case strings.Contains(rel, `"last"`):
				r.LastPage = page
			}
		}
	}
}

// PaginationOptions specifies the optional parameters for paginated requests.
type PaginationOptions struct {
	Page    int       `url:"page,omitempty"`
	PerPage int       `url:"per_page,omitempty"`
	Since   time.Time `url:"since,omitempty"`
	Mode    string    `url:"mode,omitempty"` // "extended"
	Starred *bool     `url:"starred,omitempty"`
	Read    *bool     `url:"read,omitempty"`
	IDs     []int64   `url:"ids,omitempty"` // Comma separated list of IDs
}

// addOptions adds the parameters in opt as URL query parameters to s.
func addOptions(s string, opt *PaginationOptions) (string, error) {
	if opt == nil {
		return s, nil
	}
	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}
	q := u.Query()

	if opt.Page > 0 {
		q.Set("page", strconv.Itoa(opt.Page))
	}
	if opt.PerPage > 0 {
		q.Set("per_page", strconv.Itoa(opt.PerPage))
	}
	if !opt.Since.IsZero() {
		q.Set("since", opt.Since.Format(time.RFC3339Nano))
	}
	if opt.Mode != "" {
		q.Set("mode", opt.Mode)
	}
	if opt.Starred != nil {
		q.Set("starred", strconv.FormatBool(*opt.Starred))
	}
	if opt.Read != nil {
		q.Set("read", strconv.FormatBool(*opt.Read))
	}
	if len(opt.IDs) > 0 {
		ids := make([]string, len(opt.IDs))
		for i, id := range opt.IDs {
			ids[i] = strconv.FormatInt(id, 10)
		}
		q.Set("ids", strings.Join(ids, ","))
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}
