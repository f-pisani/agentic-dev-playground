# Feedbin API Go Client

This package provides an idiomatic Go client for interacting with the [Feedbin REST API](https://github.com/feedbin/feedbin-api).

## Implementation Plan

Based on the Feedbin API documentation, here's the plan for implementing the client:

### 1. Package Structure

```
feedbin-api/
├── client.go       # Main client implementation and authentication
├── subscriptions.go # Subscription-related API methods  
├── entries.go      # Entry-related API methods
├── unread.go       # Unread entries API methods
├── starred.go      # Starred entries API methods
├── taggings.go     # Tagging-related API methods
├── tags.go         # Tags-related API methods
├── models.go       # Data models for the API
├── pagination.go   # Pagination handling
├── errors.go       # Custom error types
├── examples/       # Usage examples
```

### 2. Core Components

1. **Authentication**: 
   - HTTP Basic Auth with username and password
   - Support for testing credentials

2. **Client Configuration**:
   - Custom HTTP client option
   - Configurable base URL
   - Timeout settings

3. **Models**:
   - Subscriptions
   - Entries
   - Tags
   - Taggings
   - Saved Searches

4. **Features**:
   - Pagination handling with link header parsing
   - Error handling for API responses
   - Proper date/time handling (ISO 8601 format)
   - ETag and If-Modified-Since support for HTTP caching

### 3. Implementation Order

1. Base client with authentication
2. Subscriptions API
3. Entries API
4. Unread entries API
5. Starred entries API
6. Tags and Taggings API
7. Saved Searches API
8. Additional endpoints (Icons, Imports, etc.)
9. Example usage

### 4. Error Handling Strategy

- Custom error types for different API errors
- Proper parsing of HTTP status codes
- Inclusion of response body in error messages when available

### 5. Testing Strategy

- Unit tests for model parsing
- Mock HTTP client for API tests
- Integration examples (not automated tests)

## Usage Example

```go
// Initialize client
client := feedbin.NewClient("username", "password")

// Get subscriptions
subs, err := client.GetSubscriptions(nil)
if err != nil {
    log.Fatal(err)
}

// Print subscription titles
for _, sub := range subs {
    fmt.Println(sub.Title)
}
```

## API Endpoints Implemented

- [x] Authentication
- [ ] Subscriptions
- [ ] Entries
- [ ] Unread Entries
- [ ] Starred Entries
- [ ] Taggings
- [ ] Tags
- [ ] Saved Searches
- [ ] Recently Read Entries
- [ ] Updated Entries
- [ ] Icons
- [ ] Imports
- [ ] Pages
