# Feedbin API Client for Go

This package provides an idiomatic Go client for the Feedbin REST API (v2). It allows Go applications to interact with Feedbin's services using a clean, type-safe interface.

## Implementation Plan

### 1. Package Structure

```
feedbin-api/windsurf-cascade-claude-3.7/
├── README.md                 # This file
├── client.go                 # Core client implementation
├── auth.go                   # Authentication functionality
├── models.go                 # Data models/structs
├── subscriptions.go          # Subscriptions API
├── entries.go                # Entries API
├── unread_entries.go         # Unread entries API
├── starred_entries.go        # Starred entries API
├── taggings.go               # Taggings API
├── tags.go                   # Tags API
├── saved_searches.go         # Saved searches API
├── recently_read_entries.go  # Recently read entries API
├── updated_entries.go        # Updated entries API
├── icons.go                  # Icons API
├── imports.go                # Imports API
├── pages.go                  # Pages API
└── examples/                 # Example usage
```

### 2. Core Components

#### Client

The core client will handle:
- HTTP request creation and execution
- Authentication
- Base URL configuration
- HTTP caching (ETag and Last-Modified support)
- Error handling

#### Models

Define Go structs that map to Feedbin API JSON responses, including:
- Subscription
- Entry
- Tagging
- SavedSearch
- Icon
- Import
- Page

### 3. API Implementation Approach

Each API endpoint will be implemented as methods on the client, following these patterns:

1. **GET Collection**: Return a slice of resources
2. **GET Single Resource**: Return a single resource
3. **POST Create**: Accept a resource or parameters and return the created resource
4. **DELETE Resource**: Accept an ID and return success/error
5. **PATCH/PUT Update**: Accept an ID and updated resource, return updated resource

### 4. Pagination Support

Implement pagination helpers for endpoints that support it, with options for:
- Page-based navigation
- Automatic iteration through all pages

### 5. Error Handling

Create a consistent error handling approach:
- Custom error types for API-specific errors
- HTTP status code mapping to meaningful errors
- Detailed error messages from the API

### 6. Authentication

Support HTTP Basic Authentication with:
- Username/password authentication
- Credential validation

### 7. Implementation Order

1. Core client and authentication
2. Models
3. Subscriptions API
4. Entries API
5. Unread/Starred entries API
6. Taggings and Tags API
7. Saved Searches API
8. Recently Read/Updated Entries API
9. Icons API
10. Imports API
11. Pages API

## Usage Example

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/yourusername/feedbin"
)

func main() {
    // Create a new client
    client := feedbin.NewClient("username", "password")
    
    // Get subscriptions
    subs, err := client.GetSubscriptions(nil)
    if err != nil {
        log.Fatal(err)
    }
    
    // Print subscription titles
    for _, sub := range subs {
        fmt.Printf("Feed: %s\n", sub.Title)
    }
}
```

## Design Decisions

1. **Standard Library Only**: Using only the Go standard library for HTTP requests and JSON parsing.
2. **Idiomatic Go**: Following Go conventions for naming, error handling, and package structure.
3. **Type Safety**: Using strong typing for all API responses and requests.
4. **Configurability**: Allowing customization of client behavior through options.
5. **Testability**: Designing the client to be easily testable.
