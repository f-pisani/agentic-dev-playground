package feedbin

import "time"

// Subscription represents a feed subscription.
type Subscription struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	FeedID    int64     `json:"feed_id"`
	Title     string    `json:"title"`
	FeedURL   string    `json:"feed_url"`
	SiteURL   string    `json:"site_url"`
	JSONFeed  *JSONFeed `json:"json_feed,omitempty"`
}

// JSONFeed contains additional metadata for JSON feeds.
type JSONFeed struct {
	Favicon     string `json:"favicon"`
	FeedURL     string `json:"feed_url"`
	Icon        string `json:"icon"`
	Version     string `json:"version"`
	HomePageURL string `json:"home_page_url"`
	Title       string `json:"title"`
}

// Entry represents a feed entry.
type Entry struct {
	ID                  int64              `json:"id"`
	FeedID              int64              `json:"feed_id"`
	Title               *string            `json:"title"`
	URL                 *string            `json:"url"`
	ExtractedContentURL *string            `json:"extracted_content_url"`
	Author              *string            `json:"author"`
	Content             *string            `json:"content"`
	Summary             *string            `json:"summary"`
	Published           time.Time          `json:"published"`
	CreatedAt           time.Time          `json:"created_at"`
	Original            *OriginalEntry     `json:"original,omitempty"`
	Images              *EntryImages       `json:"images,omitempty"`
	Enclosure           *EntryEnclosure    `json:"enclosure,omitempty"`
	TwitterID           int64              `json:"twitter_id,omitempty"`
	TwitterThreadIDs    []int64            `json:"twitter_thread_ids,omitempty"`
	ExtractedArticles   []ExtractedArticle `json:"extracted_articles,omitempty"`
}

// OriginalEntry contains the original entry data if it has been updated.
type OriginalEntry struct {
	Author    string                 `json:"author"`
	Content   string                 `json:"content"`
	Title     string                 `json:"title"`
	URL       string                 `json:"url"`
	EntryID   string                 `json:"entry_id"`
	Published time.Time              `json:"published"`
	Data      map[string]interface{} `json:"data"`
}

// EntryImages contains images extracted by Feedbin.
type EntryImages struct {
	OriginalURL string `json:"original_url"`
	Size1       *Image `json:"size_1,omitempty"`
}

// Image represents an image resource.
type Image struct {
	CDNURL string `json:"cdn_url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// EntryEnclosure represents a podcast or RSS enclosure.
type EntryEnclosure struct {
	EnclosureURL    string `json:"enclosure_url"`
	EnclosureType   string `json:"enclosure_type"`
	EnclosureLength string `json:"enclosure_length"`
	ItunesDuration  string `json:"itunes_duration"`
	ItunesImage     string `json:"itunes_image"`
}

// ExtractedArticle represents content extracted from a link (e.g. in a tweet).
type ExtractedArticle struct {
	URL     string `json:"url"`
	Title   string `json:"title"`
	Host    string `json:"host"`
	Author  string `json:"author"`
	Content string `json:"content"`
}

// Tagging represents a tag applied to a feed.
type Tagging struct {
	ID     int64  `json:"id"`
	FeedID int64  `json:"feed_id"`
	Name   string `json:"name"`
}

// SavedSearch represents a saved search.
type SavedSearch struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Query string `json:"query"`
}

// Icon represents a feed icon.
type Icon struct {
	Host string `json:"host"`
	URL  string `json:"url"`
}

// ImportStatus constants.
const (
	ImportStatusPending  = "pending"
	ImportStatusComplete = "complete"
	ImportStatusFailed   = "failed"
)

// Import represents an OPML import.
type Import struct {
	ID          int64        `json:"id"`
	Complete    bool         `json:"complete"`
	CreatedAt   time.Time    `json:"created_at"`
	ImportItems []ImportItem `json:"import_items,omitempty"`
}

// ImportItem represents a single feed in an import.
type ImportItem struct {
	Title   string `json:"title"`
	FeedURL string `json:"feed_url"`
	Status  string `json:"status"` // pending, complete, failed
}
