package feedbin

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultBaseURL = "https://api.feedbin.com/v2/"
	defaultTimeout = 10 * time.Second
)

type Client struct {
	httpClient    *http.Client
	baseURL       string
	authenticator Authenticator
	secret        []byte
}

type ClientOption func(*Client)

func NewClient(username, password string, options ...ClientOption) *Client {
	c := &Client{
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		baseURL:       defaultBaseURL,
		authenticator: &BasicAuth{Username: username, Password: password},
	}

	for _, opt := range options {
		opt(c)
	}

	return c
}

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

func WithContentExtractorSecret(secret string) ClientOption {
	return func(c *Client) {
		c.secret = []byte(secret)
	}
}

func (c *Client) GetContentExtractorURL(pageURL string) (string, error) {
	if c.secret == nil {
		return "", fmt.Errorf("content extractor secret not configured")
	}
	mac := hmac.New(sha1.New, c.secret)
	mac.Write([]byte(pageURL))
	signature := hex.EncodeToString(mac.Sum(nil))
	encodedURL := base64.StdEncoding.EncodeToString([]byte(pageURL))
	return fmt.Sprintf("/parser/feedbin/%s?base64_url=%s", signature, encodedURL), nil
}

func (c *Client) newRequest(method, path string, body interface{}, opts ...RequestOption) (*http.Request, error) {
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	for _, opt := range opts {
		opt(q)
	}
	u.RawQuery = q.Encode()

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

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	c.authenticator.Authenticate(req)

	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return resp, &APIError{
			StatusCode: resp.StatusCode,
			Body:       string(bodyBytes),
			Message:    http.StatusText(resp.StatusCode),
		}
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}

	return resp, err
}