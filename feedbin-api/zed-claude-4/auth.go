package feedbin

import (
	"context"
	"encoding/base64"
	"net/http"
)

// Credentials represents user authentication credentials
type Credentials struct {
	Email    string
	Password string
}

// BasicAuth returns the base64 encoded basic auth string
func (c *Credentials) BasicAuth() string {
	auth := c.Email + ":" + c.Password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// Validate checks if the credentials are valid
func (c *Credentials) Validate() error {
	if c.Email == "" {
		return &ValidationError{
			Field:   "email",
			Message: "email is required",
		}
	}
	if c.Password == "" {
		return &ValidationError{
			Field:   "password",
			Message: "password is required",
		}
	}
	return nil
}

// setBasicAuth sets the Authorization header for HTTP basic authentication
func setBasicAuth(req *http.Request, creds *Credentials) {
	req.Header.Set("Authorization", "Basic "+creds.BasicAuth())
}

// Authenticate verifies the user's credentials with the Feedbin API
func (c *Client) Authenticate(ctx context.Context) error {
	req, err := c.newRequest(ctx, "GET", "/authentication.json", nil)
	if err != nil {
		return &AuthenticationError{Underlying: err}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return &NetworkError{
			Operation:  "authenticate",
			URL:        req.URL.String(),
			Underlying: err,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return &AuthenticationError{
			Underlying: &APIError{
				StatusCode: resp.StatusCode,
				Status:     resp.Status,
				Message:    "invalid email or password",
			},
		}
	}

	if err := checkResponse(resp); err != nil {
		return &AuthenticationError{Underlying: err}
	}

	return nil
}

// AuthenticateContext is an alias for Authenticate for backwards compatibility
func (c *Client) AuthenticateContext(ctx context.Context) error {
	return c.Authenticate(ctx)
}

// Authenticate without context (uses background context)
func (c *Client) AuthenticateWithoutContext() error {
	return c.Authenticate(context.Background())
}
