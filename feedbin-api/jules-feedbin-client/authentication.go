package feedbin

import (
	"net/http"
)

// AuthenticationService handles operations related to API authentication.
type AuthenticationService struct {
	client *Client
}

// NewAuthenticationService creates a new service for authentication related operations.
func NewAuthenticationService(client *Client) *AuthenticationService {
	return &AuthenticationService{client: client}
}

// Verify checks if the provided credentials are valid by making a simple API call
// that requires authentication, for example, trying to get subscriptions.
// Feedbin API docs state: "The Feedbin API uses HTTP Basic authentication".
// A successful request to any authenticated endpoint can serve as verification.
// We can use the root endpoint or a lightweight one like listing subscriptions with per_page=1.
// For now, let's assume a successful call to a protected resource means authentication is OK.
// The API doesn't have a dedicated "verify credentials" endpoint.
// We will try to fetch the first page of subscriptions.
func (s *AuthenticationService) Verify() (bool, *http.Response, error) {
	// Use a lightweight request to check authentication.
	// Requesting a single subscription or a specific known entry could also work.
	// For this example, we'll try to list subscriptions with a limit of 1.
	// The actual endpoint for this might be "/v2/subscriptions.json?per_page=1"
	// The NewRequest method in client.go will handle prepending /v2/
	req, err := s.client.NewRequest(http.MethodGet, "subscriptions.json?per_page=1", nil)
	if err != nil {
		return false, nil, err
	}

	// We don't need to decode the body, just check the response status.
	resp, err := s.client.Do(req, nil)
	if err != nil {
		// CheckResponse in client.Do will typically return an ErrorResponse
		// if status code is 401 Unauthorized or other errors.
		return false, resp, err
	}

	// If we get here, the request was successful (status 2xx),
	// implying authentication was successful.
	return true, resp, nil
}

// Add this to client.go or a new services.go file to attach services to the client:
/*
func (c *Client) initServices() {
	c.Authentication = NewAuthenticationService(c)
	// Initialize other services here
}

// And call it in NewClient:
func NewClient(username, password string) *Client {
	// ... existing code ...
	c := &Client{ ... }
	c.initServices() // Call this
	return c
}

// And add to Client struct:
type Client struct {
    // ... existing fields ...
    Authentication *AuthenticationService
    // Other services
}
*/
