package feedbin

import (
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
	JSONFeed  *JSONFeed `json:"json_feed,omitempty"` // Available in extended mode
}

// JSONFeed contains JSON Feed metadata
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
	ID                  int        `json:"id"`
	FeedID              int        `json:"feed_id"`
	Title               *string    `json:"title"` // Can be null
	URL                 string     `json:"url"`
	ExtractedContentURL string     `json:"extracted_content_url"`
	Author              *string    `json:"author"`  // Can be null
	Content             *string    `json:"content"` // Can be null
	Summary             *string    `json:"summary"` // Can be null
	Published           time.Time  `json:"published"`
	CreatedAt           time.Time  `json:"created_at"`
	Original            *Original  `json:"original,omitempty"`           // Only with include_original=true
	Enclosure           *Enclosure `json:"enclosure,omitempty"`          // Only with include_enclosure=true
	Images              *Images    `json:"images,omitempty"`             // Only in extended mode
	TwitterID           *int64     `json:"twitter_id,omitempty"`         // Only in extended mode
	TwitterThreadIDs    []int64    `json:"twitter_thread_ids,omitempty"` // Only in extended mode
	ExtractedArticles   []Article  `json:"extracted_articles,omitempty"` // Only in extended mode
	JSONFeed            *JSONFeed  `json:"json_feed,omitempty"`          // Only in extended mode
}

// Original represents the original entry data if the entry has been updated
type Original struct {
	Author    *string     `json:"author"`
	Content   *string     `json:"content"`
	Title     *string     `json:"title"`
	URL       string      `json:"url"`
	EntryID   string      `json:"entry_id"`
	Published time.Time   `json:"published"`
	Data      interface{} `json:"data"`
}

// Enclosure represents podcast/RSS enclosure data
type Enclosure struct {
	EnclosureURL    string `json:"enclosure_url"`
	EnclosureType   string `json:"enclosure_type,omitempty"`
	EnclosureLength string `json:"enclosure_length,omitempty"`
	ItunesDuration  string `json:"itunes_duration,omitempty"`
	ItunesImage     string `json:"itunes_image,omitempty"`
}

// Images represents images associated with an entry
type Images struct {
	OriginalURL string           `json:"original_url"`
	Size1       *ImageAttributes `json:"size_1,omitempty"`
}

// ImageAttributes represents attributes of an image
type ImageAttributes struct {
	CDNURL string `json:"cdn_url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// Article represents an extracted article
type Article struct {
	URL     string `json:"url"`
	Title   string `json:"title"`
	Host    string `json:"host"`
	Author  string `json:"author"`
	Content string `json:"content"`
}

// Tagging represents a tag assignment to a feed
type Tagging struct {
	ID     int    `json:"id"`
	FeedID int    `json:"feed_id"`
	Name   string `json:"name"`
}

// SavedSearch represents a saved search
type SavedSearch struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Query string `json:"query"`
}
