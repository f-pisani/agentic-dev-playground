package feedbinapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

// TestNewClient ensures that NewClient creates a client with correct defaults
// and that services are initialized.
func TestNewClient(t *testing.T) {
	username := "testuser"
	password := "testpass"
	c := NewClient(username, password)

	if c.username != username {
		t.Errorf("NewClient username = %q; want %q", c.username, username)
	}
	if c.password != password {
		t.Errorf("NewClient password = %q; want %q", c.password, password)
	}
	if c.UserAgent != UserAgent { // UserAgent is the default constant
		t.Errorf("NewClient UserAgent = %q; want %q", c.UserAgent, UserAgent)
	}
	expectedBaseURL, _ := url.Parse(BaseURL)
	if c.baseURL.String() != expectedBaseURL.String() {
		t.Errorf("NewClient BaseURL = %q; want %q", c.baseURL.String(), expectedBaseURL.String())
	}

	// Check if services are initialized (not nil)
	if c.Authentication == nil {
		t.Error("Authentication service not initialized")
	}
	if c.Subscriptions == nil {
		t.Error("Subscriptions service not initialized")
	}
	if c.Entries == nil {
		t.Error("Entries service not initialized")
	}
	// Add checks for other services as well
	if c.UnreadEntries == nil {
		t.Error("UnreadEntries service not initialized")
	}
	if c.StarredEntries == nil {
		t.Error("StarredEntries service not initialized")
	}
	if c.Taggings == nil {
		t.Error("Taggings service not initialized")
	}
	if c.Tags == nil {
		t.Error("Tags service not initialized")
	}
	if c.SavedSearches == nil {
		t.Error("SavedSearches service not initialized")
	}
	if c.RecentlyReadEntries == nil {
		t.Error("RecentlyReadEntries service not initialized")
	}
	if c.UpdatedEntries == nil {
		t.Error("UpdatedEntries service not initialized")
	}
	if c.Icons == nil {
		t.Error("Icons service not initialized")
	}
	if c.Imports == nil {
		t.Error("Imports service not initialized")
	}
	if c.Pages == nil {
		t.Error("Pages service not initialized")
	}
	if c.Extract == nil {
		t.Error("Extract service not initialized")
	}
}

// TestClient_Setters tests the various setter methods on the client.
func TestClient_Setters(t *testing.T) {
	c := NewClient("user", "pass")

	newURL := "http://localhost:8080/api/v2/"
	err := c.SetBaseURL(newURL)
	if err != nil {
		t.Fatalf("SetBaseURL returned an error: %v", err)
	}
	if c.baseURL.String() != newURL {
		t.Errorf("SetBaseURL base URL = %q; want %q", c.baseURL.String(), newURL)
	}

	newUA := "MyCustomAgent/1.0"
	c.SetUserAgent(newUA)
	if c.UserAgent != newUA {
		t.Errorf("SetUserAgent UserAgent = %q; want %q", c.UserAgent, newUA)
	}

	newTimeout := 60 * time.Second
	c.SetTimeout(newTimeout)
	if c.client.Timeout != newTimeout {
		t.Errorf("SetTimeout Timeout = %v; want %v", c.client.Timeout, newTimeout)
	}
}

// setupTestServer creates a new httptest.Server and a new feedbin.Client that targets it.
// It also returns a mux for registering handlers and a teardown func.
func setupTestServer(t *testing.T) (client *Client, mux *http.ServeMux, teardown func()) {
	mux = http.NewServeMux()
	server := httptest.NewServer(mux)

	client = NewClient("testuser", "testpass")
	err := client.SetBaseURL(server.URL + "/") // Ensure trailing slash for test server
	if err != nil {
		server.Close()
		t.Fatalf("Failed to set base URL for test client: %v", err)
	}

	return client, mux, server.Close
}


// TestAuthentication_Verify_Success tests a successful authentication.
func TestAuthentication_Verify_Success(t *testing.T) {
	client, mux, teardown := setupTestServer(t)
	defer teardown()

	mux.HandleFunc("/subscriptions.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Request method = %v, want %v", r.Method, http.MethodGet)
		}
		// Check basic auth header
		user, pass, ok := r.BasicAuth()
		if !ok || user != "testuser" || pass != "testpass" {
			w.WriteHeader(http.StatusUnauthorized)
			t.Error("Basic auth header not set or incorrect")
			return
		}
		if r.URL.Query().Get("per_page") != "1" {
			t.Errorf("per_page query param = %q, want %q", r.URL.Query().Get("per_page"), "1")
		}
		fmt.Fprint(w, `[]`) // Empty array of subscriptions
	})

	ok, _, err := client.Authentication.Verify()
	if err != nil {
		t.Fatalf("Authentication.Verify returned error: %v", err)
	}
	if !ok {
		t.Error("Authentication.Verify returned false, want true for success")
	}
}

// TestAuthentication_Verify_Failure tests a failed authentication (401).
func TestAuthentication_Verify_Failure(t *testing.T) {
	client, mux, teardown := setupTestServer(t)
	defer teardown()

	mux.HandleFunc("/subscriptions.json", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized) // Simulate auth failure
		fmt.Fprint(w, `{"error": "Unauthorized"}`)
	})

	ok, resp, err := client.Authentication.Verify()
	if err == nil {
		t.Fatal("Authentication.Verify expected error for 401, got nil")
	}
	if ok {
		t.Error("Authentication.Verify returned true, want false for failure")
	}
	if resp == nil || resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected HTTP 401, got status %d", resp.StatusCode)
	}
}

// TestSubscriptions_List_Simple is a very basic test for listing subscriptions.
func TestSubscriptions_List_Simple(t *testing.T) {
	client, mux, teardown := setupTestServer(t)
	defer teardown()

	mux.HandleFunc("/subscriptions.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Request method = %v, want %v", r.Method, http.MethodGet)
		}
		// Minimal valid response
		fmt.Fprint(w, `[{"id": 1, "feed_id": 101, "title": "Test Feed", "feed_url": "http://example.com/feed", "site_url": "http://example.com", "created_at": "2023-01-01T00:00:00Z"}]`)
	})

	subs, _, err := client.Subscriptions.List(nil)
	if err != nil {
		t.Fatalf("Subscriptions.List returned error: %v", err)
	}

	if len(subs) != 1 {
		t.Errorf("Expected 1 subscription, got %d", len(subs))
	}
	if subs[0].ID != 1 || subs[0].Title != "Test Feed" {
		t.Errorf("Unexpected subscription data: %+v", subs[0])
	}
}

// Add more tests for other services, e.g., Entries.List, UnreadEntries.List, etc.
// Example for Entries.List
func TestEntries_List_Simple(t *testing.T) {
	client, mux, teardown := setupTestServer(t)
	defer teardown()

	mux.HandleFunc("/entries.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[{"id": 1, "feed_id": 101, "title": "Test Entry", "url": "http://example.com/entry1", "published_at": "2023-01-01T10:00:00Z", "created_at": "2023-01-01T10:00:00Z"}]`)
	})

	entries, _, err := client.Entries.List(nil)
	if err != nil {
		t.Fatalf("Entries.List returned error: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(entries))
	}
	if entries[0].ID != 1 || entries[0].Title != "Test Entry" {
		t.Errorf("Unexpected entry data: %+v", entries[0])
	}
}

// TestParseFeedbinTime tests the time parsing logic.
func TestParseFeedbinTime(t *testing.T) {
	tests := []struct {
		name    string
		timeStr string
		wantErr bool
	}{
		{"RFC3339Nano UTC", "2013-02-19T15:33:38.449047Z", false},
		{"RFC3339Nano Offset", "2013-02-19T07:33:38.449047-08:00", false},
		{"Microsecond UTC", "2006-01-02T15:04:05.999999Z", false},
		{"Microsecond Offset", "2006-01-02T15:04:05.999999-07:00", false},
		{"No Fractional UTC", "2006-01-02T15:04:05Z", false},
		{"No Fractional Offset", "2006-01-02T15:04:05-07:00", false},
		{"Invalid", "not-a-time", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseFeedbinTime(tt.timeStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFeedbinTime(%q) error = %v, wantErr %v", tt.timeStr, err, tt.wantErr)
			}
		})
	}
}

// TestFormatFeedbinTime tests time formatting.
func TestFormatFeedbinTime(t *testing.T) {
	// Example time: 2023-10-26 10:30:45.123456 UTC
	loc, _ := time.LoadLocation("UTC")
	tm := time.Date(2023, 10, 26, 10, 30, 45, 123456000, loc)
	expected := "2023-10-26T10:30:45.123456Z"
	formatted := FormatFeedbinTime(tm)
	if formatted != expected {
		t.Errorf("FormatFeedbinTime() = %q, want %q", formatted, expected)
	}
}
