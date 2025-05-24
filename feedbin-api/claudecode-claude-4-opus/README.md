# Feedbin API Go Client

A comprehensive Go client library for the Feedbin RSS reader API v2, implemented using only the Go standard library.

## Implementation Plan

### Architecture Overview

The client follows a modular design with clear separation of concerns:

1. **Core Client (`client.go`)**: Main client struct handling HTTP requests, authentication, and response processing
2. **Models (`models.go`)**: Data structures representing API entities (Entry, Subscription, Feed, etc.)
3. **Authentication (`auth.go`)**: HTTP Basic Authentication implementation
4. **Error Handling (`errors.go`)**: Custom error types for API-specific errors
5. **Pagination (`pagination.go`)**: Utilities for handling paginated responses
6. **API Endpoints**: Separate files for each resource type (entries.go, subscriptions.go, etc.)

### Key Features

- **Authentication**: HTTP Basic Authentication with email/password
- **Full API Coverage**: All documented endpoints implemented
- **Pagination Support**: Automatic handling of paginated responses
- **Error Handling**: Proper HTTP status code handling with descriptive errors
- **Extended Mode**: Support for additional metadata via `mode=extended`
- **Bulk Operations**: Efficient handling of bulk entry operations (mark read/unread/starred)
- **Content Extraction**: Support for Mercury Parser integration
- **Zero Dependencies**: Uses only Go standard library

### Implementation Details

#### Client Structure
```go
type Client struct {
    BaseURL    string
    Email      string
    Password   string
    HTTPClient *http.Client
}
```

#### Core Methods
- HTTP request builder with authentication
- JSON response parsing with error handling
- Pagination iterator for large result sets
- Query parameter builder for filtering

#### Resource Methods
Each resource type will have dedicated methods:
- **Entries**: Get, Create, Update, Delete, Bulk operations
- **Subscriptions**: List, Create, Update, Delete
- **Feeds**: Get feed information
- **Tags/Taggings**: Organize feeds with tags
- **Saved Searches**: Manage search queries
- **Import/Export**: OPML support
- **Pages**: Save web pages as entries

### Error Handling Strategy

Custom error types for different scenarios:
- `AuthenticationError`: 401 responses
- `NotFoundError`: 404 responses
- `ValidationError`: 400 responses with validation details
- `RateLimitError`: Rate limiting errors
- `ServerError`: 5xx responses

### Usage Example

```go
// Create client
client := feedbin.NewClient("user@example.com", "password")

// Get all subscriptions
subs, err := client.GetSubscriptions()

// Get unread entries
unreadIDs, err := client.GetUnreadEntries()
entries, err := client.GetEntries(feedbin.WithIDs(unreadIDs[:100]))

// Mark entries as read
err = client.MarkEntriesRead([]int{1, 2, 3})

// Create subscription
sub, err := client.CreateSubscription("https://example.com/feed.xml")
```

### File Structure

```
feedbin-api/claudecode-claude-4-opus/
├── README.md           # This file
├── go.mod             # Go module definition
├── client.go          # Core client implementation
├── auth.go            # Authentication helpers
├── models.go          # Data structures
├── errors.go          # Error types
├── pagination.go      # Pagination utilities
├── entries.go         # Entries API
├── subscriptions.go   # Subscriptions API
├── feeds.go           # Feeds API
├── tags.go            # Tags API
├── taggings.go        # Taggings API
├── saved_searches.go  # Saved searches API
├── icons.go           # Icons API
├── imports.go         # Import API
├── pages.go           # Pages API
├── extract.go         # Content extraction
└── example/
    └── main.go        # Example usage
```

### Testing Strategy

While not required, the implementation will be structured to be easily testable:
- Interface-based design for mocking
- Separate HTTP logic from business logic
- Example code demonstrating all major features

### Development Phases

1. **Phase 1**: Core infrastructure (client, auth, models, errors)
2. **Phase 2**: Essential endpoints (entries, subscriptions, feeds)
3. **Phase 3**: Organization features (tags, saved searches)
4. **Phase 4**: Advanced features (import/export, content extraction)
5. **Phase 5**: Polish and examples