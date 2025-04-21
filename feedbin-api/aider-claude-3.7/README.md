# Feedbin API Go Client

This is a Go client library for the Feedbin REST API V2. It provides a simple, idiomatic way to interact with the Feedbin API from Go applications.

## Implementation Plan

### 1. Package Structure

```
feedbin-api/aider-claude-3.7/
├── client.go       # Main client implementation
├── auth.go         # Authentication handling
├── models.go       # Data models for API objects
├── subscriptions.go # Subscriptions API endpoints
├── entries.go      # Entries API endpoints
├── unread.go       # Unread entries API endpoints
├── starred.go      # Starred entries API endpoints
├── taggings.go     # Taggings API endpoints
├── tags.go         # Tags API endpoints
├── searches.go     # Saved searches API endpoints
├── recently_read.go # Recently read entries API endpoints
├── updated.go      # Updated entries API endpoints
├── icons.go        # Icons API endpoints
├── imports.go      # Imports API endpoints
├── pages.go        # Pages API endpoints
└── examples/       # Example usage
```

### 2. Core Components

1. **Client**: The main client struct that handles HTTP requests, authentication, and provides methods for each API endpoint.
2. **Authentication**: HTTP Basic authentication as specified in the API docs.
3. **Models**: Go structs representing the JSON objects returned by the API.
4. **Pagination**: Helper functions to handle pagination for endpoints that support it.
5. **Error Handling**: Proper error handling for non-2xx HTTP responses.

### 3. Implementation Approach

1. Create a base client that handles authentication and common HTTP operations.
2. Implement models for each API object.
3. Implement methods for each API endpoint, organized by resource type.
4. Add pagination support for endpoints that require it.
5. Implement proper error handling and logging.
6. Create examples demonstrating usage of the client.

### 4. API Coverage

The client will support all endpoints documented in the Feedbin API V2 specs:

- Authentication
- Subscriptions (GET, POST, DELETE, PATCH)
- Entries (GET)
- Unread Entries (GET, POST, DELETE)
- Starred Entries (GET, POST, DELETE)
- Taggings (GET, POST, DELETE)
- Tags (GET, DELETE)
- Saved Searches (GET, POST, DELETE, UPDATE)
- Recently Read Entries (GET, POST)
- Updated Entries (GET)
- Icons (GET)
- Imports (GET, POST)
- Pages (GET)

### 5. Usage Example

```go
package main

import (
    "fmt"
    "log"
    
    feedbin "github.com/yourusername/feedbin-api/aider-claude-3.7"
)

func main() {
    client := feedbin.NewClient("your-email@example.com", "your-password")
    
    // Get all subscriptions
    subscriptions, err := client.GetSubscriptions()
    if err != nil {
        log.Fatalf("Error getting subscriptions: %v", err)
    }
    
    for _, sub := range subscriptions {
        fmt.Printf("Feed: %s (%s)\n", sub.Title, sub.FeedURL)
    }
}
```

## Implementation Notes

- Using only the Go standard library as required
- Following Go best practices for API client design
- Implementing proper error handling
- Supporting pagination for relevant endpoints
- Handling dates in ISO 8601 format as specified in the API docs
