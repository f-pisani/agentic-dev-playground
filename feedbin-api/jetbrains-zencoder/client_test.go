package feedbin

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient("username", "password")

	if client.Username != "username" {
		t.Errorf("Expected username to be 'username', got '%s'", client.Username)
	}

	if client.Password != "password" {
		t.Errorf("Expected password to be 'password', got '%s'", client.Password)
	}

	if client.BaseURL != BaseURL {
		t.Errorf("Expected BaseURL to be '%s', got '%s'", BaseURL, client.BaseURL)
	}

	if client.HTTPClient == nil {
		t.Error("Expected HTTPClient to be initialized")
	}
}

func TestNewClientWithOptions(t *testing.T) {
	customBaseURL := "https://custom.example.com/api"
	customTimeout := 60 * time.Second

	client := NewClientWithOptions("username", "password", customBaseURL, customTimeout)

	if client.Username != "username" {
		t.Errorf("Expected username to be 'username', got '%s'", client.Username)
	}

	if client.Password != "password" {
		t.Errorf("Expected password to be 'password', got '%s'", client.Password)
	}

	if client.BaseURL != customBaseURL {
		t.Errorf("Expected BaseURL to be '%s', got '%s'", customBaseURL, client.BaseURL)
	}

	if client.HTTPClient == nil {
		t.Error("Expected HTTPClient to be initialized")
	}

	if client.HTTPClient.Timeout != customTimeout {
		t.Errorf("Expected timeout to be %v, got %v", customTimeout, client.HTTPClient.Timeout)
	}
}

func TestValidateAuth(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that the request is properly formed
		if r.URL.Path != "/authentication.json" {
			t.Errorf("Expected path to be '/authentication.json', got '%s'", r.URL.Path)
		}

		if r.Method != http.MethodGet {
			t.Errorf("Expected method to be 'GET', got '%s'", r.Method)
		}

		// Check for basic auth
		username, password, ok := r.BasicAuth()
		if !ok {
			t.Error("Expected basic auth to be set")
		}

		if username != "valid_user" {
			t.Errorf("Expected username to be 'valid_user', got '%s'", username)
		}

		if password != "valid_pass" {
			t.Errorf("Expected password to be 'valid_pass', got '%s'", password)
		}

		// Return a 200 OK response
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create a client that uses the test server
	client := NewClientWithOptions("valid_user", "valid_pass", server.URL, 30*time.Second)

	// Test the ValidateAuth method
	valid, err := client.ValidateAuth()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !valid {
		t.Error("Expected authentication to be valid")
	}
}

func TestParsePaginationLinks(t *testing.T) {
	// Create a test response with Link header
	resp := &http.Response{
		Header: http.Header{},
	}

	// Set the Link header
	resp.Header.Set("Link",
		`<https://api.feedbin.com/v2/entries.json?page=2>; rel="next", 
		<https://api.feedbin.com/v2/entries.json?page=5>; rel="last", 
		<https://api.feedbin.com/v2/entries.json?page=1>; rel="first", 
		<https://api.feedbin.com/v2/entries.json?page=1>; rel="prev"`)

	// Parse the pagination links
	links := parsePaginationLinks(resp)

	// Check that the links were parsed correctly
	if links.Next != "https://api.feedbin.com/v2/entries.json?page=2" {
		t.Errorf("Expected next link to be 'https://api.feedbin.com/v2/entries.json?page=2', got '%s'", links.Next)
	}

	if links.Last != "https://api.feedbin.com/v2/entries.json?page=5" {
		t.Errorf("Expected last link to be 'https://api.feedbin.com/v2/entries.json?page=5', got '%s'", links.Last)
	}

	if links.First != "https://api.feedbin.com/v2/entries.json?page=1" {
		t.Errorf("Expected first link to be 'https://api.feedbin.com/v2/entries.json?page=1', got '%s'", links.First)
	}

	if links.Prev != "https://api.feedbin.com/v2/entries.json?page=1" {
		t.Errorf("Expected prev link to be 'https://api.feedbin.com/v2/entries.json?page=1', got '%s'", links.Prev)
	}
}
