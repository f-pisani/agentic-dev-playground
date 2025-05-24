# Feedbin API Go Client

A comprehensive, idiomatic Go client for the Feedbin API v2. This client provides full coverage of all Feedbin API endpoints with proper error handling, pagination support, HTTP caching, and follows Go best practices.

## Features

- **Complete API Coverage**: All Feedbin API v2 endpoints supported
- **HTTP Basic Authentication**: Secure credential management
- **Pagination Support**: Automatic pagination handling with iterators
- **HTTP Caching**: ETag and Last-Modified header support
- **Context Support**: All operations support context.Context for cancellation and timeouts
- **Error Handling**: Comprehensive error types with detailed information
- **Standard Library Only**: No external dependencies
- **Idiomatic Go**: Follows Go conventions and best practices
- **Well Documented**: Comprehensive documentation and examples

## Installation

```bash
go get github.com/your-org/feedbin-api/vscode-cline-claude4
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/your-org/feedbin-api/vscode-cline-claude4"
)

func main() {
    // Create a new client
    client := feedbin.NewClient(&feedbin.Config{
        Username: "your-email@example.com",
        Password: "your-password",
    })
    
    // Validate credentials
    ctx := context.Background()
    if err := client.ValidateCredentials(ctx); err != nil {
        log.Fatal("Invalid credentials:", err)
    }
    
    // Get subscriptions
    subscriptions, _, err := client.GetSubscriptions(ctx, nil)
    if err != nil {
        log.Fatal("Failed to get subscriptions:", err)
    }
    
    fmt.Printf("Found %d subscriptions\n", len(subscriptions))
}
```

## API Coverage

### Core Services

- **Authentication**: Credential validation
- **Subscriptions**: Create, read, update, delete subscriptions
- **Entries**: Retrieve and filter entries with extensive options
- **Unread Entries**: Mark entries as read/unread
- **Starred Entries**: Star/unstar entries
- **Tags**: Rename and delete tags
- **Taggings**: Create and manage feed tags
- **Saved Searches**: Manage saved searches
- **Recently Read**: Track recently read entries
- **Updated Entries**: Track entry updates
- **Icons**: Manage feed icons
- **Imports**: OPML import operations
- **Pages**: Pages API support

### Advanced Features

- **Pagination**: Automatic handling with `EntryIterator`
- **Bulk Operations**: Efficient batch processing
- **HTTP Caching**: ETag and Last-Modified support
- **Rate Limiting**: Built-in rate limit awareness
- **Retry Logic**: Automatic retry with exponential backoff

## Examples

### Working with Entries

```go
// Get entries with filtering
since := time.Now().AddDate(0, 0, -7) // Last 7 days
entries, pagination, err := client.GetEntries(ctx, &feedbin.EntryOptions{
    Since:   &since,
    PerPage: &[]int{50}[0],
    Read:    &[]bool{false}[0], // Unread only
})

// Iterate through all pages
iterator := client.NewEntryIterator(ctx, &feedbin.EntryOptions{
    Read: &[]bool{false}[0], // Unread only
})

for iterator.HasMore() {
    entries, err := iterator.Next(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    for _, entry := range entries {
        fmt.Printf("Entry: %s\n", *entry.Title)
    }
}
```

### Managing Subscriptions

```go
// Create a subscription
subscription, err := client.CreateSubscription(ctx, "https://example.com/feed.xml")
if err != nil {
    // Handle multiple choices error
    if multiErr, ok := err.(*feedbin.MultipleChoicesError); ok {
        fmt.Println("Multiple feeds found:")
        for _, choice := range multiErr.GetChoices() {
            fmt.Printf("- %s (%s)\n", choice.Title, choice.FeedURL)
        }
    } else {
        log.Fatal(err)
    }
} else {
    fmt.Printf("Created: %s\n", subscription.Title)
}

// Update subscription title
updated, err := client.UpdateSubscription(ctx, subscription.ID, "Custom Title")
if err != nil {
    log.Fatal(err)
}

// Delete subscription
err = client.DeleteSubscription(ctx, subscription.ID)
if err != nil {
    log.Fatal(err)
}
```

### Bulk Operations

```go
// Get unread entry IDs
unreadIDs, err := client.GetUnreadEntries(ctx)
if err != nil {
    log.Fatal(err)
}

// Mark multiple entries as read
markedRead, err := client.MarkAsRead(ctx, unreadIDs[:10])
if err != nil {
    log.Fatal(err)
}

// Star multiple entries
starred, err := client.StarEntries(ctx, []int{1, 2, 3, 4, 5})
if err != nil {
    log.Fatal(err)
}
```

### Working with Tags

```go
// Get all taggings
taggings, err := client.GetTaggings(ctx)
if err != nil {
    log.Fatal(err)
}

// Create a new tagging
tagging, err := client.CreateTagging(ctx, feedID, "Technology")
if err != nil {
    log.Fatal(err)
}

// Rename a tag
updatedTaggings, err := client.RenameTag(ctx, "Old Name", "New Name")
if err != nil {
    log.Fatal(err)
}

// Delete a tag
remainingTaggings, err := client.DeleteTag(ctx, "Tag Name")
if err != nil {
    log.Fatal(err)
}
```

## Error Handling

The client provides detailed error information:

```go
entries, err := client.GetEntries(ctx, nil)
if err != nil {
    switch e := err.(type) {
    case *feedbin.APIError:
        fmt.Printf("API Error: %d - %s\n", e.StatusCode, e.Message)
        if e.Retryable {
            // Implement retry logic
        }
    case *feedbin.ValidationError:
        fmt.Printf("Validation Error: %s - %s\n", e.Field, e.Message)
    case *feedbin.MultipleChoicesError:
        fmt.Printf("Multiple feeds found: %d choices\n", len(e.GetChoices()))
    default:
        fmt.Printf("Unknown error: %v\n", err)
    }
}
```

## Configuration

```go
config := &feedbin.Config{
    Username:   "your-email@example.com",
    Password:   "your-password",
    BaseURL:    "https://api.feedbin.com/v2/", // Optional, uses default
    HTTPClient: &http.Client{                  // Optional, custom HTTP client
        Timeout: 30 * time.Second,
    },
    UserAgent:  "MyApp/1.0",                   // Optional, custom user agent
    EnableCache: true,                         // Optional, enable HTTP caching
}

client := feedbin.NewClient(config)
```

## HTTP Caching

The client supports HTTP caching to improve performance:

```go
// Enable caching in config
config := &feedbin.Config{
    Username:    "your-email@example.com",
    Password:    "your-password",
    EnableCache: true,
}

client := feedbin.NewClient(config)

// Subsequent requests will use ETag/Last-Modified headers
// and return cached data when appropriate (304 Not Modified)
```

## Best Practices

1. **Use Context**: Always pass a context with appropriate timeout
2. **Handle Errors**: Check for specific error types and handle appropriately
3. **Respect Rate Limits**: The client includes built-in rate limiting awareness
4. **Use Pagination**: For large datasets, use pagination iterators
5. **Enable Caching**: Use HTTP caching for better performance
6. **Batch Operations**: Use bulk operations for efficiency

## Implementation Details

### Architecture

The client uses a service-based architecture where each major API area has its own service. This provides:

- **Modularity**: Easy to test and maintain individual services
- **Separation of Concerns**: Each service handles its specific domain
- **Interface-Based Design**: Easy to mock for testing

### HTTP Features

- **Automatic Retries**: Exponential backoff for transient failures
- **Connection Pooling**: Efficient HTTP connection reuse
- **Request/Response Middleware**: Extensible request/response processing
- **Comprehensive Logging**: Detailed logging for debugging

### Performance

- **Minimal Allocations**: Optimized for low memory usage
- **Efficient JSON Parsing**: Streaming JSON parsing where possible
- **Connection Reuse**: HTTP/1.1 keep-alive support
- **Concurrent Safe**: Thread-safe client operations

## API Reference

### Client Methods

#### Authentication
- `ValidateCredentials(ctx context.Context) error`

#### Subscriptions
- `GetSubscriptions(ctx context.Context, opts *SubscriptionOptions) ([]Subscription, *PaginationInfo, error)`
- `GetSubscription(ctx context.Context, id int) (*Subscription, error)`
- `CreateSubscription(ctx context.Context, feedURL string) (*Subscription, error)`
- `UpdateSubscription(ctx context.Context, id int, title string) (*Subscription, error)`
- `DeleteSubscription(ctx context.Context, id int) error`

#### Entries
- `GetEntries(ctx context.Context, opts *EntryOptions) ([]Entry, *PaginationInfo, error)`
- `GetEntry(ctx context.Context, id int) (*Entry, error)`
- `GetFeedEntries(ctx context.Context, feedID int, opts *EntryOptions) ([]Entry, *PaginationInfo, error)`
- `GetEntriesByIDs(ctx context.Context, ids []int) ([]Entry, error)`
- `NewEntryIterator(ctx context.Context, opts *EntryOptions) *EntryIterator`

#### Unread Entries
- `GetUnreadEntries(ctx context.Context) ([]int, error)`
- `MarkAsRead(ctx context.Context, entryIDs []int) ([]int, error)`
- `MarkAsUnread(ctx context.Context, entryIDs []int) ([]int, error)`
- `MarkAllAsRead(ctx context.Context) ([]int, error)`

#### Starred Entries
- `GetStarredEntries(ctx context.Context) ([]int, error)`
- `StarEntries(ctx context.Context, entryIDs []int) ([]int, error)`
- `UnstarEntries(ctx context.Context, entryIDs []int) ([]int, error)`
- `UnstarAllEntries(ctx context.Context) ([]int, error)`

#### Tags and Taggings
- `GetTaggings(ctx context.Context) ([]Tagging, error)`
- `CreateTagging(ctx context.Context, feedID int, name string) (*Tagging, error)`
- `DeleteTagging(ctx context.Context, id int) error`
- `RenameTag(ctx context.Context, oldName, newName string) ([]Tagging, error)`
- `DeleteTag(ctx context.Context, name string) ([]Tagging, error)`

#### Saved Searches
- `GetSavedSearches(ctx context.Context) ([]SavedSearch, error)`
- `CreateSavedSearch(ctx context.Context, name, query string) (*SavedSearch, error)`
- `UpdateSavedSearch(ctx context.Context, id int, name, query string) (*SavedSearch, error)`
- `DeleteSavedSearch(ctx context.Context, id int) error`

#### Other APIs
- `GetRecentlyReadEntries(ctx context.Context) ([]int, error)`
- `GetUpdatedEntries(ctx context.Context) ([]int, error)`
- `GetIcons(ctx context.Context) ([]Icon, error)`
- `GetImports(ctx context.Context) ([]Import, error)`
- `GetPages(ctx context.Context) ([]Page, error)`

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass with `go test ./...`
5. Run `go vet` and `gofmt`
6. Submit a pull request

## License

MIT License - see LICENSE file for details.
