package feedbin

import (
	"time"
)

// Subscription represents a feed subscription in Feedbin
type Subscription struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	FeedID    int       `json:"feed_id"`
	Title     string    `json:"title"`
	FeedURL   string    `json:"feed_url"`
	SiteURL   string    `json:"site_url"`
	JSONFeed  *JSONFeed `json:"json_feed,omitempty"`
}

// JSONFeed represents additional metadata for a JSON Feed
type JSONFeed struct {
	Favicon     string `json:"favicon,omitempty"`
	FeedURL     string `json:"feed_url,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Version     string `json:"version,omitempty"`
	HomePageURL string `json:"home_page_url,omitempty"`
	Title       string `json:"title,omitempty"`
}

// Entry represents an entry in a feed
type Entry struct {
	ID                  int        `json:"id"`
	FeedID              int        `json:"feed_id"`
	Title               *string    `json:"title"`
	URL                 string     `json:"url"`
	ExtractedContentURL string     `json:"extracted_content_url"`
	Author              *string    `json:"author"`
	Content             *string    `json:"content"`
	Summary             *string    `json:"summary"`
	Published           time.Time  `json:"published"`
	CreatedAt           time.Time  `json:"created_at"`
	Original            *Original  `json:"original,omitempty"`
	Images              *Images    `json:"images,omitempty"`
	Enclosure           *Enclosure `json:"enclosure,omitempty"`
	TwitterID           *int64     `json:"twitter_id,omitempty"`
	TwitterThreadIDs    []int64    `json:"twitter_thread_ids,omitempty"`
	ExtractedArticles   []Article  `json:"extracted_articles,omitempty"`
	JSONFeed            *JSONFeed  `json:"json_feed,omitempty"`
}

// Original represents the original entry data if the entry has been updated
type Original struct {
	Author    string                 `json:"author"`
	Content   string                 `json:"content"`
	Title     string                 `json:"title"`
	URL       string                 `json:"url"`
	EntryID   string                 `json:"entry_id"`
	Published time.Time              `json:"published"`
	Data      map[string]interface{} `json:"data"`
}

// Images represents images associated with an entry
type Images struct {
	OriginalURL string     `json:"original_url"`
	Size1       *ImageSize `json:"size_1,omitempty"`
}

// ImageSize represents a specific image size
type ImageSize struct {
	CDNURL string `json:"cdn_url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// Enclosure represents podcast/RSS enclosure data
type Enclosure struct {
	URL            string `json:"enclosure_url"`
	Type           string `json:"enclosure_type"`
	Length         string `json:"enclosure_length"`
	ITunesDuration string `json:"itunes_duration,omitempty"`
	ITunesImage    string `json:"itunes_image,omitempty"`
}

// Article represents an extracted article from a tweet
type Article struct {
	URL     string `json:"url"`
	Title   string `json:"title"`
	Host    string `json:"host"`
	Author  string `json:"author"`
	Content string `json:"content"`
}

// Tag represents a tag in Feedbin
type Tag struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	FeedIDs []int  `json:"feed_ids,omitempty"`
}

// Tagging represents a connection between a tag and a feed
type Tagging struct {
	ID     int    `json:"id"`
	FeedID int    `json:"feed_id"`
	TagID  int    `json:"tag_id"`
	Name   string `json:"name,omitempty"`
}

// SavedSearch represents a saved search
type SavedSearch struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Query     string    `json:"query"`
	CreatedAt time.Time `json:"created_at"`
}

// ExtractedContent represents content extracted from an article
type ExtractedContent struct {
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Author        string    `json:"author"`
	DatePublished time.Time `json:"date_published"`
	LeadImageURL  *string   `json:"lead_image_url"`
	Dek           *string   `json:"dek"`
	NextPageURL   *string   `json:"next_page_url"`
	URL           string    `json:"url"`
	Domain        string    `json:"domain"`
	Excerpt       string    `json:"excerpt"`
	WordCount     int       `json:"word_count"`
	Direction     string    `json:"direction"`
	TotalPages    int       `json:"total_pages"`
	RenderedPages int       `json:"rendered_pages"`
}

// Page represents a page in Feedbin
type Page struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
}

// Icon represents a feed icon
type Icon struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}
