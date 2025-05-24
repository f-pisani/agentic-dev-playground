package feedbin

import (
	"time"
)

// Entry represents a feed entry/article
type Entry struct {
	ID                  int                    `json:"id"`
	FeedID              int                    `json:"feed_id"`
	Title               string                 `json:"title"`
	Author              string                 `json:"author,omitempty"`
	Content             string                 `json:"content"`
	Summary             string                 `json:"summary"`
	URL                 string                 `json:"url"`
	ExtractedContentURL string                 `json:"extracted_content_url,omitempty"`
	Published           time.Time              `json:"published"`
	CreatedAt           time.Time              `json:"created_at"`
	Original            *OriginalEntry         `json:"original,omitempty"`
	Images              map[string]string      `json:"images,omitempty"`
	Enclosure           *Enclosure             `json:"enclosure,omitempty"`
	TwitterID           string                 `json:"twitter_id,omitempty"`
	TwitterThreadIDs    []string               `json:"twitter_thread_ids,omitempty"`
	ExtractedArticles   []ExtractedArticle     `json:"extracted_articles,omitempty"`
	JSONFeed            map[string]interface{} `json:"json_feed,omitempty"`
}

// OriginalEntry contains the original entry data before updates
type OriginalEntry struct {
	Author          string    `json:"author,omitempty"`
	Content         string    `json:"content"`
	Title           string    `json:"title"`
	URL             string    `json:"url"`
	EntryID         string    `json:"entry_id"`
	PublishedAt     time.Time `json:"published_at"`
	Data            string    `json:"data"`
	ExtractedAuthor string    `json:"extracted_author,omitempty"`
}

// Enclosure represents podcast/media attachments
type Enclosure struct {
	URL      string `json:"url"`
	Type     string `json:"type"`
	Length   string `json:"length"`
	Title    string `json:"title,omitempty"`
	Duration string `json:"duration,omitempty"`
}

// ExtractedArticle represents content extracted via Mercury Parser
type ExtractedArticle struct {
	URL     string `json:"url"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author,omitempty"`
}

// Subscription represents a feed subscription
type Subscription struct {
	ID        int                    `json:"id"`
	FeedID    int                    `json:"feed_id"`
	Title     string                 `json:"title"`
	FeedURL   string                 `json:"feed_url"`
	SiteURL   string                 `json:"site_url"`
	CreatedAt time.Time              `json:"created_at"`
	JSONFeed  map[string]interface{} `json:"json_feed,omitempty"`
}

// Feed represents a feed
type Feed struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	FeedURL string `json:"feed_url"`
	SiteURL string `json:"site_url"`
}

// Tag represents a tag for organizing feeds
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Tagging represents the association between a feed and a tag
type Tagging struct {
	ID     int    `json:"id"`
	FeedID int    `json:"feed_id"`
	Name   string `json:"name"`
}

// SavedSearch represents a saved search query
type SavedSearch struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Query string `json:"query"`
}

// Icon represents a feed's favicon
type Icon struct {
	Host string `json:"host"`
	Icon string `json:"icon"`
}

// Import represents an OPML import operation
type Import struct {
	ID          int          `json:"id"`
	Complete    bool         `json:"complete"`
	CreatedAt   time.Time    `json:"created_at"`
	ImportItems []ImportItem `json:"import_items"`
}

// ImportItem represents a single item in an import
type ImportItem struct {
	ID            int     `json:"id"`
	ImportID      int     `json:"import_id"`
	Details       string  `json:"details"`
	FeedID        *int    `json:"feed_id,omitempty"`
	Status        string  `json:"status"`
	OriginalTitle string  `json:"original_title"`
	OriginalURL   string  `json:"original_url"`
	Message       *string `json:"message,omitempty"`
}

// Page represents a saved web page
type Page struct {
	ID    int    `json:"id"`
	URL   string `json:"url"`
	Title string `json:"title,omitempty"`
}

// RecentlyReadEntry represents a recently read entry ID
type RecentlyReadEntry struct {
	ID int `json:"id"`
}

// UpdatedEntry represents an updated entry
type UpdatedEntry struct {
	ID               int            `json:"id"`
	FeedID           int            `json:"feed_id"`
	Title            string         `json:"title"`
	URL              string         `json:"url"`
	PublishedAt      time.Time      `json:"published_at"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	OriginalTitle    string         `json:"original_title,omitempty"`
	OriginalContent  string         `json:"original_content,omitempty"`
	Content          string         `json:"content"`
	TitleContentDiff string         `json:"title_content_diff,omitempty"`
	Original         *OriginalEntry `json:"original,omitempty"`
}

// Authentication represents the authentication check response
type Authentication struct {
	Email string `json:"email,omitempty"`
}

// CreateSubscriptionRequest represents a subscription creation request
type CreateSubscriptionRequest struct {
	FeedURL string `json:"feed_url"`
}

// UpdateSubscriptionRequest represents a subscription update request
type UpdateSubscriptionRequest struct {
	Title string `json:"title,omitempty"`
}

// MultipleChoiceItem represents an item in a multiple choice response
type MultipleChoiceItem struct {
	FeedURL string `json:"feed_url"`
	Title   string `json:"title"`
}

// ExtractedContent represents content extracted via Mercury Parser
type ExtractedContent struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author,omitempty"`
	URL     string `json:"url"`
}

// EntryFilter represents filtering options for entries
type EntryFilter struct {
	IDs                []int
	Read               *bool
	Starred            *bool
	Since              *time.Time
	PerPage            int
	Page               int
	Mode               string
	IncludeOriginal    bool
	IncludeEnclosure   bool
	IncludeContentDiff bool
}
