// Package feedbin provides a Go client for the Feedbin API V2.
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
	BaseURL = "https://api.feedbin.com/v2"
	
	// UserAgent is the user agent used for API requests
	UserAgent = "Feedbin Go Client/1.0"
	
	// DefaultPerPage is the default number of items per page
	DefaultPerPage = 100
)

// Client represents a Feedbin API client
type Client struct {
	// HTTP client used to communicate with the API
	client *http.Client
	
	// Base URL for API requests
	baseURL *url.URL
	
	// User credentials for authentication
	username string
	password string
}

// NewClient returns a new Feedbin API client
func NewClient(username, password string) *Client {
	baseURL, _ := url.Parse(BaseURL)
	
	return &Client{
		client:   &http.Client{Timeout: 30 * time.Second},
		baseURL:  baseURL,
		username: username,
		password: password,
	}
}

// NewRequest creates an API request with authentication
func (c *Client) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(path)
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
	
	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", "application/json; charset=utf-8")
	
	if body != nil {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}
	
	return req, nil
}

// Do sends an API request and returns the API response
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	
	defer resp.Body.Close()
	
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Try to read error details if available
		var errorResponse ErrorResponse
		body, _ := io.ReadAll(resp.Body)
		
		// Log the response body for debugging
		errorMsg := fmt.Sprintf("API error: %s", resp.Status)
		if len(body) > 0 {
			errorMsg += fmt.Sprintf(" - Response: %s", string(body))
		}
		
		return resp, fmt.Errorf(errorMsg)
	}
	
	if v != nil && resp.StatusCode != http.StatusNoContent {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return resp, fmt.Errorf("JSON decode error: %v", err)
		}
	}
	
	return resp, nil
}

// GetPaginationLinks extracts pagination links from the Link header
func GetPaginationLinks(resp *http.Response) map[string]string {
	links := make(map[string]string)
	
	if linkHeader := resp.Header.Get("Link"); linkHeader != "" {
		for _, link := range strings.Split(linkHeader, ",") {
			segments := strings.Split(strings.TrimSpace(link), ";")
			if len(segments) < 2 {
				continue
			}
			
			// Extract URL from the first segment which is like: <https://api.feedbin.com/v2/feeds/1/entries.json?page=2>
			url := strings.Trim(segments[0], "<>")
			
			// Extract rel from the second segment which is like: rel="next"
			rel := strings.TrimSpace(segments[1])
			rel = strings.Trim(rel, "rel=")
			rel = strings.Trim(rel, "\"")
			
			links[rel] = url
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
		return 0
	}
	
	return count
}

// ParseFeedbinTime parses a time string in Feedbin's ISO 8601 format
func ParseFeedbinTime(timeStr string) (time.Time, error) {
	layouts := []string{
		"2006-01-02T15:04:05.999999Z",      // UTC format
		"2006-01-02T15:04:05.999999-07:00", // Timezone offset format
	}
	
	var t time.Time
	var err error
	
	for _, layout := range layouts {
		t, err = time.Parse(layout, timeStr)
		if err == nil {
			return t, nil
		}
	}
	
	return time.Time{}, fmt.Errorf("unable to parse time: %s", timeStr)
}

// FormatFeedbinTime formats a time.Time as a string in Feedbin's ISO 8601 format
func FormatFeedbinTime(t time.Time) string {
	return t.UTC().Format("2006-01-02T15:04:05.999999Z")
}
