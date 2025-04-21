// Package feedbin provides a Go client for the Feedbin API v2.
package feedbin

import (
	"time"
)

// Subscription represents a Feedbin subscription.
type Subscription struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	FeedID    int       `json:"feed_id"`
	Title     string    `json:"title"`
	FeedURL   string    `json:"feed_url"`
	SiteURL   string    `json:"site_url"`
	// Extended mode fields
	JSONFeed *JSONFeed `json:"json_feed,omitempty"`
}

// JSONFeed represents the additional metadata available in extended mode.
type JSONFeed struct {
	Favicon     string `json:"favicon,omitempty"`
	FeedURL     string `json:"feed_url,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Version     string `json:"version,omitempty"`
	HomePageURL string `json:"home_page_url,omitempty"`
	Title       string `json:"title,omitempty"`
}

// Entry represents a Feedbin entry.
type Entry struct {
	ID                  int       `json:"id"`
	FeedID              int       `json:"feed_id"`
	Title               *string   `json:"title"`
	URL                 string    `json:"url"`
	ExtractedContentURL string    `json:"extracted_content_url"`
	Author              *string   `json:"author"`
	Content             *string   `json:"content"`
	Summary             *string   `json:"summary"`
	Published           time.Time `json:"published"`
	CreatedAt           time.Time `json:"created_at"`
	// Extended mode fields
	Original          *EntryOriginal     `json:"original,omitempty"`
	Images            *EntryImages       `json:"images,omitempty"`
	Enclosure         *EntryEnclosure    `json:"enclosure,omitempty"`
	TwitterID         *int64             `json:"twitter_id,omitempty"`
	TwitterThreadIDs  []int64            `json:"twitter_thread_ids,omitempty"`
	ExtractedArticles []ExtractedArticle `json:"extracted_articles,omitempty"`
	JSONFeed          *JSONFeed          `json:"json_feed,omitempty"`
}

// EntryOriginal represents the original entry data if the entry has been updated.
type EntryOriginal struct {
	Author    string      `json:"author"`
	Content   string      `json:"content"`
	Title     string      `json:"title"`
	URL       string      `json:"url"`
	EntryID   string      `json:"entry_id"`
	Published time.Time   `json:"published"`
	Data      interface{} `json:"data"`
}

// EntryImages represents the images associated with an entry.
type EntryImages struct {
	OriginalURL string     `json:"original_url"`
	Size1       *ImageSize `json:"size_1,omitempty"`
}

// ImageSize represents an image size.
type ImageSize struct {
	CDNURL string `json:"cdn_url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// EntryEnclosure represents podcast/RSS enclosure data.
type EntryEnclosure struct {
	EnclosureURL    string `json:"enclosure_url"`
	EnclosureType   string `json:"enclosure_type"`
	EnclosureLength string `json:"enclosure_length"`
	ItunesDuration  string `json:"itunes_duration,omitempty"`
	ItunesImage     string `json:"itunes_image,omitempty"`
}

// ExtractedArticle represents an article extracted from a tweet.
type ExtractedArticle struct {
	URL     string `json:"url"`
	Title   string `json:"title"`
	Host    string `json:"host"`
	Author  string `json:"author"`
	Content string `json:"content"`
}

// Tag represents a Feedbin tag.
type Tag struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	FeedIDs []int  `json:"feed_ids"`
}

// Tagging represents a Feedbin tagging (association between a tag and a feed).
type Tagging struct {
	ID     int    `json:"id"`
	FeedID int    `json:"feed_id"`
	TagID  int    `json:"tag_id"`
	Name   string `json:"name"`
}

// SavedSearch represents a Feedbin saved search.
type SavedSearch struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Query     string    `json:"query"`
	CreatedAt time.Time `json:"created_at"`
}

// Icon represents a Feedbin feed icon.
type Icon struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

// Import represents a Feedbin import.
type Import struct {
	ID           int       `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Complete     bool      `json:"complete"`
	Success      bool      `json:"success"`
	ErrorMessage string    `json:"error_message,omitempty"`
	OpmlURL      string    `json:"opml_url,omitempty"`
	OpmlFile     string    `json:"opml_file,omitempty"`
}

// Page represents a Feedbin page.
type Page struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	Body      string    `json:"body"`
	Published time.Time `json:"published"`
}

// PaginationLinks represents the pagination links from the Link header.
type PaginationLinks struct {
	First string
	Prev  string
	Next  string
	Last  string
}
