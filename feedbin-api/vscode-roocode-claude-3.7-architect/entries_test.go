package feedbin

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestEntriesList(t *testing.T) {
	// Sample response data
	entries := []Entry{
		{
			ID:                  1001,
			FeedID:              101,
			Title:               String("Test Entry 1"),
			URL:                 "https://example.com/entry1",
			ExtractedContentURL: "https://extract.feedbin.com/parser/feedbin/abc123?base64_url=xyz",
			Author:              String("Test Author"),
			Content:             String("<p>Test content 1</p>"),
			Summary:             String("Test summary 1"),
			Published:           time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			CreatedAt:           time.Date(2023, 1, 1, 12, 5, 0, 0, time.UTC),
		},
		{
			ID:                  1002,
			FeedID:              101,
			Title:               String("Test Entry 2"),
			URL:                 "https://example.com/entry2",
			ExtractedContentURL: "https://extract.feedbin.com/parser/feedbin/def456?base64_url=uvw",
			Author:              String("Test Author"),
			Content:             String("<p>Test content 2</p>"),
			Summary:             String("Test summary 2"),
			Published:           time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
			CreatedAt:           time.Date(2023, 1, 2, 12, 5, 0, 0, time.UTC),
		},
	}

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that the request is properly formed
		if r.URL.Path != "/entries.json" {
			t.Errorf("Expected request to '/entries.json', got '%s'", r.URL.Path)
		}

		// Check that the request method is GET
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		// Check for query parameters
		query := r.URL.Query()

		// Handle pagination
		page := query.Get("page")
		if page == "2" {
			// Return empty array for page 2 in this test
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Link", `<https://api.feedbin.com/v2/entries.json?page=1>; rel="first", <https://api.feedbin.com/v2/entries.json?page=1>; rel="prev"`)
			w.Header().Set("X-Feedbin-Record-Count", "2")
			json.NewEncoder(w).Encode([]Entry{})
			return
		}

		// Set pagination headers for page 1
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Link", `<https://api.feedbin.com/v2/entries.json?page=1>; rel="first", <https://api.feedbin.com/v2/entries.json?page=2>; rel="next", <https://api.feedbin.com/v2/entries.json?page=2>; rel="last"`)
		w.Header().Set("X-Feedbin-Record-Count", "2")

		// Check for other query parameters
		if starred := query.Get("starred"); starred == "true" {
			// If starred parameter is provided, we'd filter entries
			// For this test, we'll just return the same entries
			t.Logf("Starred parameter provided: %s", starred)
		}

		if mode := query.Get("mode"); mode == "extended" {
			// If extended mode is requested, add extended data
			for i := range entries {
				entries[i].Original = &EntryOriginal{
					Author:    String("Original Author"),
					Content:   String("<p>Original content</p>"),
					Title:     String("Original Title"),
					URL:       "https://example.com/original",
					EntryID:   "original-id",
					Published: time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
				}
			}
		}

		// Return the entries as JSON
		json.NewEncoder(w).Encode(entries)
	}))
	defer server.Close()

	// Create a client that uses the test server
	client := NewClient(
		"test@example.com",
		"password",
		WithBaseURL(server.URL+"/"),
	)

	// Test listing entries
	opts := &ListEntriesOptions{
		Page:    Int(1),
		Starred: Bool(true),
	}

	entries, pagination, err := client.Entries.List(opts)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(entries))
	}

	// Check the first entry
	if entries[0].ID != 1001 {
		t.Errorf("Expected ID 1001, got %d", entries[0].ID)
	}
	if *entries[0].Title != "Test Entry 1" {
		t.Errorf("Expected title 'Test Entry 1', got '%s'", *entries[0].Title)
	}

	// Check pagination links
	if pagination.First == "" {
		t.Error("Expected First pagination link to be non-empty")
	}
	if pagination.Next == "" {
		t.Error("Expected Next pagination link to be non-empty")
	}
	if pagination.Last == "" {
		t.Error("Expected Last pagination link to be non-empty")
	}
	if pagination.Prev != "" {
		t.Error("Expected Prev pagination link to be empty for first page")
	}

	// Test listing entries with extended mode
	opts = &ListEntriesOptions{
		Page:    Int(1),
		Starred: Bool(true),
		Mode:    "extended",
	}

	entries, _, err = client.Entries.List(opts)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check extended data
	if entries[0].Original == nil {
		t.Error("Expected Original to be non-nil in extended mode")
	} else {
		if *entries[0].Original.Title != "Original Title" {
			t.Errorf("Expected Original.Title to be 'Original Title', got '%s'", *entries[0].Original.Title)
		}
	}

	// Test listing entries with page 2
	opts = &ListEntriesOptions{
		Page: Int(2),
	}

	entries, pagination, err = client.Entries.List(opts)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("Expected 0 entries for page 2, got %d", len(entries))
	}

	// Check pagination links for page 2
	if pagination.First == "" {
		t.Error("Expected First pagination link to be non-empty")
	}
	if pagination.Prev == "" {
		t.Error("Expected Prev pagination link to be non-empty for second page")
	}
	if pagination.Next != "" {
		t.Error("Expected Next pagination link to be empty for last page")
	}
	if pagination.Last != "" {
		t.Error("Expected Last pagination link to be empty for last page")
	}
}

func TestEntriesGet(t *testing.T) {
	// Sample response data
	entry := Entry{
		ID:                  1001,
		FeedID:              101,
		Title:               String("Test Entry 1"),
		URL:                 "https://example.com/entry1",
		ExtractedContentURL: "https://extract.feedbin.com/parser/feedbin/abc123?base64_url=xyz",
		Author:              String("Test Author"),
		Content:             String("<p>Test content 1</p>"),
		Summary:             String("Test summary 1"),
		Published:           time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		CreatedAt:           time.Date(2023, 1, 1, 12, 5, 0, 0, time.UTC),
	}

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that the request is properly formed
		if r.URL.Path != "/entries/1001.json" {
			t.Errorf("Expected request to '/entries/1001.json', got '%s'", r.URL.Path)
		}

		// Check that the request method is GET
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		// Return the entry as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(entry)
	}))
	defer server.Close()

	// Create a client that uses the test server
	client := NewClient(
		"test@example.com",
		"password",
		WithBaseURL(server.URL+"/"),
	)

	// Test getting an entry
	e, err := client.Entries.Get(1001, nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if e.ID != 1001 {
		t.Errorf("Expected ID 1001, got %d", e.ID)
	}
	if *e.Title != "Test Entry 1" {
		t.Errorf("Expected title 'Test Entry 1', got '%s'", *e.Title)
	}
}
