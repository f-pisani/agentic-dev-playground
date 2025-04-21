# Feedbin API Go Client

This is a Go client library for the [Feedbin REST API v2](https://github.com/feedbin/feedbin-api).

## Installation

```shell
go get github.com/feedbin/go-client
```

## Usage Examples

### Authentication

```go
import "github.com/feedbin/go-client"

// Create a new client
client := feedbin.NewClient("user@example.com", "password")

// Verify credentials
valid, err := client.Authentication.Validate()
if err != nil {
    // Handle error
}

if valid {
    // Credentials are valid
} else {
    // Credentials are invalid
}
```

### Subscriptions

```go
// List all subscriptions
subscriptions, resp, err := client.Subscriptions.List(nil)

// Get subscription by ID
subscription, resp, err := client.Subscriptions.Get(123)

// Create a new subscription
subscription, resp, err := client.Subscriptions.Create("https://example.com/feed.xml")

// Update a subscription
subscription, resp, err := client.Subscriptions.Update(123, "New Title")

// Delete a subscription
resp, err := client.Subscriptions.Delete(123)
```

### Entries

```go
// List all entries
entries, resp, err := client.Entries.List(nil)

// Get entries with options (pagination, filters)
opts := &feedbin.EntryOptions{
    PageOptions: feedbin.PageOptions{
        Page: feedbin.Int(2),
        PerPage: feedbin.Int(50),
    },
    Read: feedbin.Bool(false),         // Only unread entries
    Starred: feedbin.Bool(true),       // Only starred entries
    Mode: "extended",                  // Include extended metadata
    IncludeEnclosure: feedbin.Bool(true), // Include podcast data
}
entries, resp, err := client.Entries.List(opts)

// Get entries for a specific feed
entries, resp, err := client.Entries.ListByFeed(123, nil)

// Get a specific entry
entry, resp, err := client.Entries.Get(456)

// Get entries by IDs
entries, resp, err := client.Entries.GetByIDs([]int{1, 2, 3})
```

### Unread Entries

```go
// Get all unread entry IDs
ids, resp, err := client.Unread.List()

// Mark entries as unread
markedIds, resp, err := client.Unread.MarkAsUnread([]int{1, 2, 3})

// Mark entries as read
markedIds, resp, err := client.Unread.MarkAsRead([]int{1, 2, 3})
```

### Starred Entries

```go
// Get all starred entry IDs
ids, resp, err := client.Starred.List()

// Star entries
starredIds, resp, err := client.Starred.Star([]int{1, 2, 3})

// Unstar entries
unstarredIds, resp, err := client.Starred.Unstar([]int{1, 2, 3})
```

### Tags and Taggings

```go
// Get all tags
tags, resp, err := client.Tags.List()

// Get all taggings
taggings, resp, err := client.Taggings.List()

// Add a tag to a feed
tagging, resp, err := client.Taggings.Create(123, 0, "New Tag")

// Add a feed to an existing tag
tagging, resp, err := client.Taggings.Create(123, 456, "")

// Remove a tagging
resp, err := client.Taggings.Delete(789)
```

### Saved Searches

```go
// Get all saved searches
searches, resp, err := client.SavedSearches.List()

// Get a specific saved search
search, resp, err := client.SavedSearches.Get(123)

// Create a new saved search
search, resp, err := client.SavedSearches.Create("Tech News", "tech news")

// Update a saved search
search, resp, err := client.SavedSearches.Update(123, "New Title", "new query")

// Delete a saved search
resp, err := client.SavedSearches.Delete(123)
```

## Handling Pagination

The Feedbin API uses pagination for endpoints that return lists of items. The pagination information is included in the response headers:

```go
// Get entries
opts := &feedbin.EntryOptions{
    PageOptions: feedbin.PageOptions{
        Page: feedbin.Int(1),
    },
}
entries, resp, err := client.Entries.List(opts)

// Parse pagination information
linkHeader, err := feedbin.ParseLinkHeader(resp)
if err != nil {
    // Handle error
}

// Check if there are more pages
if linkHeader != nil && linkHeader.NextURL != nil {
    // There are more pages
    fmt.Printf("Total records: %d\n", linkHeader.Count)
}
```

## Error Handling

Errors are returned as `error` values. For API errors, you can check if the error is of type `*feedbin.ErrorResponse` to get more information:

```go
_, resp, err := client.Subscriptions.Get(123)
if err != nil {
    if errorResponse, ok := err.(*feedbin.ErrorResponse); ok {
        // This is an API error
        statusCode := errorResponse.Response.StatusCode
        message := errorResponse.Message
        
        if statusCode == 404 {
            // Handle not found error
        } else if statusCode == 403 {
            // Handle forbidden error
        }
    } else {
        // This is some other kind of error (network, JSON parsing, etc.)
    }
}
```

## API Implementation Status

- ✅ Authentication
- ✅ Subscriptions
- ✅ Entries
- ✅ Unread Entries
- ✅ Starred Entries
- ✅ Tags
- ✅ Taggings
- ✅ Saved Searches

## License

This library is provided under the MIT License. 