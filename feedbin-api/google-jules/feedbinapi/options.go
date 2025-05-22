package feedbinapi

import (
	"fmt"
	"net/http"
	"net/url"
)

// ClientOption is a function that configures a Client.
type ClientOption func(*Client) error

// WithHTTPClient sets the HTTP client for the Feedbin API client.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) error {
		if httpClient == nil {
			return fmt.Errorf("http client cannot be nil")
		}
		c.httpClient = httpClient
		return nil
	}
}

// WithBaseURL sets a custom base URL for the Feedbin API client.
// The default is https://api.feedbin.com/v2/.
func WithBaseURL(baseURLStr string) ClientOption {
	return func(c *Client) error {
		if baseURLStr == "" {
			return fmt.Errorf("base URL cannot be empty")
		}
		parsedURL, err := url.Parse(baseURLStr)
		if err != nil {
			return fmt.Errorf("failed to parse base URL: %w", err)
		}
		if parsedURL.Scheme == "" || parsedURL.Host == "" {
			return fmt.Errorf("base URL must be an absolute URL")
		}
		c.baseURL = parsedURL
		return nil
	}
}
