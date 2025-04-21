package feedbin

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthenticationVerify(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that the request is properly formed
		if r.URL.Path != "/authentication.json" {
			t.Errorf("Expected request to '/authentication.json', got '%s'", r.URL.Path)
		}

		// Check that the request includes basic auth
		username, password, ok := r.BasicAuth()
		if !ok {
			t.Error("Expected request to include basic auth")
		}

		// Check that the credentials are correct
		if username != "test@example.com" || password != "password" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Return a successful response
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create a client that uses the test server
	client := NewClient(
		"test@example.com",
		"password",
		WithBaseURL(server.URL+"/"),
	)

	// Test with valid credentials
	valid, err := client.Authentication.Verify()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !valid {
		t.Error("Expected valid credentials")
	}

	// Test with invalid credentials
	client = NewClient(
		"test@example.com",
		"wrong-password",
		WithBaseURL(server.URL+"/"),
	)

	valid, err = client.Authentication.Verify()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if valid {
		t.Error("Expected invalid credentials")
	}
}
