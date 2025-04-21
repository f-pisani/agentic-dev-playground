# Feedbin API Client for Go

This package provides an idiomatic Go client for the Feedbin REST API. It allows Go applications to interact with the Feedbin service for managing RSS feeds, entries, and related functionality.

## Implementation Plan

### 1. Package Structure

The client will be organized into the following package structure:

```
feedbin/
├── client.go       # Main client implementation
├── auth.go         # Authentication functionality
├── models/         # Data models for API objects
│   ├── entry.go
│   ├── feed.go
│   ├── subscription.go
│   └── ...
├── endpoints/      # API endpoint implementations
│   ├── entries.go
│   ├── subscriptions.go
│   ├── unread.go
│   ├── starred.go
│   └── ...
└── errors/         # Custom error types
    └── errors.go
```

### 2. Core Components

#### 2.1 Client

The main client will provide:
- Configuration options (base URL, timeout settings)
- HTTP client management
- Authentication handling
- Methods to access different API endpoints

```go
// Client represents a Feedbin API client
type Client struct {
    // BaseURL is the base URL for API requests
    BaseURL string
    
    // HTTPClient is the HTTP client used for making requests
    HTTPClient *http.Client
    
    // Authentication credentials
    Email    string
    Password string
    
    // API endpoints
    Subscriptions *SubscriptionsService
    Entries       *EntriesService
    UnreadEntries *UnreadEntriesService
    StarredEntries *StarredEntriesService
    Tags          *TagsService
    Taggings      *TaggingsService
    SavedSearches *SavedSearchesService
    // Additional services...
}

// NewClient creates a new Feedbin API client
func NewClient(email, password string, options ...ClientOption) *Client {
    // Initialize client with default values
    // Apply options
    // Initialize services
    // Return client
}
```

#### 2.2 Authentication

The client will support HTTP Basic Authentication as required by the Feedbin API.

```go
// AuthenticationService handles authentication-related operations
type AuthenticationService struct {
    client *Client
}

// Verify checks if the provided credentials are valid
func (s *AuthenticationService) Verify() (bool, error) {
    // Make request to /authentication.json
    // Return true if status code is 200, false otherwise
}
```

#### 2.3 Models

Data structures representing Feedbin objects:

```go
// Subscription represents a Feedbin subscription
type Subscription struct {
    ID        int       `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    FeedID    int       `json:"feed_id"`
    Title     string    `json:"title"`
    FeedURL   string    `json:"feed_url"`
    SiteURL   string    `json:"site_url"`
    JSONFeed  *JSONFeed `json:"json_feed,omitempty"`
}

// Entry represents a Feedbin entry
type Entry struct {
    ID                  int       `json:"id"`
    FeedID              int       `json:"feed_id"`
    Title               *string   `json:"title"`
    URL                 string    `json:"url"`
    ExtractedContentURL string    `json:"extracted_content_url"`
    Author              *string   `json:"author"`
    Content             *string   `json:"content"`
    Summary             *string   `json:"summary"`
    Published           time.Time `json:"published"`
    CreatedAt           time.Time `json:"created_at"`
    
    // Extended mode fields
    Original          *EntryOriginal     `json:"original,omitempty"`
    Images            *EntryImages       `json:"images,omitempty"`
    Enclosure         *EntryEnclosure    `json:"enclosure,omitempty"`
    TwitterID         *int64             `json:"twitter_id,omitempty"`
    TwitterThreadIDs  []int64            `json:"twitter_thread_ids,omitempty"`
    ExtractedArticles []ExtractedArticle `json:"extracted_articles,omitempty"`
    JSONFeed          *JSONFeed          `json:"json_feed,omitempty"`
}
```

#### 2.4 Error Handling

Custom error types that include:

```go
// APIError represents an error returned by the Feedbin API
type APIError struct {
    StatusCode int
    Body       string
    Message    string
}

// Error returns the error message
func (e *APIError) Error() string {
    return e.Message
}
```

#### 2.5 Pagination

Support for handling paginated responses:

```go
// PaginationLinks represents pagination links from the API
type PaginationLinks struct {
    First string
    Prev  string
    Next  string
    Last  string
}

// parsePaginationLinks parses the Link header for pagination links
func parsePaginationLinks(header string) PaginationLinks {
    // Parse Link header and extract URLs
}

// getTotalRecords extracts the total record count from the X-Feedbin-Record-Count header
func getTotalRecords(header string) (int, error) {
    // Parse and return the record count
}
```

### 3. API Endpoints Implementation

#### 3.1 Subscriptions

```go
// SubscriptionsService handles subscription-related operations
type SubscriptionsService struct {
    client *Client
}

// List returns all subscriptions
func (s *SubscriptionsService) List(since *time.Time, extended bool) ([]models.Subscription, error) {
    // Make request to /subscriptions.json with optional parameters
}

// Get returns a specific subscription
func (s *SubscriptionsService) Get(id int) (*models.Subscription, error) {
    // Make request to /subscriptions/{id}.json
}

// Create creates a new subscription
func (s *SubscriptionsService) Create(feedURL string) (*models.Subscription, error) {
    // Make request to /subscriptions.json with feed_url
}

// Delete deletes a subscription
func (s *SubscriptionsService) Delete(id int) error {
    // Make request to delete /subscriptions/{id}.json
}

// Update updates a subscription
func (s *SubscriptionsService) Update(id int, title string) (*models.Subscription, error) {
    // Make request to patch /subscriptions/{id}.json with title
}
```

#### 3.2 Entries

```go
// EntriesService handles entry-related operations
type EntriesService struct {
    client *Client
}

// ListOptions represents options for listing entries
type ListEntriesOptions struct {
    Page             *int
    Since            *time.Time
    IDs              []int
    Read             *bool
    Starred          *bool
    PerPage          *int
    Mode             string
    IncludeOriginal  bool
    IncludeEnclosure bool
    IncludeContentDiff bool
}

// List returns entries based on the provided options
func (s *EntriesService) List(opts *ListEntriesOptions) ([]models.Entry, *PaginationLinks, error) {
    // Make request to /entries.json with options
}

// ListByFeed returns entries for a specific feed
func (s *EntriesService) ListByFeed(feedID int, opts *ListEntriesOptions) ([]models.Entry, *PaginationLinks, error) {
    // Make request to /feeds/{id}/entries.json with options
}

// Get returns a specific entry
func (s *EntriesService) Get(id int, opts *ListEntriesOptions) (*models.Entry, error) {
    // Make request to /entries/{id}.json with options
}
```

#### 3.3 Unread Entries

```go
// UnreadEntriesService handles unread entry operations
type UnreadEntriesService struct {
    client *Client
}

// List returns all unread entry IDs
func (s *UnreadEntriesService) List() ([]int, error) {
    // Make request to /unread_entries.json
}

// MarkAsUnread marks entries as unread
func (s *UnreadEntriesService) MarkAsUnread(ids []int) ([]int, error) {
    // Make request to post /unread_entries.json with ids
}

// MarkAsRead marks entries as read
func (s *UnreadEntriesService) MarkAsRead(ids []int) ([]int, error) {
    // Make request to delete /unread_entries.json with ids
}
```

#### 3.4 Starred Entries

```go
// StarredEntriesService handles starred entry operations
type StarredEntriesService struct {
    client *Client
}

// List returns all starred entry IDs
func (s *StarredEntriesService) List() ([]int, error) {
    // Make request to /starred_entries.json
}

// Star stars entries
func (s *StarredEntriesService) Star(ids []int) ([]int, error) {
    // Make request to post /starred_entries.json with ids
}

// Unstar unstars entries
func (s *StarredEntriesService) Unstar(ids []int) ([]int, error) {
    // Make request to delete /starred_entries.json with ids
}
```

### 4. Implementation Approach

1. **Core Client**: Implement the base client with authentication and HTTP handling
2. **Models**: Define data structures for API objects
3. **Error Types**: Create custom error types for API responses
4. **Endpoints**: Implement each API endpoint, starting with core functionality:
   - Authentication
   - Subscriptions
   - Entries
   - Unread Entries
   - Starred Entries
5. **Pagination**: Add pagination support to relevant endpoints
6. **Caching**: Implement HTTP caching support
7. **Testing**: Add examples and tests for key functionality

### 5. API Coverage

The client will support all endpoints documented in the Feedbin API:

- Authentication
- Subscriptions
- Entries
- Unread Entries
- Starred Entries
- Taggings
- Tags
- Saved Searches
- Recently Read Entries
- Updated Entries
- Icons
- Imports
- Pages

### 6. Usage Examples

The client will include usage examples for common operations:

```go
// Create a new client
client := feedbin.NewClient("user@example.com", "password")

// Verify authentication
valid, err := client.Authentication.Verify()
if err != nil {
    log.Fatalf("Error verifying credentials: %v", err)
}
if !valid {
    log.Fatal("Invalid credentials")
}

// List subscriptions
subscriptions, err := client.Subscriptions.List(nil, false)
if err != nil {
    log.Fatalf("Error listing subscriptions: %v", err)
}
for _, sub := range subscriptions {
    fmt.Printf("Subscription: %s (%d)\n", sub.Title, sub.ID)
}

// Get entries
opts := &feedbin.ListEntriesOptions{
    PerPage: feedbin.Int(50),
    Starred: feedbin.Bool(true),
}
entries, pagination, err := client.Entries.List(opts)
if err != nil {
    log.Fatalf("Error listing entries: %v", err)
}
fmt.Printf("Found %d entries\n", len(entries))

// Mark entries as read
ids := []int{1234, 5678}
_, err = client.UnreadEntries.MarkAsRead(ids)
if err != nil {
    log.Fatalf("Error marking entries as read: %v", err)
}
```

## Implementation Timeline

1. Core client and authentication (1-2 hours)
2. Models and error types (1-2 hours)
3. Basic endpoints (subscriptions, entries) (2-3 hours)
4. Additional endpoints (unread, starred, tags, etc.) (3-4 hours)
5. Pagination and caching support (1-2 hours)
6. Documentation and examples (1-2 hours)

Total estimated time: 9-15 hours