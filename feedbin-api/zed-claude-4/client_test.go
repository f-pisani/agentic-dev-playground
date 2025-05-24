package feedbin

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient("test@example.com", "password")

	if client == nil {
		t.Fatal("NewClient returned nil")
	}

	if client.credentials.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", client.credentials.Email)
	}

	if client.credentials.Password != "password" {
		t.Errorf("Expected password 'password', got '%s'", client.credentials.Password)
	}
}

func TestCredentialsValidation(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		password string
		wantErr  bool
	}{
		{"valid credentials", "test@example.com", "password", false},
		{"empty email", "", "password", true},
		{"empty password", "test@example.com", "", true},
		{"both empty", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			creds := &Credentials{
				Email:    tt.email,
				Password: tt.password,
			}

			err := creds.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBasicAuth(t *testing.T) {
	creds := &Credentials{
		Email:    "test@example.com",
		Password: "password",
	}

	auth := creds.BasicAuth()
	expected := "dGVzdEBleGFtcGxlLmNvbTpwYXNzd29yZA==" // base64 of "test@example.com:password"

	if auth != expected {
		t.Errorf("BasicAuth() = %s, want %s", auth, expected)
	}
}

func TestAuthenticate(t *testing.T) {
	tests := []struct {
		name          string
		statusCode    int
		expectedError bool
		errorType     string
	}{
		{"successful auth", http.StatusOK, false, ""},
		{"unauthorized", http.StatusUnauthorized, true, "*feedbin.AuthenticationError"},
		{"server error", http.StatusInternalServerError, true, "*feedbin.AuthenticationError"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/authentication.json" {
					t.Errorf("Expected path '/authentication.json', got '%s'", r.URL.Path)
				}

				if r.Method != "GET" {
					t.Errorf("Expected method 'GET', got '%s'", r.Method)
				}

				// Check authorization header
				auth := r.Header.Get("Authorization")
				if auth == "" {
					t.Error("Missing Authorization header")
				}

				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := NewClient("test@example.com", "password")
			client.SetBaseURL(server.URL)

			err := client.Authenticate(context.Background())

			if (err != nil) != tt.expectedError {
				t.Errorf("Authenticate() error = %v, wantErr %v", err, tt.expectedError)
			}
		})
	}
}

func TestGetSubscriptions(t *testing.T) {
	subscriptions := []Subscription{
		{
			ID:        1,
			CreatedAt: time.Now(),
			FeedID:    101,
			Title:     "Test Feed",
			FeedURL:   "https://example.com/feed.xml",
			SiteURL:   "https://example.com",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/subscriptions.json" {
			t.Errorf("Expected path '/subscriptions.json', got '%s'", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(subscriptions)
	}))
	defer server.Close()

	client := NewClient("test@example.com", "password")
	client.SetBaseURL(server.URL)

	result, pagination, err := client.GetSubscriptions(context.Background(), nil)
	if err != nil {
		t.Fatalf("GetSubscriptions() error = %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 subscription, got %d", len(result))
	}

	if result[0].Title != "Test Feed" {
		t.Errorf("Expected title 'Test Feed', got '%s'", result[0].Title)
	}

	if pagination == nil {
		t.Error("Expected pagination info, got nil")
	}
}

func TestCreateSubscription(t *testing.T) {
	subscription := Subscription{
		ID:        1,
		CreatedAt: time.Now(),
		FeedID:    101,
		Title:     "New Feed",
		FeedURL:   "https://example.com/feed.xml",
		SiteURL:   "https://example.com",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/subscriptions.json" {
			t.Errorf("Expected path '/subscriptions.json', got '%s'", r.URL.Path)
		}

		if r.Method != "POST" {
			t.Errorf("Expected method 'POST', got '%s'", r.Method)
		}

		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json; charset=utf-8" {
			t.Errorf("Expected content type 'application/json; charset=utf-8', got '%s'", contentType)
		}

		var req CreateSubscriptionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}

		if req.FeedURL != "https://example.com/feed.xml" {
			t.Errorf("Expected feed URL 'https://example.com/feed.xml', got '%s'", req.FeedURL)
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(subscription)
	}))
	defer server.Close()

	client := NewClient("test@example.com", "password")
	client.SetBaseURL(server.URL)

	result, choices, err := client.CreateSubscription(context.Background(), "https://example.com/feed.xml")
	if err != nil {
		t.Fatalf("CreateSubscription() error = %v", err)
	}

	if choices != nil {
		t.Error("Expected no choices for successful creation")
	}

	if result.Title != "New Feed" {
		t.Errorf("Expected title 'New Feed', got '%s'", result.Title)
	}
}

func TestCreateSubscriptionMultipleChoices(t *testing.T) {
	choices := []FeedChoice{
		{FeedURL: "https://example.com/feed1.xml", Title: "Feed 1"},
		{FeedURL: "https://example.com/feed2.xml", Title: "Feed 2"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMultipleChoices)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(choices)
	}))
	defer server.Close()

	client := NewClient("test@example.com", "password")
	client.SetBaseURL(server.URL)

	result, _, err := client.CreateSubscription(context.Background(), "https://example.com")
	if err == nil {
		t.Fatal("Expected error for multiple choices")
	}

	if result != nil {
		t.Error("Expected nil result for multiple choices")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("Expected APIError, got %T", err)
	}

	if !apiErr.IsMultipleChoices() {
		t.Error("Expected multiple choices error")
	}

	// Note: In a real implementation, choices would be properly decoded
	// This test demonstrates the structure
}

func TestGetUnreadEntries(t *testing.T) {
	entryIDs := []int{1, 2, 3, 4, 5}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/unread_entries.json" {
			t.Errorf("Expected path '/unread_entries.json', got '%s'", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(entryIDs)
	}))
	defer server.Close()

	client := NewClient("test@example.com", "password")
	client.SetBaseURL(server.URL)

	result, err := client.GetUnreadEntries(context.Background())
	if err != nil {
		t.Fatalf("GetUnreadEntries() error = %v", err)
	}

	if len(result) != 5 {
		t.Errorf("Expected 5 entry IDs, got %d", len(result))
	}

	for i, id := range result {
		if id != entryIDs[i] {
			t.Errorf("Expected entry ID %d, got %d", entryIDs[i], id)
		}
	}
}

func TestAPIErrorTypes(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		checkFunc  func(*APIError) bool
	}{
		{"not found", http.StatusNotFound, (*APIError).IsNotFound},
		{"unauthorized", http.StatusUnauthorized, (*APIError).IsUnauthorized},
		{"forbidden", http.StatusForbidden, (*APIError).IsForbidden},
		{"unsupported media type", http.StatusUnsupportedMediaType, (*APIError).IsUnsupportedMediaType},
		{"multiple choices", http.StatusMultipleChoices, (*APIError).IsMultipleChoices},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := NewClient("test@example.com", "password")
			client.SetBaseURL(server.URL)

			_, _, err := client.GetSubscriptions(context.Background(), nil)
			if err == nil {
				t.Fatal("Expected error")
			}

			apiErr, ok := err.(*APIError)
			if !ok {
				t.Fatalf("Expected APIError, got %T", err)
			}

			if !tt.checkFunc(apiErr) {
				t.Errorf("Error type check failed for status %d", tt.statusCode)
			}
		})
	}
}

func TestPaginationInfo(t *testing.T) {
	pagination := &PaginationInfo{
		Next:     "https://api.example.com/page2",
		Previous: "https://api.example.com/page1",
	}

	if !pagination.HasNext() {
		t.Error("Expected HasNext() to return true")
	}

	if !pagination.HasPrevious() {
		t.Error("Expected HasPrevious() to return true")
	}

	emptyPagination := &PaginationInfo{}
	if emptyPagination.HasNext() {
		t.Error("Expected HasNext() to return false for empty pagination")
	}

	if emptyPagination.HasPrevious() {
		t.Error("Expected HasPrevious() to return false for empty pagination")
	}
}
