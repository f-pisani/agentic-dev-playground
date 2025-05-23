package feedbin

import (
	"net/http"
	"time"
)

// Subscription represents a Feedbin subscription
type Subscription struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	FeedID    int       `json:"feed_id"`
	Title     string    `json:"title"`
	FeedURL   string    `json:"feed_url"`
	SiteURL   string    `json:"site_url"`
	JSONFeed  *JSONFeed `json:"json_feed,omitempty"`
}

// JSONFeed represents additional metadata for JSON feeds
type JSONFeed struct {
	Favicon     string `json:"favicon,omitempty"`
	FeedURL     string `json:"feed_url,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Version     string `json:"version,omitempty"`
	HomePageURL string `json:"home_page_url,omitempty"`
	Title       string `json:"title,omitempty"`
}

// Entry represents a Feedbin entry
type Entry struct {
	ID                  int                `json:"id"`
	FeedID              int                `json:"feed_id"`
	Title               *string            `json:"title"`
	URL                 string             `json:"url"`
	ExtractedContentURL string             `json:"extracted_content_url"`
	Author              *string            `json:"author"`
	Content             *string            `json:"content"`
	Summary             string             `json:"summary"`
	Published           time.Time          `json:"published"`
	CreatedAt           time.Time          `json:"created_at"`
	Original            *OriginalEntry     `json:"original,omitempty"`
	Images              *EntryImages       `json:"images,omitempty"`
	Enclosure           *Enclosure         `json:"enclosure,omitempty"`
	TwitterID           *int64             `json:"twitter_id,omitempty"`
	TwitterThreadIDs    []int64            `json:"twitter_thread_ids,omitempty"`
	ExtractedArticles   []ExtractedArticle `json:"extracted_articles,omitempty"`
	JSONFeed            *JSONFeed          `json:"json_feed,omitempty"`
}

// OriginalEntry represents the original entry data before updates
type OriginalEntry struct {
	Author    *string                `json:"author"`
	Content   *string                `json:"content"`
	Title     *string                `json:"title"`
	URL       string                 `json:"url"`
	EntryID   string                 `json:"entry_id"`
	Published time.Time              `json:"published"`
	Data      map[string]interface{} `json:"data"`
}

// EntryImages represents images associated with an entry
type EntryImages struct {
	OriginalURL string     `json:"original_url"`
	Size1       *ImageSize `json:"size_1,omitempty"`
}

// ImageSize represents a resized image
type ImageSize struct {
	CDNURL string `json:"cdn_url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// Enclosure represents podcast/RSS enclosure data
type Enclosure struct {
	EnclosureURL    string `json:"enclosure_url"`
	EnclosureType   string `json:"enclosure_type"`
	EnclosureLength string `json:"enclosure_length"`
	ItunesDuration  string `json:"itunes_duration,omitempty"`
	ItunesImage     string `json:"itunes_image,omitempty"`
}

// ExtractedArticle represents an extracted article from a tweet
type ExtractedArticle struct {
	URL     string `json:"url"`
	Title   string `json:"title"`
	Host    string `json:"host"`
	Author  string `json:"author"`
	Content string `json:"content"`
}

// Tagging represents a feed tag
type Tagging struct {
	ID     int    `json:"id"`
	FeedID int    `json:"feed_id"`
	Name   string `json:"name"`
}

// SavedSearch represents a saved search
type SavedSearch struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Query     string    `json:"query"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Icon represents a feed icon
type Icon struct {
	ID        int       `json:"id"`
	FeedID    int       `json:"feed_id"`
	Host      string    `json:"host"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Import represents an OPML import
type Import struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Page represents a page
type Page struct {
	ID        int       `json:"id"`
	URL       string    `json:"url"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Feed represents a feed
type Feed struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	FeedURL   string    `json:"feed_url"`
	SiteURL   string    `json:"site_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PaginationInfo contains pagination metadata
type PaginationInfo struct {
	CurrentPage int
	TotalCount  int
	NextURL     string
	PrevURL     string
	FirstURL    string
	LastURL     string
}

// PaginatedResponse wraps paginated API responses
type PaginatedResponse struct {
	Data       interface{}
	Pagination *PaginationInfo
}

// APIError represents an API error response
type APIError struct {
	StatusCode int
	Message    string
	Response   *http.Response
	Retryable  bool
}

// Error implements the error interface
func (e *APIError) Error() string {
	return e.Message
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
	Value   interface{}
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	return e.Message
}

// CacheManager defines the interface for HTTP caching
type CacheManager interface {
	Get(key string) (*CachedResponse, bool)
	Set(key string, response *CachedResponse)
}

// CachedResponse represents a cached HTTP response
type CachedResponse struct {
	Data         []byte
	ETag         string
	LastModified time.Time
}

// MemoryCache implements CacheManager using in-memory storage
type MemoryCache struct {
	cache map[string]*CachedResponse
}

// NewMemoryCache creates a new in-memory cache
func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		cache: make(map[string]*CachedResponse),
	}
}

// Get retrieves a cached response
func (m *MemoryCache) Get(key string) (*CachedResponse, bool) {
	response, found := m.cache[key]
	return response, found
}

// Set stores a response in the cache
func (m *MemoryCache) Set(key string, response *CachedResponse) {
	m.cache[key] = response
}

// SubscriptionOptions holds options for subscription requests
type SubscriptionOptions struct {
	Since *time.Time
	Mode  string
}

// EntryOptions holds options for entry requests
type EntryOptions struct {
	Page               *int
	Since              *time.Time
	IDs                []int
	Read               *bool
	Starred            *bool
	PerPage            *int
	Mode               string
	IncludeOriginal    bool
	IncludeEnclosure   bool
	IncludeContentDiff bool
}

// CreateSubscriptionRequest represents a subscription creation request
type CreateSubscriptionRequest struct {
	FeedURL string `json:"feed_url"`
}

// UpdateSubscriptionRequest represents a subscription update request
type UpdateSubscriptionRequest struct {
	Title string `json:"title"`
}

// CreateTaggingRequest represents a tagging creation request
type CreateTaggingRequest struct {
	FeedID int    `json:"feed_id"`
	Name   string `json:"name"`
}

// RenameTagRequest represents a tag rename request
type RenameTagRequest struct {
	OldName string `json:"old_name"`
	NewName string `json:"new_name"`
}

// DeleteTagRequest represents a tag deletion request
type DeleteTagRequest struct {
	Name string `json:"name"`
}

// UnreadEntriesRequest represents a request to mark entries as unread
type UnreadEntriesRequest struct {
	UnreadEntries []int `json:"unread_entries"`
}

// StarredEntriesRequest represents a request to star entries
type StarredEntriesRequest struct {
	StarredEntries []int `json:"starred_entries"`
}

// CreateSavedSearchRequest represents a saved search creation request
type CreateSavedSearchRequest struct {
	Name  string `json:"name"`
	Query string `json:"query"`
}

// UpdateSavedSearchRequest represents a saved search update request
type UpdateSavedSearchRequest struct {
	Name  string `json:"name"`
	Query string `json:"query"`
}

// CreatePageRequest represents a page creation request
type CreatePageRequest struct {
	URL string `json:"url"`
}

// MultipleChoicesResponse represents a response when multiple feeds are found
type MultipleChoicesResponse []struct {
	FeedURL string `json:"feed_url"`
	Title   string `json:"title"`
}
