# Feedbin API Client

A Go client library for the Feedbin REST API v2, implemented using only the Go standard library.

## Implementation Plan

### Overview
This package provides a complete Go client for the Feedbin API, supporting all documented endpoints with proper authentication, pagination, and error handling.

### Architecture

#### Core Components
- **Client**: Main client struct with HTTP configuration and authentication
- **Models**: Type-safe representations of all API responses
- **Authentication**: HTTP Basic Auth with email/password credentials
- **Pagination**: Link header parsing for paginated responses
- **Error Handling**: Structured error types for different API error scenarios

#### Package Structure
```
feedbin/
├── client.go          # Main client implementation
├── models.go          # API response structs
├── auth.go            # Authentication methods
├── subscriptions.go   # Subscription operations
├── entries.go         # Entry operations
├── tags.go            # Tag and tagging operations
├── pagination.go      # Pagination utilities
├── errors.go          # Error handling
└── examples/          # Usage examples
```

### Supported Operations

#### Authentication
- `Authenticate()` - Verify user credentials

#### Subscriptions
- `GetSubscriptions(opts)` - List all subscriptions with optional filtering
- `GetSubscription(id)` - Get specific subscription
- `CreateSubscription(feedURL)` - Subscribe to a feed
- `UpdateSubscription(id, title)` - Update subscription title
- `DeleteSubscription(id)` - Delete subscription

#### Entries
- `GetEntries(opts)` - Get paginated entries with filtering
- `GetFeedEntries(feedID, opts)` - Get entries for specific feed
- `GetEntry(id, opts)` - Get specific entry with optional extended mode

#### Unread Entries
- `GetUnreadEntries()` - Get list of unread entry IDs
- `MarkAsUnread(entryIDs)` - Mark entries as unread
- `MarkAsRead(entryIDs)` - Mark entries as read

#### Starred Entries
- `GetStarredEntries()` - Get list of starred entry IDs
- `StarEntries(entryIDs)` - Star entries
- `UnstarEntries(entryIDs)` - Unstar entries

#### Tags & Taggings
- `GetTaggings()` - Get all taggings
- `GetTagging(id)` - Get specific tagging
- `CreateTagging(feedID, name)` - Create new tagging
- `DeleteTagging(id)` - Delete tagging
- `RenameTag(oldName, newName)` - Rename tag
- `DeleteTag(name)` - Delete tag

### Key Features

#### Pagination Support
- Automatic Link header parsing
- `PaginatedResponse` wrapper for paginated results
- Support for `page`, `per_page`, and navigation links

#### Flexible Options
- Option structs for all list operations
- Support for `since`, `mode`, `read`, `starred` filters
- Extended mode support for detailed metadata

#### Error Handling
- Custom error types for different HTTP status codes
- Detailed error messages with API response context
- Proper handling of 401, 403, 404, 415, and other status codes

#### HTTP Caching Support
- ETag and Last-Modified header handling
- Built-in support for conditional requests

### Usage Example

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/user/feedbin"
)

func main() {
    client := feedbin.NewClient("user@example.com", "password")
    
    // Verify authentication
    if err := client.Authenticate(); err != nil {
        log.Fatal("Authentication failed:", err)
    }
    
    // Get subscriptions
    subs, err := client.GetSubscriptions(nil)
    if err != nil {
        log.Fatal("Failed to get subscriptions:", err)
    }
    
    fmt.Printf("Found %d subscriptions\n", len(subs))
    
    // Get recent entries
    opts := &feedbin.EntryOptions{
        PerPage: 20,
        Read:    feedbin.Bool(false), // Only unread
    }
    
    entries, pagination, err := client.GetEntries(opts)
    if err != nil {
        log.Fatal("Failed to get entries:", err)
    }
    
    fmt.Printf("Found %d unread entries\n", len(entries))
    
    // Handle pagination
    if pagination.HasNext() {
        nextEntries, _, err := client.GetEntriesFromURL(pagination.Next)
        if err == nil {
            fmt.Printf("Next page has %d entries\n", len(nextEntries))
        }
    }
}
```

### Implementation Notes

#### Standard Library Only
- Uses `net/http` for HTTP requests
- Uses `encoding/json` for JSON handling
- Uses `net/url` for URL construction and parsing
- No external dependencies required

#### Idiomatic Go
- Proper error handling with custom error types
- Pointer fields for optional values in option structs
- Context support for request cancellation
- Interface-based design for testability

#### Thread Safety
- Client is safe for concurrent use
- HTTP client reuse with proper configuration
- No shared mutable state

### Testing Strategy
- Unit tests for all public methods
- Mock HTTP server for integration tests
- Examples that double as documentation tests
- Error case coverage for all status codes

## Implementation Status

✅ **Completed Features:**
- HTTP Basic Authentication with email/password
- Complete CRUD operations for subscriptions
- Full entry retrieval with filtering and pagination
- Unread/read status management
- Starred entries management
- Tag and tagging operations
- Comprehensive error handling with custom error types
- Pagination support with Link header parsing
- Request/response validation
- Context support for cancellation
- Standard library only implementation

✅ **Code Quality:**
- Passes `go vet` without warnings
- All tests passing
- Idiomatic Go patterns
- Thread-safe client design
- Proper error wrapping and handling

## Quick Start

```bash
# Install the package
go get github.com/feedbin/feedbin-go

# Set environment variables (for example)
export FEEDBIN_EMAIL="your-email@example.com"
export FEEDBIN_PASSWORD="your-password"

# Run the example
cd examples
go run main.go
```

## Package Documentation

The client provides both context-aware methods and convenience methods without context:

**Context Methods (Recommended):**
- `client.GetSubscriptions(ctx, opts)`
- `client.Authenticate(ctx)`
- All CRUD operations accept context

**Convenience Methods:**
- `client.GetSubscriptionsWithoutContext(opts)`
- `client.AuthenticateWithoutContext()`
- Use `context.Background()` internally

## Error Handling

The client provides structured error handling:

```go
if err != nil {
    if apiErr, ok := err.(*feedbin.APIError); ok {
        switch {
        case apiErr.IsUnauthorized():
            // Handle authentication failure
        case apiErr.IsNotFound():
            // Handle resource not found
        case apiErr.IsMultipleChoices():
            // Handle multiple feed choices
        default:
            // Handle other API errors
        }
    } else {
        // Handle network or other errors
    }
}
```

## Notes

- All timestamps use RFC3339 format as required by the API
- Entry content may contain null values for title, author, and content
- Pagination is automatically handled via Link headers
- Rate limiting should be implemented by the caller if needed
- HTTPS is enforced (HTTP will fail)

## Contributing

This implementation covers the complete Feedbin API v2 specification. For additional features or bug reports, please refer to the official API documentation in the `specs/` directory.