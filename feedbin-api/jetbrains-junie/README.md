# Feedbin API Client for Go

This is a Go client library for the Feedbin API V2. It provides a simple, idiomatic way to interact with the Feedbin API from Go applications.

## Implementation Plan

### 1. Package Structure

```
jetbrains-junie/
├── client.go       # Main client implementation
├── auth.go         # Authentication handling
├── models.go       # Data models/structs
├── subscriptions.go # Subscriptions API
├── entries.go      # Entries API
├── unread.go       # Unread entries API
├── starred.go      # Starred entries API
├── taggings.go     # Taggings API
├── tags.go         # Tags API
├── saved_searches.go # Saved searches API
├── recently_read.go # Recently read entries API
├── updated.go      # Updated entries API
├── icons.go        # Icons API
├── imports.go      # Imports API
├── pages.go        # Pages API
├── pagination.go   # Pagination utilities
├── errors.go       # Error handling
└── examples/       # Example usage
```

### 2. Core Components

#### 2.1 Client

The main client will handle:
- HTTP request creation and execution
- Authentication
- Base URL configuration
- HTTP caching support (ETag and Last-Modified)
- Error handling and response parsing

#### 2.2 Authentication

Implement HTTP Basic authentication as required by the Feedbin API.

#### 2.3 Models

Define Go structs for all API objects:
- Subscription
- Entry
- Tag
- Tagging
- SavedSearch
- etc.

#### 2.4 Pagination

Implement pagination support:
- Parse Link headers
- Handle page parameters
- Provide utilities for iterating through paginated results

#### 2.5 Error Handling

Implement error handling:
- HTTP status code checking
- Error response parsing
- Custom error types for different API errors

### 3. API Endpoints

Implement all API endpoints as described in the documentation:

#### 3.1 Authentication
- Verify credentials

#### 3.2 Subscriptions
- Get all subscriptions
- Get a specific subscription
- Create a subscription
- Delete a subscription
- Update a subscription

#### 3.3 Entries
- Get all entries
- Get entries for a feed
- Get a specific entry
- Support for various parameters (since, ids, read, starred, etc.)

#### 3.4 Unread Entries
- Get unread entries
- Mark entries as unread
- Mark entries as read

#### 3.5 Starred Entries
- Get starred entries
- Star entries
- Unstar entries

#### 3.6 Taggings
- Get all taggings
- Create a tagging
- Delete a tagging

#### 3.7 Tags
- Get all tags
- Rename a tag
- Delete a tag

#### 3.8 Saved Searches
- Get all saved searches
- Get a specific saved search
- Create a saved search
- Update a saved search
- Delete a saved search

#### 3.9 Recently Read Entries
- Get recently read entries
- Add recently read entries
- Delete recently read entries

#### 3.10 Updated Entries
- Get updated entries

#### 3.11 Icons
- Get feed icons

#### 3.12 Imports
- Get imports
- Create an import
- Get an import
- Delete an import

#### 3.13 Pages
- Get pages
- Get a specific page
- Create a page
- Update a page
- Delete a page

### 4. Implementation Approach

1. Start with the core client and authentication
2. Implement models for all API objects
3. Add pagination support
4. Implement error handling
5. Implement each API endpoint group
6. Add examples and documentation
7. Test all functionality

### 5. Testing Strategy

- Unit tests for client functionality
- Integration tests for API endpoints (optional)
- Example code to demonstrate usage

## Usage Example

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/example/feedbin"
)

func main() {
    // Create a new client
    client := feedbin.NewClient("username", "password")
    
    // Get all subscriptions
    subscriptions, err := client.GetSubscriptions()
    if err != nil {
        log.Fatal(err)
    }
    
    // Print subscription titles
    for _, sub := range subscriptions {
        fmt.Printf("Subscription: %s\n", sub.Title)
    }
}
```