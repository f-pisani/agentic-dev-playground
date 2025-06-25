# Feedbin API Client Implementation Plan

This document outlines the plan for creating a Go client for the Feedbin API.

## 1. Project Structure

The project will be organized into the following Go packages and files:

*   `feedbin/` (root package directory)
    *   `client.go`: Defines the main `Client` struct and methods for interacting with the API.
    *   `models.go`: Contains the Go structs representing the API data models (e.g., `Entry`, `Feed`, `Subscription`).
    *   `options.go`: Defines optional parameters for API calls (e.g., for pagination, filtering).
    *   `errors.go`: Defines custom error types for the client.
    *   `authentication.go`: Handles authentication logic.
    *   `endpoints.go`: Contains endpoint definitions and methods for each resource.

## 2. Main `Client` Struct

The primary `Client` struct will be the entry point for all API interactions.

```go
// in client.go
package feedbin

import "net/http"

type Client struct {
    httpClient    *http.Client
    baseURL       string
    authenticator Authenticator
}

// NewClient creates a new Feedbin API client.
func NewClient(username, password string, options ...ClientOption) *Client {
    // ... implementation ...
}
```

## 3. Authentication

Authentication will be handled using HTTP Basic Auth. An `Authenticator` interface will be created to allow for different authentication strategies in the future, though the initial implementation will focus on basic auth.

```go
// in authentication.go
package feedbin

import "net/http"

type Authenticator interface {
    Authenticate(req *http.Request)
}

type BasicAuth struct {
    Username string
    Password string
}

func (b *BasicAuth) Authenticate(req *http.Request) {
    req.SetBasicAuth(b.Username, b.Password)
}
```

## 4. API Methods

The client will implement methods for all API resources, grouped by resource. Each resource will have its own file (e.g., `entries.go`, `feeds.go`).

### Entries

*   `GetEntries(opts ...RequestOption) ([]*Entry, error)`
*   `GetFeedEntries(feedID int64, opts ...RequestOption) ([]*Entry, error)`
*   `GetEntry(entryID int64, opts ...RequestOption) (*Entry, error)`

### Feeds

*   `GetFeed(feedID int64) (*Feed, error)`

### Subscriptions

*   `GetSubscriptions(opts ...RequestOption) ([]*Subscription, error)`
*   `GetSubscription(subID int64) (*Subscription, error)`
*   `CreateSubscription(feedURL string) (*Subscription, error)`
*   `UpdateSubscription(subID int64, title string) (*Subscription, error)`
*   `DeleteSubscription(subID int64) error`

### Taggings

*   `GetTaggings() ([]*Tagging, error)`
*   `GetTagging(taggingID int64) (*Tagging, error)`
*   `CreateTagging(feedID int64, name string) (*Tagging, error)`
*   `DeleteTagging(taggingID int64) error`

### Starred Entries

*   `GetStarredEntryIDs() ([]int64, error)`
*   `StarEntries(entryIDs []int64) error`
*   `UnstarEntries(entryIDs []int64) error`

### Unread Entries

*   `GetUnreadEntryIDs() ([]int64, error)`
*   `MarkEntriesAsRead(entryIDs []int64) error`
*   `MarkEntriesAsUnread(entryIDs []int64) error`

### Recently Read Entries

*   `GetRecentlyReadEntryIDs() ([]int64, error)`
*   `CreateRecentlyReadEntries(entryIDs []int64) error`

### Saved Searches

*   `GetSavedSearches() ([]*SavedSearch, error)`
*   `GetSavedSearch(searchID int64, opts ...RequestOption) ([]int64, error)`
*   `CreateSavedSearch(name, query string) (*SavedSearch, error)`
*   `UpdateSavedSearch(searchID int64, name string) (*SavedSearch, error)`
*   `DeleteSavedSearch(searchID int64) error`

### And so on for all other resources...

## 5. Pagination

Pagination will be handled by passing `page` and `per_page` parameters in the request options. The client will not automatically handle multiple pages; it will be the responsibility of the user of the client to request subsequent pages.

```go
// in options.go
package feedbin

// ...

func WithPage(page int) RequestOption {
    // ...
}

func WithPerPage(perPage int) RequestOption {
    // ...
}
```

## 6. Error Handling

A custom error struct will be used to provide more context on API errors.

```go
// in errors.go
package feedbin

import "fmt"

type APIError struct {
    StatusCode int
    Body       string
    Message    string
}

func (e *APIError) Error() string {
    return fmt.Sprintf("API Error: %d %s - %s", e.StatusCode, e.Message, e.Body)
}
```

The client will check the HTTP status code of each response and return an `*APIError` if the request was not successful.

## 7. Data Models

Go structs will be defined in `models.go` to represent the JSON data structures returned by the API.

```go
// in models.go
package feedbin

import "time"

type Entry struct {
    ID                  int64      `json:"id"`
    FeedID              int64      `json:"feed_id"`
    Title               *string    `json:"title"`
    URL                 string     `json:"url"`
    Author              *string    `json:"author"`
    Content             *string    `json:"content"`
    Summary             string     `json:"summary"`
    Published           time.Time  `json:"published"`
    CreatedAt           time.Time  `json:"created_at"`
    ExtractedContentURL string     `json:"extracted_content_url"`
    // ... other fields
}

type Feed struct {
    ID      int64  `json:"id"`
    Title   string `json:"title"`
    FeedURL string `json:"feed_url"`
    SiteURL string `json:"site_url"`
}

// ... and so on for all other models
```

## 8. Content Extraction

The client will provide a helper function to generate the signed URL for the content extraction service.

```go
// in client.go
package feedbin

func (c *Client) GetContentExtractorURL(url string) (string, error) {
    // ... implementation to generate signed URL ...
}
```

This function will require the user to provide their content extraction username and secret.