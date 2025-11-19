# Feedbin API Client (Go)

This is an idiomatic Go client for the [Feedbin API](https://github.com/feedbin/feedbin-api), implemented as part of a coding task.

## Implementation Plan

### 1. Package Structure

The package will be self-contained in `feedbin-api/antigravity-gemini-3-pro`.

-   `client.go`: Main client struct, initialization, and common HTTP logic (request creation, auth, error handling).
-   `models.go`: Go structs representing API resources (Subscription, Entry, etc.).
-   `subscriptions.go`: Methods for interacting with the Subscriptions API.
-   `entries.go`: Methods for interacting with the Entries API.
-   `unread_entries.go`: Methods for managing unread entries.
-   `starred_entries.go`: Methods for managing starred entries.
-   `taggings.go`: Methods for managing taggings.
-   `tags.go`: Methods for managing tags.
-   `saved_searches.go`: Methods for managing saved searches.
-   `recently_read.go`: Methods for managing recently read entries.
-   `updated_entries.go`: Methods for managing updated entries.
-   `icons.go`: Methods for fetching icons.
-   `imports.go`: Methods for managing imports.
-   `pages.go`: Methods for creating pages.
-   `authentication.go`: Methods for checking authentication.
-   `client_test.go`: Tests for the client.

### 2. Client Design

The `Client` struct will hold:
-   `BaseURL`: API base URL (defaulting to `https://api.feedbin.com/v2/`).
-   `HTTPClient`: `*http.Client` for making requests.
-   `User`: Username (email) for Basic Auth.
-   `Password`: Password for Basic Auth.

```go
type Client struct {
    BaseURL    string
    HTTPClient *http.Client
    User       string
    Password   string
}
```

### 3. Authentication

The client will use HTTP Basic Authentication for all requests. A helper method `NewRequest` will attach the `Authorization` header.

### 4. Pagination

Pagination will be supported via a `PaginationOptions` struct (for request parameters like `page`, `per_page`, `since`) and by parsing the `Link` header in responses to find next/prev pages.

### 5. Error Handling

Non-2xx responses will be converted into a `APIError` struct containing the status code and message.

### 6. Supported Endpoints

-   **Authentication**: `GET /authentication.json`
-   **Subscriptions**:
    -   `GET /subscriptions.json`
    -   `GET /subscriptions/:id.json`
    -   `POST /subscriptions.json`
    -   `DELETE /subscriptions/:id.json`
    -   `PATCH /subscriptions/:id.json`
-   **Entries**:
    -   `GET /entries.json` (with pagination)
    -   `GET /entries.json?ids=...`
    -   `GET /feeds/:feed_id/entries.json`
    -   `GET /entries/:id.json`
-   **Unread Entries**:
    -   `GET /unread_entries.json`
    -   `POST /unread_entries.json` (mark as unread)
    -   `DELETE /unread_entries.json` (mark as read)
    -   `POST /unread_entries/delete.json` (alternative)
-   **Starred Entries**:
    -   `GET /starred_entries.json`
    -   `POST /starred_entries.json` (star)
    -   `DELETE /starred_entries.json` (unstar)
    -   `POST /starred_entries/delete.json` (alternative)
-   **Taggings**:
    -   `GET /taggings.json`
    -   `GET /taggings/:id.json`
    -   `POST /taggings.json`
    -   `DELETE /taggings/:id.json`
-   **Tags**:
    -   `POST /tags.json` (rename)
    -   `DELETE /tags.json` (delete)
-   **Saved Searches**:
    -   `GET /saved_searches.json`
    -   `GET /saved_searches/:id.json`
    -   `POST /saved_searches.json`
    -   `DELETE /saved_searches/:id.json`
    -   `PATCH /saved_searches/:id.json`
    -   `POST /saved_searches/:id/update.json` (alternative)
-   **Recently Read Entries**:
    -   `GET /recently_read_entries.json`
    -   `POST /recently_read_entries.json`
-   **Updated Entries**:
    -   `GET /updated_entries.json`
    -   `DELETE /updated_entries.json` (mark as read)
    -   `POST /updated_entries/delete.json` (alternative)
-   **Icons**:
    -   `GET /icons.json`
-   **Imports**:
    -   `GET /imports.json`
    -   `GET /imports/:id.json`
    -   `POST /imports.json`
-   **Pages**:
    -   `POST /pages.json`

### 7. Testing

Basic tests will be added to verify struct definitions and request formation.

## Usage

```go
client := feedbin.NewClient("user@example.com", "password")
subs, err := client.GetSubscriptions(context.Background(), nil)
if err != nil {
    log.Fatal(err)
}
for _, sub := range subs {
    fmt.Printf("Sub: %s\n", sub.Title)
}
```
