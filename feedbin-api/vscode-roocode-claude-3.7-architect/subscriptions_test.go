package feedbin

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSubscriptionsList(t *testing.T) {
	// Sample response data
	subscriptions := []Subscription{
		{
			ID:        1,
			CreatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			FeedID:    101,
			Title:     "Test Feed 1",
			FeedURL:   "https://example.com/feed1.xml",
			SiteURL:   "https://example.com",
		},
		{
			ID:        2,
			CreatedAt: time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
			FeedID:    102,
			Title:     "Test Feed 2",
			FeedURL:   "https://example.com/feed2.xml",
			SiteURL:   "https://example.com",
		},
	}

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that the request is properly formed
		if r.URL.Path != "/subscriptions.json" {
			t.Errorf("Expected request to '/subscriptions.json', got '%s'", r.URL.Path)
		}

		// Check that the request method is GET
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		// Check for query parameters
		query := r.URL.Query()
		if since := query.Get("since"); since != "" {
			// If since parameter is provided, filter subscriptions
			// This is just a simple example, in a real test you might want to actually filter
			t.Logf("Since parameter provided: %s", since)
		}

		if mode := query.Get("mode"); mode == "extended" {
			// If extended mode is requested, add JSONFeed data
			for i := range subscriptions {
				subscriptions[i].JSONFeed = &JSONFeed{
					Favicon:     "https://example.com/favicon.ico",
					FeedURL:     subscriptions[i].FeedURL,
					Icon:        "https://example.com/icon.png",
					Version:     "https://jsonfeed.org/version/1",
					HomePageURL: subscriptions[i].SiteURL,
					Title:       subscriptions[i].Title,
				}
			}
		}

		// Return the subscriptions as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(subscriptions)
	}))
	defer server.Close()

	// Create a client that uses the test server
	client := NewClient(
		"test@example.com",
		"password",
		WithBaseURL(server.URL+"/"),
	)

	// Test listing subscriptions
	subs, err := client.Subscriptions.List(nil, false)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(subs) != 2 {
		t.Errorf("Expected 2 subscriptions, got %d", len(subs))
	}

	// Check the first subscription
	if subs[0].ID != 1 {
		t.Errorf("Expected ID 1, got %d", subs[0].ID)
	}
	if subs[0].Title != "Test Feed 1" {
		t.Errorf("Expected title 'Test Feed 1', got '%s'", subs[0].Title)
	}
	if subs[0].JSONFeed != nil {
		t.Error("Expected JSONFeed to be nil in normal mode")
	}

	// Test listing subscriptions in extended mode
	subs, err = client.Subscriptions.List(nil, true)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if subs[0].JSONFeed == nil {
		t.Error("Expected JSONFeed to be non-nil in extended mode")
	}
	if subs[0].JSONFeed.Title != "Test Feed 1" {
		t.Errorf("Expected JSONFeed title 'Test Feed 1', got '%s'", subs[0].JSONFeed.Title)
	}
}

func TestSubscriptionsGet(t *testing.T) {
	// Sample response data
	subscription := Subscription{
		ID:        1,
		CreatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		FeedID:    101,
		Title:     "Test Feed 1",
		FeedURL:   "https://example.com/feed1.xml",
		SiteURL:   "https://example.com",
	}

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that the request is properly formed
		if r.URL.Path != "/subscriptions/1.json" {
			t.Errorf("Expected request to '/subscriptions/1.json', got '%s'", r.URL.Path)
		}

		// Check that the request method is GET
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		// Return the subscription as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(subscription)
	}))
	defer server.Close()

	// Create a client that uses the test server
	client := NewClient(
		"test@example.com",
		"password",
		WithBaseURL(server.URL+"/"),
	)

	// Test getting a subscription
	sub, err := client.Subscriptions.Get(1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if sub.ID != 1 {
		t.Errorf("Expected ID 1, got %d", sub.ID)
	}
	if sub.Title != "Test Feed 1" {
		t.Errorf("Expected title 'Test Feed 1', got '%s'", sub.Title)
	}
}

func TestSubscriptionsCreate(t *testing.T) {
	// Sample response data
	subscription := Subscription{
		ID:        1,
		CreatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		FeedID:    101,
		Title:     "Test Feed 1",
		FeedURL:   "https://example.com/feed1.xml",
		SiteURL:   "https://example.com",
	}

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that the request is properly formed
		if r.URL.Path != "/subscriptions.json" {
			t.Errorf("Expected request to '/subscriptions.json', got '%s'", r.URL.Path)
		}

		// Check that the request method is POST
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Check that the request body contains the feed URL
		var req SubscriptionCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Error decoding request body: %v", err)
		}
		if req.FeedURL != "https://example.com/feed1.xml" {
			t.Errorf("Expected feed URL 'https://example.com/feed1.xml', got '%s'", req.FeedURL)
		}

		// Return the subscription as JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(subscription)
	}))
	defer server.Close()

	// Create a client that uses the test server
	client := NewClient(
		"test@example.com",
		"password",
		WithBaseURL(server.URL+"/"),
	)

	// Test creating a subscription
	sub, err := client.Subscriptions.Create("https://example.com/feed1.xml")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if sub.ID != 1 {
		t.Errorf("Expected ID 1, got %d", sub.ID)
	}
	if sub.Title != "Test Feed 1" {
		t.Errorf("Expected title 'Test Feed 1', got '%s'", sub.Title)
	}
}
