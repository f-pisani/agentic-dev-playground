# Feedbin API Go Client

This is an idiomatic Go client for the Feedbin API V2. It provides a simple interface to interact with the Feedbin REST API.

## Features

- HTTP Basic Authentication
- Support for all Feedbin API endpoints
- Pagination handling
- Error handling
- HTTP caching support (ETag and Last-Modified)
- Fully documented code

## Installation

```bash
go get github.com/yourusername/feedbin-go
```

## Package Structure

```
feedbin-api/jetbrains-augmentcode-agent-auto/
├── README.md
├── client.go         # Main client implementation
├── authentication.go # Authentication service
├── subscriptions.go  # Subscriptions service
├── entries.go        # Entries service
├── unread_entries.go # Unread entries service
├── starred_entries.go # Starred entries service
├── taggings.go       # Taggings service
├── tags.go           # Tags service
├── saved_searches.go # Saved searches service
├── updated_entries.go # Updated entries service
├── icons.go          # Icons service
├── imports.go        # Imports service
├── pages.go          # Pages service
├── extract.go        # Full content extraction service
├── models.go         # Data models
└── examples/         # Usage examples
    └── main.go
```

## Implemented API Endpoints

- Authentication
- Subscriptions
- Entries
- Unread Entries
- Starred Entries
- Taggings
- Tags
- Saved Searches
- Updated Entries
- Icons
- Imports
- Pages
- Full Content Extraction

## Usage

```go
package main

import (
    "fmt"
    "log"

    feedbin "github.com/feedbin/go-client"
)

func main() {
    // Create a new client
    client := feedbin.NewClient("username", "password")

    // Check authentication
    if err := client.Authentication.Verify(); err != nil {
        log.Fatalf("Authentication failed: %v", err)
    }

    // Get subscriptions
    subscriptions, _, err := client.Subscriptions.List(nil)
    if err != nil {
        log.Fatalf("Failed to get subscriptions: %v", err)
    }

    // Print subscriptions
    for _, sub := range subscriptions {
        fmt.Printf("Subscription: %s (%s)\n", sub.Title, sub.FeedURL)
    }
}
```

### Authentication

```go
// Create a new client
client := feedbin.NewClient("username", "password")

// Check authentication
if err := client.Authentication.Verify(); err != nil {
    log.Fatalf("Authentication failed: %v", err)
}
```

### Subscriptions

```go
// Get all subscriptions
subscriptions, _, err := client.Subscriptions.List(nil)

// Get subscriptions with options
opts := &feedbin.SubscriptionListOptions{
    Since: time.Now().Add(-24 * time.Hour), // Get subscriptions created in the last 24 hours
    Mode: "extended", // Get extended information
}
subscriptions, resp, err := client.Subscriptions.List(opts)

// Get a single subscription
subscription, _, err := client.Subscriptions.Get(12345, "")

// Create a subscription
newSub, _, err := client.Subscriptions.Create(&feedbin.CreateSubscriptionOptions{
    FeedURL: "https://example.com/feed.xml",
})

// Update a subscription
updatedSub, _, err := client.Subscriptions.Update(12345, &feedbin.UpdateSubscriptionOptions{
    Title: "My Custom Title",
})

// Delete a subscription
_, err := client.Subscriptions.Delete(12345)
```

### Entries

```go
// Get all entries
entries, _, err := client.Entries.List(nil)

// Get entries with options
opts := &feedbin.EntryListOptions{
    Page: 1,
    PerPage: 50,
    Read: feedbin.Bool(false), // Only unread entries
    Starred: feedbin.Bool(true), // Only starred entries
    Mode: "extended", // Get extended information
    IncludeOriginal: true, // Include original entry data if updated
}
entries, resp, err := client.Entries.List(opts)

// Get entries for a specific feed
entries, _, err := client.Entries.ListByFeed(12345, nil)

// Get a single entry
entry, _, err := client.Entries.Get(12345, nil)
```

### Unread Entries

```go
// Get all unread entry IDs
unreadIDs, _, err := client.UnreadEntries.List()

// Mark entries as unread
markedIDs, _, err := client.UnreadEntries.Create([]int64{12345, 12346, 12347})

// Mark entries as read
markedIDs, _, err := client.UnreadEntries.Delete([]int64{12345, 12346, 12347})
```

### Starred Entries

```go
// Get all starred entry IDs
starredIDs, _, err := client.StarredEntries.List()

// Star entries
starredIDs, _, err := client.StarredEntries.Create([]int64{12345, 12346, 12347})

// Unstar entries
unstarredIDs, _, err := client.StarredEntries.Delete([]int64{12345, 12346, 12347})
```

### Full Content Extraction

```go
// Set the extract service secret
client.Extract.SetSecret("your-extract-service-secret")

// Extract full content from a URL
article, err := client.Extract.Extract("https://example.com/article")
if err != nil {
    log.Fatalf("Failed to extract content: %v", err)
}

fmt.Printf("Title: %s\n", article.Title)
fmt.Printf("Content: %s\n", article.Content)
```
