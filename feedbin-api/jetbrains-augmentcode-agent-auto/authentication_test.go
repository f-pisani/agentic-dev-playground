package feedbin

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthenticationService_Verify(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the request method
		if r.Method != http.MethodGet {
			t.Errorf("Expected 'GET' request, got '%s'", r.Method)
		}

		// Check the request path
		if r.URL.Path != "/v2/authentication.json" {
			t.Errorf("Expected request to '/v2/authentication.json', got '%s'", r.URL.Path)
		}

		// Check the auth header
		username, password, ok := r.BasicAuth()
		if !ok {
			t.Error("Expected Basic Auth header to be set")
		}
		if username != "username" {
			t.Errorf("Expected username to be 'username', got '%s'", username)
		}
		if password != "password" {
			t.Errorf("Expected password to be 'password', got '%s'", password)
		}

		// Return a successful response
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create a client using the test server URL
	client := NewClient("username", "password")
	client.SetBaseURL(server.URL)

	// Test the Verify method
	err := client.Authentication.Verify()
	if err != nil {
		t.Errorf("Authentication.Verify returned error: %v", err)
	}
}

func TestAuthenticationService_VerifyFail(t *testing.T) {
	// Create a test server that returns an unauthorized response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	// Create a client using the test server URL
	client := NewClient("username", "password")
	client.SetBaseURL(server.URL)

	// Test the Verify method
	err := client.Authentication.Verify()
	if err == nil {
		t.Error("Authentication.Verify should have returned an error")
	}
}
