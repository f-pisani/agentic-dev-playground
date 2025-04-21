package client

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
)

// Client is the Feedbin API client.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Username   string
	Password   string
}

// NewClient creates a new Feedbin API client.
func NewClient(username, password string) *Client {
	return &Client{
		BaseURL:    "https://api.feedbin.com/v2",
		HTTPClient: &http.Client{},
		Username:   username,
		Password:   password,
	}
}

// newRequest creates an HTTP request with basic auth.
func (c *Client) newRequest(method, path string, body io.Reader) (*http.Request, error) {
	u, err := url.Parse(c.BaseURL + path)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}
	// Basic Auth header
	req.Header.Set("Authorization", "Basic "+basicAuth(c.Username, c.Password))
	if method == http.MethodPost || method == http.MethodPatch {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}
	return req, nil
}

func basicAuth(username, password string) string {
	creds := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(creds))
}

// do sends the request and decodes the response.
func (c *Client) do(req *http.Request, v interface{}) error {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("Feedbin API error: %s", resp.Status)
		return errors.New(resp.Status)
	}
	if v != nil {
		return json.NewDecoder(resp.Body).Decode(v)
	}
	return nil
}
