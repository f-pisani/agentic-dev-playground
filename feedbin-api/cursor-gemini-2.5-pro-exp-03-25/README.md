# Go Client for Feedbin API v2

This package provides a Go client for interacting with the [Feedbin API v2](https://github.com/feedbin/feedbin-api/tree/main/specs).

## Implementation Plan

1.  **Package Structure:**
    *   Root package directory: `feedbin-api/cursor-gemini-2.5-pro-exp-03-25`.
    *   Main client logic: `client.go`.
    *   API data structures (structs): `models.go`.
    *   Internal utility functions (request building, error handling, pagination parsing): `internal/util/util.go` (or in `client.go`).
    *   This `README.md`.
    *   `go.mod` file for dependency management (standard library only).

2.  **Client (`client.go`):**
    *   `Client` struct: Holds `baseURL`, `httpClient (*http.Client)`, `username`, `password`.
    *   `NewClient(username, password string) (*Client, error)`: Constructor, initializes with credentials and default `http.Client`. Allows optional custom client.
    *   `doRequest(ctx context.Context, method, path string, params url.Values, body io.Reader) (*http.Response, error)`: Helper for building and executing HTTP requests with Basic Auth, correct headers (`Content-Type`, `Accept`), URL construction, and context handling.
    *   `handleResponse(resp *http.Response, v interface{}) error`: Helper for checking status codes (returning `APIError` for non-2xx), decoding JSON responses, and closing the response body.

3.  **Models (`models.go`):**
    *   Define Go structs with `json:"..."` tags for API objects:
        *   `Subscription`: `ID`, `CreatedAt`, `FeedID`, `Title`, `FeedURL`, `SiteURL`, `JSONFeed` (optional).
        *   `Entry`: `ID`, `FeedID`, `Title` (\*string), `URL`, `ExtractedContentURL`, `Author` (\*string), `Content` (\*string), `Summary`, `Published` (time.Time), `CreatedAt` (time.Time), `Original` (\*OriginalEntry, optional), `Images` (\*Images, optional), `Enclosure` (\*Enclosure, optional), `TwitterID` (\*int64, optional), `TwitterThreadIDs` ([]int64, optional), `ExtractedArticles` ([]ExtractedArticle, optional), `JSONFeed` (optional). Define nested structs as needed. Handle potential `null` JSON values appropriately (e.g., using pointers).
        *   `FeedChoice`: `FeedURL`, `Title`.
        *   `APIError`: `StatusCode` (int), `Message` (string), `Body` (string). Implements `error` interface.
        *   `PaginationInfo`: `NextPageURL`, `LastPageURL`, `PrevPageURL`, `FirstPageURL`, `TotalRecords` (int).

4.  **Endpoints Implementation (in `client.go`):**
    *   **Authentication:**
        *   `VerifyCredentials(ctx context.Context) (bool, error)`: `GET /v2/authentication.json`.
    *   **Subscriptions:**
        *   `ListSubscriptions(ctx context.Context, options *ListSubscriptionsOptions) ([]Subscription, error)`: `GET /v2/subscriptions.json`. Handles `since`, `mode` options.
        *   `GetSubscription(ctx context.Context, id int64, options *GetSubscriptionOptions) (*Subscription, error)`: `GET /v2/subscriptions/{id}.json`. Handles `mode`.
        *   `CreateSubscription(ctx context.Context, feedURL string) (*Subscription, []FeedChoice, error)`: `POST /v2/subscriptions.json`. Handles 201, 302, 300 status codes.
        *   `UpdateSubscription(ctx context.Context, id int64, title string) (*Subscription, error)`: `PATCH /v2/subscriptions/{id}.json`.
        *   `DeleteSubscription(ctx context.Context, id int64) error`: `DELETE /v2/subscriptions/{id}.json`. Checks for 204.
    *   **Entries:**
        *   `ListEntries(ctx context.Context, options *ListEntriesOptions) ([]Entry, *PaginationInfo, error)`: `GET /v2/entries.json`. Handles various options (`page`, `since`, `ids`, `read`, `starred`, `per_page`, `mode`, includes). Parses `Link` and `X-Feedbin-Record-Count` headers.
        *   `ListFeedEntries(ctx context.Context, feedID int64, options *ListEntriesOptions) ([]Entry, *PaginationInfo, error)`: `GET /v2/feeds/{feedID}/entries.json`. Similar options and pagination handling.
        *   `GetEntry(ctx context.Context, id int64, options *GetEntryOptions) (*Entry, error)`: `GET /v2/entries/{id}.json`. Handles options.

5.  **Pagination:**
    *   Implement `parseLinkHeader(header string) map[string]string` helper (likely in `internal/util/util.go` or `client.go`).
    *   Populate `PaginationInfo` in list methods using `Link` and `X-Feedbin-Record-Count` headers.

6.  **Error Handling:**
    *   Use the custom `APIError` struct for non-2xx HTTP responses.

7.  **Standard Library Only:** Adhere strictly to the Go standard library.

8.  **Verification:** Ensure code passes `go vet ./...`.

## Usage (Example)

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	feedbin "path/to/your/feedbin-api/cursor-gemini-2.5-pro-exp-03-25" // Adjust import path
)

func main() {
	username := os.Getenv("FEEDBIN_USERNAME")
	password := os.Getenv("FEEDBIN_PASSWORD")

	if username == "" || password == "" {
		log.Fatal("FEEDBIN_USERNAME and FEEDBIN_PASSWORD environment variables must be set")
	}

	client, err := feedbin.NewClient(username, password)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	ctx := context.Background()

	// Verify Credentials
	valid, err := client.VerifyCredentials(ctx)
	if err != nil {
		log.Fatalf("Error verifying credentials: %v", err)
	}
	if !valid {
		log.Fatal("Invalid credentials")
	}
	fmt.Println("Credentials verified successfully!")

	// List Subscriptions
	subscriptions, err := client.ListSubscriptions(ctx, nil) // No options
	if err != nil {
		log.Fatalf("Error listing subscriptions: %v", err)
	}
	fmt.Printf("Found %d subscriptions:\\n", len(subscriptions))
	for _, sub := range subscriptions {
		fmt.Printf("- %s (ID: %d, Feed URL: %s)\\n", sub.Title, sub.ID, sub.FeedURL)
	}

	// List first page of unread entries
	entries, pagination, err := client.ListEntries(ctx, &feedbin.ListEntriesOptions{
		Read:    feedbin.Bool(false), // Pointer to bool for optional param
		PerPage: feedbin.Int(20),     // Pointer to int
	})
	if err != nil {
		log.Fatalf("Error listing entries: %v", err)
	}
	fmt.Printf("\\nFound %d unread entries (Page 1 of potentially %d total records):\\n", len(entries), pagination.TotalRecords)
	for _, entry := range entries {
		title := "Untitled"
		if entry.Title != nil {
			title = *entry.Title
		}
		fmt.Printf("- %s (ID: %d, Published: %s)\\n", title, entry.ID, entry.Published.Format("2006-01-02"))
	}

	if pagination.NextPageURL != "" {
		fmt.Printf("Next page URL: %s\\n", pagination.NextPageURL)
		// Here you could potentially fetch the next page using the URL
	}

}

```
(Note: Helper functions like `feedbin.Bool(false)` and `feedbin.Int(20)` would be needed to easily create pointers for optional parameters). 