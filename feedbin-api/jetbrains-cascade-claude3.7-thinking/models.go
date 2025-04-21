package client

import (
	"time"
)

// Subscription represents a Feedbin subscription
type Subscription struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	FeedID    int64     `json:"feed_id"`
	Title     string    `json:"title"`
	FeedURL   string    `json:"feed_url"`
	SiteURL   string    `json:"site_url"`
	JSONFeed  *JSONFeed `json:"json_feed,omitempty"` // Extended mode only
}

// JSONFeed represents additional metadata available in JSON feeds
type JSONFeed struct {
	Favicon     string `json:"favicon,omitempty"`
	FeedURL     string `json:"feed_url,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Version     string `json:"version,omitempty"`
	HomePageURL string `json:"home_page_url,omitempty"`
	Title       string `json:"title,omitempty"`
}

// Entry represents a Feedbin entry (article)
type Entry struct {
	ID                  int64     `json:"id"`
	FeedID              int64     `json:"feed_id"`
	Title               *string   `json:"title"` // Can be null
	URL                 string    `json:"url"`
	ExtractedContentURL string    `json:"extracted_content_url"`
	Author              *string   `json:"author"`  // Can be null
	Content             *string   `json:"content"` // Can be null
	Summary             *string   `json:"summary"` // Can be null
	Published           time.Time `json:"published"`
	CreatedAt           time.Time `json:"created_at"`

	// Extended mode fields
	Original          *EntryOriginal `json:"original,omitempty"`
	TwitterID         *int64         `json:"twitter_id,omitempty"`
	TwitterThreadIDs  []int64        `json:"twitter_thread_ids,omitempty"`
	Images            *EntryImages   `json:"images,omitempty"`
	Enclosure         *Enclosure     `json:"enclosure,omitempty"`
	ExtractedArticles []Article      `json:"extracted_articles,omitempty"`
	JSONFeed          *JSONFeed      `json:"json_feed,omitempty"`
}

// EntryOriginal represents the original entry data if it has been updated
type EntryOriginal struct {
	Author    string                 `json:"author,omitempty"`
	Content   string                 `json:"content,omitempty"`
	Title     string                 `json:"title,omitempty"`
	URL       string                 `json:"url,omitempty"`
	EntryID   string                 `json:"entry_id,omitempty"`
	Published time.Time              `json:"published,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

// EntryImages represents images associated with an entry
type EntryImages struct {
	OriginalURL string     `json:"original_url,omitempty"`
	Size1       *ImageSize `json:"size_1,omitempty"`
}

// ImageSize represents size information for an image
type ImageSize struct {
	CDNURL string `json:"cdn_url,omitempty"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

// Enclosure represents podcast/RSS enclosure data
type Enclosure struct {
	EnclosureURL    string `json:"enclosure_url,omitempty"`
	EnclosureType   string `json:"enclosure_type,omitempty"`
	EnclosureLength string `json:"enclosure_length,omitempty"`
	ItunesDuration  string `json:"itunes_duration,omitempty"`
	ItunesImage     string `json:"itunes_image,omitempty"`
}

// Article represents an extracted article from a tweet
type Article struct {
	URL     string `json:"url,omitempty"`
	Title   string `json:"title,omitempty"`
	Host    string `json:"host,omitempty"`
	Author  string `json:"author,omitempty"`
	Content string `json:"content,omitempty"`
}

// Tag represents a Feedbin tag
type Tag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Tagging represents a relation between a feed and a tag
type Tagging struct {
	ID     int64  `json:"id"`
	FeedID int64  `json:"feed_id"`
	TagID  int64  `json:"tag_id"`
	Name   string `json:"name"` // Tag name
}

// SavedSearch represents a saved search in Feedbin
type SavedSearch struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Query     string    `json:"query"`
	CreatedAt time.Time `json:"created_at"`
}

// ImportStatus represents the status of a feed import
type ImportStatus struct {
	ID              int64     `json:"id"`
	CompletePercent float64   `json:"complete_percent"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Page represents a saved page in Feedbin
type Page struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	URL      string `json:"url"`
	Content  string `json:"content,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
}

// Icon represents a feed icon
type Icon struct {
	ID     int64  `json:"id"`
	URL    string `json:"url"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}
