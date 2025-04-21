package feedbin

import (
	"time"
)

// Subscription represents a Feedbin subscription
type Subscription struct {
	ID              int64     `json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	Title           string    `json:"title"`
	FeedURL         string    `json:"feed_url"`
	SiteURL         string    `json:"site_url"`
	FeedID          int64     `json:"feed_id"`
	UnreadCount     int       `json:"unread_count,omitempty"`
	Favicon         *Favicon  `json:"favicon,omitempty"`
	UpdatedAt       time.Time `json:"updated_at,omitempty"`
	LastPublishedAt time.Time `json:"last_published_entry_at,omitempty"`
}

// SubscriptionRequest represents a request to create a new subscription
type SubscriptionRequest struct {
	FeedURL string `json:"feed_url"`
}

// SubscriptionUpdateRequest represents a request to update a subscription
type SubscriptionUpdateRequest struct {
	Title string `json:"title"`
}

// Entry represents a Feedbin entry
type Entry struct {
	ID            int64     `json:"id"`
	FeedID        int64     `json:"feed_id"`
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	Summary       string    `json:"summary"`
	Content       string    `json:"content"`
	URL           string    `json:"url"`
	ExtractedURL  string    `json:"extracted_content_url,omitempty"`
	Published     time.Time `json:"published"`
	CreatedAt     time.Time `json:"created_at"`
	OriginalEntry map[string]interface{} `json:"original,omitempty"`
}

// UnreadEntry represents an unread entry ID
type UnreadEntry struct {
	ID        int64     `json:"id"`
	EntryID   int64     `json:"entry_id"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UnreadEntryRequest represents a request to mark entries as unread
type UnreadEntryRequest struct {
	UnreadEntries []int64 `json:"unread_entries"`
}

// StarredEntry represents a starred entry ID
type StarredEntry struct {
	ID        int64     `json:"id"`
	EntryID   int64     `json:"entry_id"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// StarredEntryRequest represents a request to star entries
type StarredEntryRequest struct {
	StarredEntries []int64 `json:"starred_entries"`
}

// Tagging represents a tag applied to a feed
type Tagging struct {
	ID        int64     `json:"id"`
	FeedID    int64     `json:"feed_id"`
	TagID     int64     `json:"tag_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// TaggingRequest represents a request to create a tagging
type TaggingRequest struct {
	FeedID int64  `json:"feed_id"`
	Name   string `json:"name"`
}

// Tag represents a Feedbin tag
type Tag struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count,omitempty"`
}

// SavedSearch represents a saved search
type SavedSearch struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Query     string    `json:"query"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SavedSearchRequest represents a request to create a saved search
type SavedSearchRequest struct {
	Name  string `json:"name"`
	Query string `json:"query"`
}

// RecentlyReadEntry represents a recently read entry
type RecentlyReadEntry struct {
	ID        int64     `json:"id"`
	EntryID   int64     `json:"entry_id"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

// RecentlyReadEntryRequest represents a request to mark entries as recently read
type RecentlyReadEntryRequest struct {
	EntryIDs []int64 `json:"recently_read_entries"`
}

// UpdatedEntry represents an updated entry
type UpdatedEntry struct {
	ID        int64     `json:"id"`
	EntryID   int64     `json:"entry_id"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

// Favicon represents a feed favicon
type Favicon struct {
	ID       int64  `json:"id"`
	URL      string `json:"url"`
	DataURL  string `json:"data_url,omitempty"`
	Filename string `json:"filename,omitempty"`
}

// Import represents an OPML import
type Import struct {
	ID        int64     `json:"id"`
	Complete  bool      `json:"complete"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ImportRequest represents a request to create an import
type ImportRequest struct {
	OPML string `json:"opml"`
}

// Page represents a page
type Page struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	EntryID   int64     `json:"entry_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ErrorResponse represents an error response from the API
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
