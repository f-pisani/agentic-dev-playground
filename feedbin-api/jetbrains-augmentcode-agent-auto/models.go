package feedbin

import (
	"time"
)

// Subscription represents a Feedbin subscription.
type Subscription struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	FeedID    int64     `json:"feed_id"`
	Title     string    `json:"title"`
	FeedURL   string    `json:"feed_url"`
	SiteURL   string    `json:"site_url"`
	JSONFeed  *JSONFeed `json:"json_feed,omitempty"`
}

// JSONFeed represents additional metadata for a JSON Feed.
type JSONFeed struct {
	Favicon     string `json:"favicon,omitempty"`
	FeedURL     string `json:"feed_url,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Version     string `json:"version,omitempty"`
	HomePageURL string `json:"home_page_url,omitempty"`
	Title       string `json:"title,omitempty"`
}

// Feed represents a Feedbin feed.
type Feed struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	FeedURL string `json:"feed_url"`
	SiteURL string `json:"site_url"`
}

// Entry represents a Feedbin entry.
type Entry struct {
	ID                  int64      `json:"id"`
	FeedID              int64      `json:"feed_id"`
	Title               string     `json:"title"`
	URL                 string     `json:"url"`
	ExtractedContentURL string     `json:"extracted_content_url"`
	Author              string     `json:"author"`
	Summary             string     `json:"summary"`
	Content             string     `json:"content"`
	Published           time.Time  `json:"published"`
	CreatedAt           time.Time  `json:"created_at"`
	Original            *Entry     `json:"original,omitempty"`
	ContentDiff         string     `json:"content_diff,omitempty"`
	Images              []string   `json:"images,omitempty"`
	Enclosure           *Enclosure `json:"enclosure,omitempty"`
	TwitterID           int64      `json:"twitter_id,omitempty"`
	TwitterThreadIDs    []int64    `json:"twitter_thread_ids,omitempty"`
	ExtractedArticles   []Article  `json:"extracted_articles,omitempty"`
	JSONFeed            *JSONFeed  `json:"json_feed,omitempty"`
}

// Enclosure represents podcast/RSS enclosure data.
type Enclosure struct {
	URL      string `json:"url"`
	Type     string `json:"type"`
	Length   string `json:"length"`
	Duration int    `json:"duration,omitempty"`
}

// Article represents an extracted article.
type Article struct {
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Author        string    `json:"author"`
	DatePublished time.Time `json:"date_published"`
	LeadImageURL  string    `json:"lead_image_url"`
	Dek           string    `json:"dek"`
	NextPageURL   string    `json:"next_page_url"`
	URL           string    `json:"url"`
	Domain        string    `json:"domain"`
	Excerpt       string    `json:"excerpt"`
	WordCount     int       `json:"word_count"`
	Direction     string    `json:"direction"`
	TotalPages    int       `json:"total_pages"`
	RenderedPages int       `json:"rendered_pages"`
}

// Tagging represents a Feedbin tagging.
type Tagging struct {
	ID     int64  `json:"id"`
	FeedID int64  `json:"feed_id"`
	Name   string `json:"name"`
}

// SavedSearch represents a Feedbin saved search.
type SavedSearch struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Query string `json:"query"`
}

// Icon represents a Feedbin feed icon.
type Icon struct {
	Host string `json:"host"`
	URL  string `json:"url"`
}

// Import represents a Feedbin OPML import.
type Import struct {
	ID          int64        `json:"id"`
	Complete    bool         `json:"complete"`
	CreatedAt   time.Time    `json:"created_at"`
	ImportItems []ImportItem `json:"import_items,omitempty"`
}

// ImportItem represents an item in a Feedbin OPML import.
type ImportItem struct {
	Title   string `json:"title"`
	FeedURL string `json:"feed_url"`
	Status  string `json:"status"`
}

// Page represents a Feedbin page.
type Page struct {
	URL   string `json:"url"`
	Title string `json:"title,omitempty"`
}
