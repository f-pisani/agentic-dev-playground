package feedbin

import "time"

// Entry represents a single feed entry.
type Entry struct {
	ID                  int64      `json:"id"`
	FeedID              int64      `json:"feed_id"`
	Title               *string    `json:"title"`
	URL                 *string    `json:"url"`
	ExtractedContentURL *string    `json:"extracted_content_url"`
	Author              *string    `json:"author"`
	Content             *string    `json:"content"` // HTML content
	Summary             *string    `json:"summary"`
	Published           *time.Time `json:"published"`
	CreatedAt           *time.Time `json:"created_at"`

	// Fields included with mode=extended
	Original          *OriginalEntry     `json:"original,omitempty"`
	Images            *EntryImages       `json:"images,omitempty"`
	Enclosure         *EntryEnclosure    `json:"enclosure,omitempty"`
	TwitterID         *int64             `json:"twitter_id,omitempty"`
	TwitterThreadIDs  []int64            `json:"twitter_thread_ids,omitempty"`
	ExtractedArticles []ExtractedArticle `json:"extracted_articles,omitempty"`
	JSONFeed          *JSONFeedMetadata  `json:"json_feed,omitempty"`

	// Fields included with include_original=true
	// Original           *OriginalEntry `json:"original,omitempty"` // Already defined above

	// Fields included with include_content_diff=true
	ContentDiff *string `json:"content_diff,omitempty"` // HTML diff
}

// OriginalEntry holds the data for the original version of an updated entry.
type OriginalEntry struct {
	Author    *string    `json:"author"`
	Content   *string    `json:"content"`
	Title     *string    `json:"title"`
	URL       *string    `json:"url"`
	EntryID   *string    `json:"entry_id"` // Often a tag URI
	Published *time.Time `json:"published"`
	Data      *string    `json:"data"` // Can be JSON or other structure
}

// EntryImages holds image data associated with an entry.
type EntryImages struct {
	OriginalURL *string    `json:"original_url"`
	Size1       *ImageSize `json:"size_1,omitempty"` // Example size, others might exist
	// Add other potential sizes if known (size_2, etc.)
}

// ImageSize represents a specific size of an image with its CDN URL.
type ImageSize struct {
	CDNURL *string `json:"cdn_url"`
	Width  *int    `json:"width"`
	Height *int    `json:"height"`
}

// EntryEnclosure holds podcast/enclosure related metadata.
type EntryEnclosure struct {
	EnclosureURL    *string `json:"enclosure_url"`
	EnclosureType   *string `json:"enclosure_type"`
	EnclosureLength *string `json:"enclosure_length"` // String because it might be large
	ITunesDuration  *string `json:"itunes_duration"`
	ITunesImage     *string `json:"itunes_image"`
}

// ExtractedArticle holds content extracted from links within an entry (often tweets).
type ExtractedArticle struct {
	URL     *string `json:"url"`
	Title   *string `json:"title"`
	Host    *string `json:"host"`
	Author  *string `json:"author"`
	Content *string `json:"content"` // HTML content
}

// JSONFeedMetadata holds additional metadata from JSON Feeds.
type JSONFeedMetadata map[string]interface{} // Use map for flexibility

// Feed represents a single feed.
type Feed struct {
	ID      int64   `json:"id"`
	Title   *string `json:"title"`
	FeedURL *string `json:"feed_url"`
	SiteURL *string `json:"site_url"`
}

// Icon represents a feed icon (favicon).
type Icon struct {
	Host string `json:"host"`
	URL  string `json:"url"`
}

// Import represents an OPML import task.
type Import struct {
	ID          int64        `json:"id"`
	Complete    bool         `json:"complete"`
	CreatedAt   time.Time    `json:"created_at"`
	ImportItems []ImportItem `json:"import_items,omitempty"` // Only present when getting specific import
}

// ImportItem represents a single item within an OPML import.
type ImportItem struct {
	Title   *string `json:"title"`
	FeedURL *string `json:"feed_url"`
	Status  string  `json:"status"` // "pending", "complete", "failed"
}

// Page represents a saved webpage, returned as an Entry.
// No separate struct needed, uses Entry.

// SavedSearch represents a saved search query.
type SavedSearch struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Query string `json:"query"`
}

// Subscription represents a user's subscription to a feed.
type Subscription struct {
	ID        int64            `json:"id"`
	CreatedAt time.Time        `json:"created_at"`
	FeedID    int64            `json:"feed_id"`
	Title     string           `json:"title"`
	FeedURL   string           `json:"feed_url"`
	SiteURL   string           `json:"site_url"`
	JSONFeed  JSONFeedMetadata `json:"json_feed,omitempty"` // Included with mode=extended
}

// Tagging represents a tag applied to a feed.
type Tagging struct {
	ID     int64  `json:"id"`
	FeedID int64  `json:"feed_id"`
	Name   string `json:"name"`
}

// ExtractedContent represents the content extracted by the Mercury Parser service.
type ExtractedContent struct {
	Title         *string    `json:"title"`
	Content       *string    `json:"content"` // HTML content
	Author        *string    `json:"author"`
	DatePublished *time.Time `json:"date_published"`
	LeadImageURL  *string    `json:"lead_image_url"`
	Dek           *string    `json:"dek"` // Deck/subtitle
	NextPageURL   *string    `json:"next_page_url"`
	URL           *string    `json:"url"`
	Domain        *string    `json:"domain"`
	Excerpt       *string    `json:"excerpt"`
	WordCount     *int       `json:"word_count"`
	Direction     *string    `json:"direction"` // "ltr" or "rtl"
	TotalPages    *int       `json:"total_pages"`
	RenderedPages *int       `json:"rendered_pages"`
}
