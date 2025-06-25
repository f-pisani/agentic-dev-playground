package feedbin

import "time"

// Subscription represents a user's subscription to a feed.
type Subscription struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	FeedID    int64     `json:"feed_id"`
	Title     string    `json:"title"`
	FeedURL   string    `json:"feed_url"`
	SiteURL   string    `json:"site_url"`
	JSONFeed  *JSONFeed `json:"json_feed,omitempty"`
}

// JSONFeed holds metadata for a JSON Feed.
type JSONFeed struct {
	Version     string `json:"version"`
	HomePageURL string `json:"home_page_url"`
	FeedURL     string `json:"feed_url"`
	Title       string `json:"title"`
	Icon        string `json:"icon"`
	Favicon     string `json:"favicon"`
}

// Entry represents a single feed entry.
type Entry struct {
	ID                  int64               `json:"id"`
	FeedID              int64               `json:"feed_id"`
	Title               *string             `json:"title"`
	URL                 string              `json:"url"`
	ExtractedContentURL string              `json:"extracted_content_url"`
	Author              *string             `json:"author"`
	Content             *string             `json:"content"`
	Summary             string              `json:"summary"`
	Published           time.Time           `json:"published"`
	CreatedAt           time.Time           `json:"created_at"`
	Original            *OriginalEntry      `json:"original,omitempty"`
	TwitterID           int64               `json:"twitter_id,omitempty"`
	TwitterThreadIDs    []int64             `json:"twitter_thread_ids,omitempty"`
	Images              *Images             `json:"images,omitempty"`
	Enclosure           *Enclosure          `json:"enclosure,omitempty"`
	ExtractedArticles   []*ExtractedArticle `json:"extracted_articles,omitempty"`
	JSONFeed            *JSONFeed           `json:"json_feed,omitempty"`
	ContentDiff         string              `json:"content_diff,omitempty"`
}

// OriginalEntry holds the original data of an entry that has been updated.
type OriginalEntry struct {
	Author    string    `json:"author"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	EntryID   string    `json:"entry_id"`
	Published time.Time `json:"published"`
	Data      *string   `json:"data"` // Can be null
}

// Images holds URLs and metadata for an entry's associated image.
type Images struct {
	OriginalURL string    `json:"original_url"`
	Size1       ImageSize `json:"size_1"`
}

// ImageSize holds the CDN URL and dimensions for a specific image size.
type ImageSize struct {
	CDNURL string `json:"cdn_url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// Enclosure holds podcast-related metadata.
type Enclosure struct {
	URL      string `json:"enclosure_url"`
	Type     string `json:"enclosure_type"`
	Length   string `json:"enclosure_length"`
	Duration string `json:"itunes_duration,omitempty"`
	Image    string `json:"itunes_image,omitempty"`
}

// ExtractedArticle holds content extracted from a URL within an entry.
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

// SavedSearch represents a user-saved search query.
type SavedSearch struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Query string `json:"query"`
}

// Icon represents a feed's favicon.
type Icon struct {
	Host string `json:"host"`
	URL  string `json:"url"`
}

// Import represents an OPML import job.
type Import struct {
	ID          int64        `json:"id"`
	Complete    bool         `json:"complete"`
	CreatedAt   time.Time    `json:"created_at"`
	ImportItems []ImportItem `json:"import_items"`
}

// ImportItem represents a single item within an OPML import.
type ImportItem struct {
	Title   string `json:"title"`
	FeedURL string `json:"feed_url"`
	Status  string `json:"status"`
}

// Page represents a webpage to be saved as an entry.
type Page struct {
	URL   string `json:"url"`
	Title string `json:"title,omitempty"`
}

// Feed represents a single feed's metadata.
type Feed struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	FeedURL string `json:"feed_url"`
	SiteURL string `json:"site_url"`
}
