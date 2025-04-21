package client

import (
	"encoding/json"
	"fmt"
	
)

const (
	baseURL = "https://api.feedbin.com/v2/"
)

type Client struct {
	HTTPClient *http.Client
	Username   string
	Password   string
}

func NewClient(username, password string) (*Client, error) {
	if username == "" || password == "" {
		return nil, fmt.Errorf("missing email or password")
	}
	return &Client{
		HTTPClient: &http.Client{},
		Username:   username,
		Password:   password,
	}, nil
}

func (c *Client) Do(req *http.Request, v interface{}) error {
	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API error: %s", resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(v)
}

func (c *Client) NewRequest(method, endpoint string, params map[string]string) (*http.Request, error) {
	url := baseURL + endpoint
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	return req, nil
}
